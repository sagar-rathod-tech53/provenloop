package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/middlewares"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	authController *controllers.AuthController,
	db *sql.DB,
) {

	auth := router.Group("/auth")

	// =========================
	// PUBLIC ROUTES
	// =========================
	{
		auth.POST("/register", authController.Register)
		auth.POST("/verify-otp", authController.VerifyOTP)
		auth.POST("/resend-registration-otp", authController.ResendRegistrationOTP)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh-token", authController.RefreshToken)
		auth.POST("/forgot-password", authController.ForgotPassword)
		auth.POST("/reset-password", authController.ResetPassword)
	}

	// =========================
	// PROTECTED ROUTES
	// =========================
	protected := auth.Group("/")
	protected.Use(middlewares.DeserializeUser(db))

	{
		protected.POST("/logout", authController.Logout)
		protected.POST("/change-password", authController.ChangePassword)
	}
}
