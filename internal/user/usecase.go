package user

import (
	"context"

	"github.com/bayujo/synapsis-marketplace/internal/middleware"
	"github.com/bayujo/synapsis-marketplace/pkg/error"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, username, email, password string) error
	AuthenticateUser(ctx context.Context, username, password string) (string, error)
}

type userUsecase struct {
	userRepository UserRepository
	authService    middleware.AuthService
}

func NewUserUsecase(userRepository UserRepository, authService middleware.AuthService) UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		authService:    authService,
	}
}

func (uc *userUsecase) RegisterUser(ctx context.Context, username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.ErrInvalidInput
	}

	_, err := uc.userRepository.FindByUsername(ctx, username)
	if err == nil {
		return errors.ErrUsernameExists
	}
	_, err = uc.userRepository.FindByEmail(ctx, email)
	if err == nil {
		return errors.ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	err = uc.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userUsecase) AuthenticateUser(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepository.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.ErrInvalidCredentials
	}

	token, err := uc.authService.GenerateToken(ctx, username)
	if err != nil {
		return "", err
	}

	return token, nil
}
