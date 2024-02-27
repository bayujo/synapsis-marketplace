package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
    FindByUsername(ctx context.Context, username string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Create(ctx context.Context, user *User) error
}

type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
    query := "SELECT id, username, email, password FROM users WHERE username = $1"
    row := r.db.QueryRow(query, username)
    user := &User{}
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
    query := "SELECT id, username, email, password FROM users WHERE email = $1"
    row := r.db.QueryRow(query, email)
    user := &User{}
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
    query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
    err := r.db.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID)
    if err != nil {
        return err
    }
    return nil
}
