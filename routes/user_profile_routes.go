package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/middlewares"
)

func RegisterUserProfileRoutes(
	router *gin.Engine,
	profileController *controllers.UserProfileController,
	db *sql.DB,
) {

	profile := router.Group("/api/v1/profile")

	// =========================
	// PROTECTED ROUTES
	// =========================
	profile.Use(middlewares.DeserializeUser(db))

	{
		profile.POST("/create", profileController.CreateProfile)
		profile.GET("/get/:user_id", profileController.GetProfile)
		profile.PUT("/update", profileController.UpdateProfile)
		profile.DELETE("/delete", profileController.DeleteProfile)
	}
}
