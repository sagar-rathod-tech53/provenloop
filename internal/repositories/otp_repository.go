package repositories

import (
	"context"
	"database/sql"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
)

type OTPRepository struct {
	DB *sql.DB
}

// InsertOTP saves an OTP for the user email into the database
func (r *OTPRepository) InsertOTP(ctx context.Context, email, otp string) error {
	query := `INSERT INTO otps (email, otp) VALUES ($1, $2)`
	_, err := r.DB.ExecContext(ctx, query, email, otp)
	return err
}

// SaveOTP saves an OTP record in the database
func (r *OTPRepository) SaveOTP(ctx context.Context, otp models.OTP) error {
	query := `INSERT INTO otps (id, email, otp, is_verified, created_at) 
	          VALUES (uuid_generate_v4(), $1, $2, $3, $4)`
	_, err := r.DB.ExecContext(ctx, query, otp.Email, otp.OTP, otp.IsVerified, otp.CreatedAt)
	return err
}

// GetOTPByEmail retrieves the OTP record by email
func (r *OTPRepository) GetOTPByEmail(ctx context.Context, email string) (*models.OTP, error) {
	query := `SELECT id, email, otp, is_verified, created_at FROM otps WHERE email = $1 ORDER BY created_at DESC LIMIT 1`
	row := r.DB.QueryRowContext(ctx, query, email)

	var otp models.OTP
	err := row.Scan(&otp.ID, &otp.Email, &otp.OTP, &otp.IsVerified, &otp.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &otp, nil
}

// MarkUserVerified updates a user's verification status
func (r *OTPRepository) MarkUserVerified(ctx context.Context, email string) error {
	query := `UPDATE otps SET is_verified = TRUE WHERE email = $1`
	_, err := r.DB.ExecContext(ctx, query, email)
	return err
}
