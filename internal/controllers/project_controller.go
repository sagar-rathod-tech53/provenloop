package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/internal/services"
)

type ProjectController struct {
	Service *services.ProjectService
}

func (c *ProjectController) Create(
	ctx *gin.Context,
) {

	userID, _ := ctx.Get(
		"user_id",
	)

	var payload models.Project

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

	payload.UserID = userID.(string)

	err := c.Service.Create(
		ctx,
		payload,
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
			"message": "project created",
		},
	)
}

func (c *ProjectController) GetAll(
	ctx *gin.Context,
) {

	userID := ctx.Param(
		"user_id",
	)

	data, err := c.Service.GetAll(
		ctx,
		userID,
	)

	if err != nil {

		ctx.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		200,
		gin.H{
			"data": data,
		},
	)
}

func (c *ProjectController) Update(
	ctx *gin.Context,
) {

	userID, _ := ctx.Get(
		"user_id",
	)

	var payload models.Project

	if err := ctx.ShouldBindJSON(
		&payload,
	); err != nil {

		ctx.JSON(
			400,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	payload.UserID = userID.(string)

	err := c.Service.Update(
		ctx,
		payload,
	)

	if err != nil {

		ctx.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		200,
		gin.H{
			"message": "project updated",
		},
	)
}

func (c *ProjectController) Delete(
	ctx *gin.Context,
) {

	userID, _ := ctx.Get(
		"user_id",
	)

	id := ctx.Param(
		"id",
	)

	err := c.Service.Delete(
		ctx,
		id,
		userID.(string),
	)

	if err != nil {

		ctx.JSON(
			500,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		200,
		gin.H{
			"message": "project deleted",
		},
	)
}
