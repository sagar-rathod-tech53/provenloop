package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sagar-rathod-tech53/provenloop/config"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/internal/repositories"
	"github.com/sagar-rathod-tech53/provenloop/utils"
	"golang.org/x/crypto/bcrypt"
)

var otpLength = 6 // Global variable for OTP length

type AuthService struct {
	DB                  *sql.DB
	UserRepository      repositories.UserRepository
	OTPRepository       repositories.OTPRepository
	TokenExpiration     time.Duration
	OTPLifespan         time.Duration
	BlacklistRepository repositories.TokenBlacklistRepository
	// Config              config.Config
}

// RegisterUserWithOTP handles the registration of a new user and sends an OTP.
func (s *AuthService) RegisterUserWithOTP(ctx context.Context, email, username, password string) error {
	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user model
	user := models.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save the user in the database
	err = s.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Generate OTP
	otp := utils.GenerateOTP(otpLength)

	// Create OTP record
	otpRecord := models.OTP{
		Email:      email,
		OTP:        otp,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}

	// Save OTP to the database
	err = s.OTPRepository.SaveOTP(ctx, otpRecord)
	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	// Send OTP email to the user
	subject := "Verify Your Email Address with Do Host Network"
	body := fmt.Sprintf(`Hello,

Thank you for choosing Do Host Network. Please use the One-Time Password (OTP) below to verify your email address.

Your OTP is: %s

This OTP is valid for the next 10 minutes. Please keep this code confidential and do not share it with anyone.

If you did not request this email, please contact our support team or ignore this message.

Best regards,
The Do Host Network Team`, otp)

	err = utils.SendEmail(email, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	fmt.Println("Registration and OTP email sent successfully to", email)
	return nil
}

func (s *AuthService) LoginUser(ctx context.Context, identifier, password string) (string, string, error) {
	// Step 1: Check if OTP is verified
	otpRecord, err := s.OTPRepository.GetOTPByEmail(ctx, identifier)
	if err != nil {
		// Not fatal: if identifier is a username, it won't have OTP verification
		otpRecord = nil
	}
	if otpRecord != nil && !otpRecord.IsVerified {
		return "", "", errors.New("email not verified. Please verify your email before logging in")
	}

	// Step 2: Fetch user details from repository using email or username
	user, err := s.UserRepository.GetUserByEmailOrUsername(identifier)
	if err != nil {
		return "", "", errors.New("invalid email/username or password")
	}

	// Step 3: Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid email/username or password")
	}

	// Step 4: Load config and generate JWT
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return "", "", fmt.Errorf("failed to load config: %w", err)
	}

	token, err := utils.GenerateToken(24*time.Hour, user.ID, cfg.TokenSecret)
	if err != nil {
		return "", "", err
	}

	// Return both token and user ID
	return token, user.ID, nil
}

// LogoutUser handles user logout by invalidating the JWT token
// func (s *AuthService) LogoutUser(ctx *gin.Context) {
// 	// Invalidate the token (e.g., clear cookie or blacklist the token)
// 	// Clear the token cookie
// 	// Load config and generate JWT
// 	cfg, err := config.LoadConfig(".")
// 	if err != nil {
// 		return  fmt.Errorf("failed to load config: %w", err)
// 	}
// 	ctx.SetCookie("token", "", -1, "/", cfg.COOKIEDOMAIN, false, true)

// 	// Optionally, perform any other cleanup related to the token or session.

// 	// No return needed if the action is successful
// }

// VerifyOTP handles OTP verification
func (s *AuthService) VerifyOTP(ctx context.Context, email, otp string) error {
	storedOTP, err := s.OTPRepository.GetOTPByEmail(ctx, email)
	if err != nil {
		return errors.New("OTP not found or expired")
	}

	if storedOTP.OTP != otp {
		return errors.New("invalid OTP")
	}

	return s.OTPRepository.MarkUserVerified(ctx, email)
}

// ForgotPassword generates an OTP for password reset and sends an email
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// Step 1: Check if user exists
	user, err := s.UserRepository.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "user not found" {
			return fmt.Errorf("no user registered with this email")
		}
		return fmt.Errorf("failed to check user: %w", err)
	}

	// Step 2: Generate OTP
	otp := utils.GenerateOTP(6) // Generate a 6-digit OTP

	// Step 3: Create OTP record
	otpRecord := models.OTP{
		Email:      user.Email,
		OTP:        otp,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}

	// Step 4: Save OTP to the database
	err = s.OTPRepository.SaveOTP(ctx, otpRecord)
	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	// Step 5: Send email
	subject := "Reset Your Password - Do Host Network"
	body := fmt.Sprintf(`Hello %s,

We received a request to reset the password for your Do Host Network account. Please use the One-Time Password (OTP) below to reset your password.

Your OTP is: %s

This OTP is valid for the next 10 minutes. If you did not request this password reset, please contact our support team or ignore this email.

Best regards,
The Do Host Network Team`, user.Username, otp)

	err = utils.SendEmail(user.Email, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	fmt.Println("Password reset email sent successfully to", user.Email)
	return nil
}

// ResetPassword resets the user's password
func (s *AuthService) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	if err := s.VerifyOTP(ctx, email, otp); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.UserRepository.UpdatePassword(ctx, email, hashedPassword)
}
