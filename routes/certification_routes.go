package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/controllers"
	"github.com/sagar-rathod-tech53/provenloop/middlewares"
)

func RegisterCertificationRoutes(
	router *gin.Engine,
	certificationController *controllers.CertificationController,
	db *sql.DB,
) {

	certifications := router.Group(
		"/api/v1/certifications",
	)

	certifications.Use(
		middlewares.DeserializeUser(db),
	)

	{
		certifications.POST("/create", certificationController.Create)
		certifications.GET("/:user_id", certificationController.GetAll)
		certifications.PUT("/update", certificationController.Update)
		certifications.DELETE("/delete/:id", certificationController.Delete)
	}
}
