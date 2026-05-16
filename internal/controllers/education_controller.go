package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/internal/services"
)

type EducationController struct {
	Service *services.EducationService
}

// CREATE
func (c *EducationController) Create(ctx *gin.Context) {

	userID, _ := ctx.Get("user_id")

	var payload models.Education
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	payload.UserID = userID.(string)

	err := c.Service.Create(ctx, payload)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "education added"})
}

// GET ALL
func (c *EducationController) GetAll(ctx *gin.Context) {

	userID := ctx.Param("user_id")

	data, err := c.Service.GetAll(ctx, userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": data})
}

// UPDATE
func (c *EducationController) Update(ctx *gin.Context) {

	userID, _ := ctx.Get("user_id")

	var payload models.Education
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	payload.UserID = userID.(string)

	err := c.Service.Update(ctx, payload)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "updated"})
}

// DELETE
func (c *EducationController) Delete(ctx *gin.Context) {

	userID, _ := ctx.Get("user_id")
	id := ctx.Param("id")

	err := c.Service.Delete(ctx, id, userID.(string))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "deleted"})
}
