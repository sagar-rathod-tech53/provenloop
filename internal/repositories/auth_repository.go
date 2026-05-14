package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CheckUserExists(
	ctx context.Context,
	email string,
	username string,
) (bool, error) {

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE email = $1
		OR username = $2
	)
	`

	var exists bool

	err := r.DB.QueryRowContext(
		ctx,
		query,
		email,
		username,
	).Scan(&exists)

	return exists, err
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	user models.User,
) error {

	query := `
INSERT INTO users
(
	id,
	email,
	username,
	password_hash,
	is_verified,
	created_at,
	updated_at
)
VALUES($1,$2,$3,$4,$5,$6,$7)
`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		false,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) VerifyUser(
	ctx context.Context,
	email string,
) error {

	query := `
	UPDATE users
	SET is_verified = true
	WHERE email = $1
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		email,
	)

	return err
}

func (r *UserRepository) GetUserByEmailOrUsername(
	ctx context.Context,
	emailOrUsername string,
) (*models.User, error) {

	query := `
	SELECT
		id,
		email,
		username,
		password_hash,
		is_verified,
		created_at,
		updated_at
	FROM users
	WHERE email = $1
	OR username = $1
	`

	row := r.DB.QueryRowContext(
		ctx,
		query,
		emailOrUsername,
	)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CheckUserExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE email = $1
	)
	`

	var exists bool

	err := r.DB.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(&exists)

	return exists, err
}

func (r *UserRepository) UpdatePassword(
	ctx context.Context,
	email string,
	passwordHash string,
) error {

	query := `
	UPDATE users
	SET
		password_hash = $1,
		updated_at = $2
	WHERE email = $3
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		passwordHash,
		time.Now(),
		email,
	)

	return err
}

func (r *UserRepository) GetUserByID(
	ctx context.Context,
	userID string,
) (*models.User, error) {

	query := `
	SELECT
		id,
		email,
		username,
		password_hash,
		is_verified,
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	`

	row := r.DB.QueryRowContext(
		ctx,
		query,
		userID,
	)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ChangePassword(
	ctx context.Context,
	userID string,
	newPasswordHash string,
) error {

	query := `
	UPDATE users
	SET
		password_hash = $1,
		updated_at = $2
	WHERE id = $3
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		newPasswordHash,
		time.Now(),
		userID,
	)

	return err
}
