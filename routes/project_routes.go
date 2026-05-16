package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/middlewares"
)

func ProjectRoutes(
	router *gin.Engine,
	projectController *controllers.ProjectController,
	db *sql.DB,
) {

	project := router.Group(
		"/api/v1/projects",
	)

	project.Use(
		middlewares.DeserializeUser(db),
	)

	{
		project.POST(
			"/create",
			projectController.Create,
		)

		project.GET(
			"/:user_id",
			projectController.GetAll,
		)

		project.PUT(
			"/update",
			projectController.Update,
		)

		project.DELETE(
			"/delete/:id",
			projectController.Delete,
		)
	}
}
