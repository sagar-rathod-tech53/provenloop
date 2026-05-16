package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/internal/services"
)

type UserProfileController struct {
	Service *services.UserProfileService
}

// CREATE PROFILE
func (c *UserProfileController) CreateProfile(ctx *gin.Context) {

	var payload models.UserProfile

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// from auth middleware
	userID, _ := ctx.Get("user_id")
	payload.UserID = userID.(string)

	err := c.Service.CreateProfile(ctx, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "profile created"})
}

// GET PROFILE
func (c *UserProfileController) GetProfile(ctx *gin.Context) {

	userID := ctx.Param("user_id")

	profile, err := c.Service.GetProfile(ctx, userID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data":    profile,
	})
}

// UPDATE PROFILE
func (c *UserProfileController) UpdateProfile(ctx *gin.Context) {

	userID, _ := ctx.Get("user_id")

	var payload models.UserProfile

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := c.Service.UpdateProfile(ctx, userID.(string), payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
	})
}

// DELETE PROFILE
func (c *UserProfileController) DeleteProfile(ctx *gin.Context) {

	userID, _ := ctx.Get("user_id")

	err := c.Service.DeleteProfile(ctx, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "profile deleted successfully",
	})
}
