package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

// TokenBlacklistRepository defines the methods to interact with the token blacklist storage.
type TokenBlacklistRepository interface {
	BlacklistToken(ctx context.Context, userID, token string) error
}

type tokenBlacklistRepo struct {
	DB *sql.DB
}

// NewTokenBlacklistRepository creates a new instance of TokenBlacklistRepository.
func NewTokenBlacklistRepository(db *sql.DB) TokenBlacklistRepository {
	if db == nil {
		log.Fatal("Database connection cannot be nil")
	}
	return &tokenBlacklistRepo{DB: db}
}

// BlacklistToken adds a token to the token_blacklist table
func (r *tokenBlacklistRepo) BlacklistToken(ctx context.Context, userID, token string) error {
	log.Printf("BlacklistToken: Received userID: %s, token: %s", userID, token)

	if userID == "" {
		log.Println("BlacklistToken: userID is empty")
		return errors.New("user ID cannot be empty")
	}
	if token == "" {
		log.Println("BlacklistToken: token is empty")
		return errors.New("token cannot be empty")
	}

	log.Printf("BlacklistToken: Blacklisting token for userID: %s", userID)
	_, err := r.DB.ExecContext(ctx, `INSERT INTO token_blacklist (user_id, token) VALUES ($1, $2)`, userID, token)
	if err != nil {
		log.Printf("BlacklistToken: Error blacklisting token for userID: %s, error: %v", userID, err)
		return errors.New("database error: " + err.Error())
	}

	log.Printf("BlacklistToken: Token successfully blacklisted for userID: %s", userID)
	return nil
}
