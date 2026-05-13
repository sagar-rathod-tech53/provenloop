package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

// CreateUser saves a new user into the database
func (r *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (id, email, username, password_hash, created_at, updated_at) 
	          VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5)`
	_, err := r.DB.ExecContext(ctx, query, user.Email, user.Username, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmailOrUsername(identifier string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE email = $1 OR username = $1`
	row := r.DB.QueryRow(query, identifier)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE email = $1 or username = $1`
	row := r.DB.QueryRow(query, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, email, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE email = $3`
	_, err := r.DB.ExecContext(ctx, query, passwordHash, time.Now(), email)
	return err
}
