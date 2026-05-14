package routes

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/sagar-rathod-tech53/provenloop/config"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/internal/repositories"
	"github.com/sagar-rathod-tech53/provenloop/internal/services"
	"github.com/sagar-rathod-tech53/provenloop/middlewares"
	"github.com/sagar-rathod-tech53/provenloop/migrations"
)

func SetupServer() *gin.Engine {

	// ==========================================
	// Load Config
	// ==========================================

	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// ==========================================
	// Connect PostgreSQL
	// ==========================================

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// ==========================================
	// Run Migrations
	// ==========================================

	err = migrations.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	// ==========================================
	// Connect Redis
	// ==========================================

	config.ConnectRedis()

	// ==========================================
	// Repositories
	// ==========================================

	userRepo := &repositories.UserRepository{
		DB: db,
	}

	// ==========================================
	// Services
	// ==========================================

	authService := &services.AuthService{
		UserRepo: userRepo,
	}

	// ==========================================
	// Controllers
	// ==========================================

	authController := &controllers.AuthController{
		AuthService: authService,
	}

	// ==========================================
	// Router
	// ==========================================

	router := gin.Default()

	// ==========================================
	// Public Routes
	// ==========================================

	auth := router.Group("/auth")

	{
		auth.POST(
			"/register",
			authController.Register,
		)

		auth.POST(
			"/verify-otp",
			authController.VerifyOTP,
		)

		auth.POST(
			"/resend-registration-otp",
			authController.ResendRegistrationOTP,
		)

		auth.POST(
			"/login",
			authController.Login,
		)

		auth.POST(
			"/refresh-token",
			authController.RefreshToken,
		)

		auth.POST(
			"/forgot-password",
			authController.ForgotPassword,
		)

		auth.POST(
			"/reset-password",
			authController.ResetPassword,
		)
	}

	// ==========================================
	// Protected Routes
	// ==========================================

	protected := router.Group("/auth")

	protected.Use(
		middlewares.DeserializeUser(db),
	)

	{
		protected.POST(
			"/logout",
			authController.Logout,
		)

		protected.POST(
			"/change-password",
			authController.ChangePassword,
		)
	}

	return router
}
