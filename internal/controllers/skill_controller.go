package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/internal/services"
)

type SkillController struct {
	Service *services.SkillService
}

func (c *SkillController) Create(
	ctx *gin.Context,
) {

	userID, _ := ctx.Get(
		"user_id",
	)

	var payload models.Skill

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
			500,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	ctx.JSON(
		201,
		gin.H{
			"message": "skill created",
		},
	)
}

func (c *SkillController) GetAll(
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

func (c *SkillController) Update(
	ctx *gin.Context,
) {

	userID, _ := ctx.Get(
		"user_id",
	)

	var payload models.Skill

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
			"message": "updated",
		},
	)
}

func (c *SkillController) Delete(
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
			"message": "deleted",
		},
	)
}
