package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/internal/services"
)

type AuthController struct {
	AuthService *services.AuthService
}

func (c *AuthController) Register(
	ctx *gin.Context,
) {

	var payload struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(
		&payload,
	); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err := c.AuthService.RegisterUser(
		ctx,
		payload.Email,
		payload.Username,
		payload.Password,
	)

	if err != nil {

		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"message": "user registered successfully",
		},
	)
}

func (c *AuthController) VerifyOTP(
	ctx *gin.Context,
) {

	var payload struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	// validate request
	if err := ctx.ShouldBindJSON(
		&payload,
	); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	// call service
	err := c.AuthService.VerifyOTP(
		ctx,
		payload.Email,
		payload.OTP,
	)

	if err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "otp verified successfully",
		},
	)
}

func (c *AuthController) ResendRegistrationOTP(
	ctx *gin.Context,
) {

	var body struct {
		Email string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {

		ctx.JSON(
			400,
			gin.H{
				"error": "invalid request body",
			},
		)
		return
	}

	err := c.AuthService.ResendRegistrationOTP(
		ctx,
		body.Email,
	)

	if err != nil {

		ctx.JSON(
			400,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	ctx.JSON(
		200,
		gin.H{
			"message": "otp resent successfully",
		},
	)
}

func (c *AuthController) Login(
	ctx *gin.Context,
) {

	var body struct {
		EmailOrUsername string `json:"email_or_username"`
		Password        string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			models.LoginResponse{
				Status: false,
				Error:  "invalid request body",
			},
		)
		return
	}

	response := c.AuthService.LoginUser(
		ctx,
		body.EmailOrUsername,
		body.Password,
	)

	if !response.Status {

		ctx.JSON(
			http.StatusUnauthorized,
			response,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}

func (c *AuthController) RefreshToken(
	ctx *gin.Context,
) {

	var body models.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			models.RefreshTokenResponse{
				Status: false,
				Error:  "invalid request body",
			},
		)
		return
	}

	response := c.AuthService.RefreshToken(
		ctx,
		body.RefreshToken,
	)

	if !response.Status {

		ctx.JSON(
			http.StatusUnauthorized,
			response,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}

func (c *AuthController) Logout(
	ctx *gin.Context,
) {
	var body models.LogoutRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			models.LogoutResponse{
				Status: false,
				Error:  "invalid request body",
			},
		)
		return
	}

	authHeader := ctx.GetHeader("Authorization")

	token := strings.TrimPrefix(
		authHeader,
		"Bearer ",
	)

	response := c.AuthService.LogoutUser(
		ctx,
		token,
		body.RefreshToken,
	)

	if !response.Status {
		ctx.JSON(
			http.StatusUnauthorized,
			response,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"status":  true,
			"message": "logout successful",
		},
	)
}

func (c *AuthController) ForgotPassword(
	ctx *gin.Context,
) {

	var body models.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			models.ForgotPasswordResponse{
				Status: false,
				Error:  "invalid request body",
			},
		)
		return
	}

	response := c.AuthService.ForgotPassword(
		ctx,
		body.Email,
	)

	if !response.Status {

		ctx.JSON(
			http.StatusBadRequest,
			response,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"status":  true,
			"message": "forgot password otp sent successfully",
		},
	)
}

func (c *AuthController) ResetPassword(
	ctx *gin.Context,
) {

	var body models.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			models.ResetPasswordResponse{
				Status: false,
				Error:  "invalid request body",
			},
		)
		return
	}

	response := c.AuthService.ResetPassword(
		ctx,
		body.Email,
		body.OTP,
		body.NewPassword,
	)

	if !response.Status {

		ctx.JSON(
			http.StatusBadRequest,
			response,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"status":  true,
			"message": "password reset successful",
		},
	)
}

func (c *AuthController) ChangePassword(
	ctx *gin.Context,
) {

	var body models.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {

		ctx.JSON(
			http.StatusBadRequest,
			models.ChangePasswordResponse{
				Status: false,
				Error:  "invalid request body",
			},
		)
		return
	}

	// =====================================
	// Get User From Middleware
	// =====================================

	userID := ctx.GetString(
		"user_id",
	)

	if userID == "" {

		ctx.JSON(
			http.StatusUnauthorized,
			models.ChangePasswordResponse{
				Status: false,
				Error:  "unauthorized",
			},
		)
		return
	}

	response := c.AuthService.ChangePassword(
		ctx,
		userID,
		body.OldPassword,
		body.NewPassword,
	)

	if !response.Status {

		ctx.JSON(
			http.StatusBadRequest,
			response,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"status":  true,
			"message": "password changed successfully",
		},
	)
}
