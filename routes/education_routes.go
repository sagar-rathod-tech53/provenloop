package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/middlewares"
)

func RegisterEducationRoutes(
	router *gin.Engine,
	educationController *controllers.EducationController,
	db *sql.DB,
) {

	education := router.Group("/api/v1/education")
	education.Use(middlewares.DeserializeUser(db))

	{
		education.POST("/create", educationController.Create)
		education.GET("/:user_id", educationController.GetAll)
		education.PUT("/update", educationController.Update)
		education.DELETE("/delete/:id", educationController.Delete)
	}
}
