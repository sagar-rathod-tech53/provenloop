package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sagar-rathod-tech53/provenloop/config"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/internal/repositories"
	"github.com/sagar-rathod-tech53/provenloop/utils"
)

type AuthService struct {
	UserRepo *repositories.UserRepository
}

func (s *AuthService) RegisterUser(
	ctx context.Context,
	email string,
	username string,
	password string,
) error {

	ctx, cancel := context.WithTimeout(
		ctx,
		10*time.Second,
	)
	defer cancel()

	// =====================================
	// Check Existing User
	// =====================================

	exists, err := s.UserRepo.CheckUserExists(
		ctx,
		email,
		username,
	)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf(
			"email or username already exists",
		)
	}

	// =====================================
	// Concurrent Tasks
	// =====================================

	hashChan := make(chan string, 1)
	hashErrChan := make(chan error, 1)

	otpChan := make(chan string, 1)

	// password hashing
	go func() {

		hashedPassword, err := utils.HashPassword(
			password,
		)

		if err != nil {
			hashErrChan <- err
			return
		}

		hashChan <- hashedPassword
	}()

	// otp generation
	go func() {

		otp := utils.GenerateOTP(6)

		otpChan <- otp
	}()

	// wait results
	var hashedPassword string
	var otp string

	for i := 0; i < 2; i++ {

		select {

		case err := <-hashErrChan:
			return err

		case hash := <-hashChan:
			hashedPassword = hash

		case generatedOTP := <-otpChan:
			otp = generatedOTP
		}
	}

	// =====================================
	// Create User
	// =====================================

	user := models.User{
		ID:           uuid.New().String(),
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.UserRepo.CreateUser(
		ctx,
		user,
	)

	if err != nil {
		return err
	}

	// =====================================
	// Save OTP in Redis
	// =====================================

	err = config.RDB.Set(
		ctx,
		"otp:"+email,
		otp,
		5*time.Minute,
	).Err()

	if err != nil {
		return err
	}

	// =====================================
	// Send Email Async
	// =====================================

	go func() {

		subject := "Email Verification OTP"

		body := fmt.Sprintf(`
<h1>Your OTP is: %s</h1>
`, otp)

		_ = utils.SendEmail(
			email,
			subject,
			body,
		)

	}()

	return nil
}

