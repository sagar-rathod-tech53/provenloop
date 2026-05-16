package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/middlewares"
)

func SkillRoutes(
	router *gin.Engine,
	skillController *controllers.SkillController,
	db *sql.DB,
) {

	skills := router.Group(
		"/api/v1/skills",
	)

	skills.Use(
		middlewares.DeserializeUser(db),
	)

	{
		skills.POST(
			"/create",
			skillController.Create,
		)

		skills.GET(
			"/:user_id",
			skillController.GetAll,
		)

		skills.PUT(
			"/update",
			skillController.Update,
		)

		skills.DELETE(
			"/delete/:id",
			skillController.Delete,
		)
	}
}
