package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	envconv "github.com/bayujo/synapsis-marketplace/pkg/timeconv"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GenerateToken(ctx context.Context, username string) (string, error)
	ValidatePassword(ctx context.Context, plainPassword, hashedPassword string) bool
}

type authService struct {
	secretKey string
	expiresIn string
}

func NewAuthService(secretKey string, expiresIn string) AuthService {
	return &authService{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

func (s *authService) GenerateToken(ctx context.Context, username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(envconv.ParseDuration(s.expiresIn)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s *authService) ValidatePassword(ctx context.Context, plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized: Token not present", http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized: Error parsing token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized: Token is not valid", http.StatusUnauthorized)
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "Unauthorized: Error parsing claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}