func (s *AuthService) VerifyOTP(
	ctx context.Context,
	email string,
	otp string,
) error {

	// get otp from redis
	storedOTP, err := config.RDB.Get(
		ctx,
		"otp:"+email,
	).Result()

	if err != nil {
		return fmt.Errorf("otp expired or not found")
	}

	// compare otp
	if storedOTP != otp {
		return fmt.Errorf("invalid otp")
	}

	// delete otp after verification
	err = config.RDB.Del(
		ctx,
		"otp:"+email,
	).Err()

	if err != nil {
		return err
	}

	// mark verified in database
	err = s.UserRepo.VerifyUser(
		ctx,
		email,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResendRegistrationOTP(
	ctx context.Context,
	email string,
) error {

	ctx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	// =========================================
	// Check User Exists
	// =========================================

	exists, err := s.UserRepo.CheckUserExists(
		ctx,
		email,
		"",
	)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf(
			"user not found",
		)
	}

	// =========================================
	// Generate OTP
	// =========================================

	otp := utils.GenerateOTP(6)

	// =========================================
	// Save OTP In Redis
	// =========================================

	err = config.RDB.Set(
		ctx,
		"otp:"+email,
		otp,
		10*time.Minute,
	).Err()

	if err != nil {
		return err
	}

	// =========================================
	// Send Email In Background
	// =========================================

	go func() {

		subject := "Resend Registration OTP"

		body := fmt.Sprintf(`
<h2>Email Verification OTP</h2>

<h1>%s</h1>

<p>This OTP is valid for 10 minutes.</p>
`, otp)

		err := utils.SendEmail(
			email,
			subject,
			body,
		)

		if err != nil {
			fmt.Println(
				"failed to send resend otp email:",
				err,
			)
		}
	}()

	return nil
}

func (s *AuthService) LoginUser(
	ctx context.Context,
	emailOrUsername string,
	password string,
) models.LoginResponse {

	ctx, cancel := context.WithTimeout(
		ctx,
		10*time.Second,
	)
	defer cancel()

	// =====================================
	// Get User
	// =====================================

	user, err := s.UserRepo.GetUserByEmailOrUsername(
		ctx,
		emailOrUsername,
	)

	if err != nil {

		return models.LoginResponse{
			Status: false,
			Error:  "invalid credentials",
		}
	}

	// =====================================
	// Check Verification
	// =====================================

	if !user.IsVerified {

		return models.LoginResponse{
			Status: false,
			Error:  "please verify your account",
		}
	}

	// =====================================
	// Channels
	// =====================================

	passwordChan := make(chan error, 1)

	accessTokenChan := make(chan string, 1)
	refreshTokenChan := make(chan string, 1)

	tokenErrChan := make(chan error, 2)

	// =====================================
	// Concurrent Password Verification
	// =====================================

	go func() {

		err := utils.VerifyPassword(
			user.PasswordHash,
			password,
		)

		passwordChan <- err
	}()

	// =====================================
	// Concurrent Access Token Generation
	// =====================================

	go func() {

		token, err := utils.GenerateAccessToken(
			user.ID,
		)

		if err != nil {
			tokenErrChan <- err
			return
		}

		accessTokenChan <- token
	}()

	// =====================================
	// Concurrent Refresh Token Generation
	// =====================================

	go func() {

		token, err := utils.GenerateRefreshToken(
			user.ID,
		)

		if err != nil {
			tokenErrChan <- err
			return
		}

		refreshTokenChan <- token
	}()

	var accessToken string
	var refreshToken string

	// =====================================
	// Wait Concurrent Results
	// =====================================

	for i := 0; i < 3; i++ {

		select {

		case err := <-passwordChan:

			if err != nil {

				return models.LoginResponse{
					Status: false,
					Error:  "invalid credentials",
				}
			}

		case token := <-accessTokenChan:

			accessToken = token

		case token := <-refreshTokenChan:

			refreshToken = token

		case err := <-tokenErrChan:

			return models.LoginResponse{
				Status: false,
				Error:  err.Error(),
			}
		}
	}

	// =====================================
	// Store Session In Redis Concurrently
	// =====================================

	go func() {

		config.RDB.Set(
			context.Background(),
			"session:"+user.ID,
			accessToken,
			24*time.Hour,
		)

		config.RDB.Set(
			context.Background(),
			"refresh:"+user.ID,
			refreshToken,
			7*24*time.Hour,
		)
	}()

	return models.LoginResponse{
		Status:       true,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) models.RefreshTokenResponse {

	ctx, cancel := context.WithTimeout(
		ctx,
		10*time.Second,
	)
	defer cancel()

	// =====================================
	// Validate Refresh Token
	// =====================================

	userID, err := utils.ValidateRefreshToken(
		refreshToken,
	)

	if err != nil {

		return models.RefreshTokenResponse{
			Status: false,
			Error:  "invalid refresh token",
		}
	}

	// =====================================
	// Check Redis Session
	// =====================================

	storedRefreshToken, err := config.RDB.Get(
		ctx,
		"refresh:"+userID,
	).Result()

	if err != nil {

		return models.RefreshTokenResponse{
			Status: false,
			Error:  "session expired",
		}
	}

	if storedRefreshToken != refreshToken {

		return models.RefreshTokenResponse{
			Status: false,
			Error:  "refresh token mismatch",
		}
	}

	// =====================================
	// Concurrent Token Generation
	// =====================================

	accessTokenChan := make(chan string, 1)
	refreshTokenChan := make(chan string, 1)
	errChan := make(chan error, 2)

	// generate new access token
	go func() {

		token, err := utils.GenerateAccessToken(
			userID,
		)

		if err != nil {
			errChan <- err
			return
		}

		accessTokenChan <- token
	}()

	// generate new refresh token
	go func() {

		token, err := utils.GenerateRefreshToken(
			userID,
		)

		if err != nil {
			errChan <- err
			return
		}

		refreshTokenChan <- token
	}()

	var newAccessToken string
	var newRefreshToken string

	// =====================================
	// Wait Results
	// =====================================

	for i := 0; i < 2; i++ {

		select {

		case token := <-accessTokenChan:
			newAccessToken = token

		case token := <-refreshTokenChan:
			newRefreshToken = token

		case err := <-errChan:

			return models.RefreshTokenResponse{
				Status: false,
				Error:  err.Error(),
			}
		}
	}

	// =====================================
	// Save New Tokens In Redis Concurrently
	// =====================================

	saveErrChan := make(chan error, 2)

	// save access token
	go func() {

		err := config.RDB.Set(
			ctx,
			"session:"+userID,
			newAccessToken,
			24*time.Hour,
		).Err()

		saveErrChan <- err
	}()

	// save refresh token
	go func() {

		err := config.RDB.Set(
			ctx,
			"refresh:"+userID,
			newRefreshToken,
			7*24*time.Hour,
		).Err()

		saveErrChan <- err
	}()

	for i := 0; i < 2; i++ {

		if err := <-saveErrChan; err != nil {

			return models.RefreshTokenResponse{
				Status: false,
				Error: fmt.Sprintf(
					"redis save failed: %v",
					err,
				),
			}
		}
	}

	return models.RefreshTokenResponse{
		Status:       true,
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
	}
}

func (s *AuthService) LogoutUser(
	ctx context.Context,
	accessToken string,
	refreshToken string,
) models.LogoutResponse {

	ctx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	userID, err := utils.ValidateRefreshToken(
		refreshToken,
	)

	if err != nil {
		return models.LogoutResponse{
			Status: false,
			Error:  "invalid refresh token",
		}
	}

	storedToken, err := config.RDB.Get(
		ctx,
		"refresh:"+userID,
	).Result()

	if err != nil {

		return models.LogoutResponse{
			Status: false,
			Error:  "user already logged out",
		}
	}

	if storedToken != refreshToken {

		return models.LogoutResponse{
			Status: false,
			Error:  "session mismatch",
		}
	}

	errChan := make(chan error, 3)

	// remove session
	go func() {
		errChan <- config.RDB.Del(
			ctx,
			"session:"+userID,
		).Err()
	}()

	// remove refresh
	go func() {
		errChan <- config.RDB.Del(
			ctx,
			"refresh:"+userID,
		).Err()
	}()

	// blacklist access token
	go func() {
		errChan <- config.RDB.Set(
			ctx,
			"blacklist:"+accessToken,
			"true",
			15*time.Minute,
		).Err()
	}()

	for i := 0; i < 3; i++ {

		if err := <-errChan; err != nil {

			return models.LogoutResponse{
				Status: false,
				Error:  err.Error(),
			}
		}
	}

	return models.LogoutResponse{
		Status: true,
	}
}

func (s *AuthService) ForgotPassword(
	ctx context.Context,
	email string,
) models.ForgotPasswordResponse {

	ctx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	// =====================================
	// Check User Exists
	// =====================================

	exists, err := s.UserRepo.CheckUserExistsByEmail(
		ctx,
		email,
	)

	if err != nil {

		return models.ForgotPasswordResponse{
			Status: false,
			Error:  err.Error(),
		}
	}

	if !exists {

		return models.ForgotPasswordResponse{
			Status: false,
			Error:  "user not found",
		}
	}

	// =====================================
	// Generate OTP
	// =====================================

	otp := utils.GenerateOTP(6)

	// =====================================
	// Save OTP In Redis
	// =====================================

	err = config.RDB.Set(
		ctx,
		"forgot_otp:"+email,
		otp,
		10*time.Minute,
	).Err()

	if err != nil {

		return models.ForgotPasswordResponse{
			Status: false,
			Error:  err.Error(),
		}
	}

	// =====================================
	// Send Email In Background
	// =====================================

	go func() {

		subject := "Forgot Password OTP"

		body := fmt.Sprintf(`
<h2>Password Reset OTP</h2>

<h1>%s</h1>

<p>This OTP is valid for 10 minutes.</p>
`, otp)

		_ = utils.SendEmail(
			email,
			subject,
			body,
		)

	}()

	// =====================================
	// Immediate Response
	// =====================================

	return models.ForgotPasswordResponse{
		Status: true,
	}
}

func (s *AuthService) ResetPassword(
	ctx context.Context,
	email string,
	otp string,
	newPassword string,
) models.ResetPasswordResponse {

	ctx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	// =====================================
	// Get OTP From Redis
	// =====================================

	storedOTP, err := config.RDB.Get(
		ctx,
		"forgot_otp:"+email,
	).Result()

	if err != nil {

		return models.ResetPasswordResponse{
			Status: false,
			Error:  "otp expired or not found",
		}
	}

	// =====================================
	// Compare OTP
	// =====================================

	if storedOTP != otp {

		return models.ResetPasswordResponse{
			Status: false,
			Error:  "invalid otp",
		}
	}

	// =====================================
	// Concurrent Operations
	// =====================================

	var (
		hashedPassword string
		hashErr        error
	)

	var wg sync.WaitGroup

	wg.Add(1)

	// hash password concurrently
	go func() {
		defer wg.Done()

		hashedPassword, hashErr = utils.HashPassword(
			newPassword,
		)
	}()

	wg.Wait()

	if hashErr != nil {

		return models.ResetPasswordResponse{
			Status: false,
			Error:  hashErr.Error(),
		}
	}

	// =====================================
	// Update Password
	// =====================================

	err = s.UserRepo.UpdatePassword(
		ctx,
		email,
		hashedPassword,
	)

	if err != nil {

		return models.ResetPasswordResponse{
			Status: false,
			Error:  err.Error(),
		}
	}

	// =====================================
	// Delete OTP From Redis
	// =====================================

	go func() {

		config.RDB.Del(
			context.Background(),
			"forgot_otp:"+email,
		)

	}()

	// =====================================
	// Send Email In Background
	// =====================================

	go func() {

		subject := "Password Reset Successful"

		body := `
<h2>Password Reset Successful</h2>

<p>Your password has been changed successfully.</p>
`

		_ = utils.SendEmail(
			email,
			subject,
			body,
		)

	}()

	return models.ResetPasswordResponse{
		Status: true,
	}
}

func (s *AuthService) ChangePassword(
	ctx context.Context,
	userID string,
	oldPassword string,
	newPassword string,
) models.ChangePasswordResponse {

	ctx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	// =====================================
	// Get User
	// =====================================

	user, err := s.UserRepo.GetUserByID(
		ctx,
		userID,
	)

	if err != nil {

		return models.ChangePasswordResponse{
			Status: false,
			Error:  "user not found",
		}
	}

	// =====================================
	// Verify old password
	// =====================================

	err = utils.VerifyPassword(
		user.PasswordHash,
		oldPassword,
	)

	if err != nil {

		return models.ChangePasswordResponse{
			Status: false,
			Error:  "old password incorrect",
		}
	}

	// =====================================
	// prevent same password
	// =====================================

	if oldPassword == newPassword {

		return models.ChangePasswordResponse{
			Status: false,
			Error:  "new password cannot be same as old password",
		}
	}

	// =====================================
	// hash concurrently
	// =====================================

	var (
		hashedPassword string
		hashErr        error
	)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {

		defer wg.Done()

		hashedPassword,
			hashErr = utils.HashPassword(
			newPassword,
		)

	}()

	wg.Wait()

	if hashErr != nil {

		return models.ChangePasswordResponse{
			Status: false,
			Error:  hashErr.Error(),
		}
	}

	// =====================================
	// update db
	// =====================================

	err = s.UserRepo.ChangePassword(
		ctx,
		userID,
		hashedPassword,
	)

	if err != nil {

		return models.ChangePasswordResponse{
			Status: false,
			Error:  err.Error(),
		}
	}

	// =====================================
	// background email
	// =====================================

	go func() {

		subject := "Password Changed"

		body := `
		<h2>Password Updated</h2>

		<p>Your password changed successfully.</p>
		`

		_ = utils.SendEmail(
			user.Email,
			subject,
			body,
		)

	}()

	return models.ChangePasswordResponse{
		Status: true,
	}
}
