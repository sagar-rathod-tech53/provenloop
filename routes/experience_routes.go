package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/middlewares"
)

func ExperienceRoutes(
	router *gin.Engine,
	expController *controllers.ExperienceController,
	db *sql.DB,
) {

	exp := router.Group("/api/v1/experience")
	exp.Use(middlewares.DeserializeUser(db))

	{
		exp.POST("/create", expController.Create)
		exp.GET("/:user_id", expController.GetAll)
		exp.PUT("/update", expController.Update)
		exp.DELETE("/delete/:id", expController.Delete)
	}
}
