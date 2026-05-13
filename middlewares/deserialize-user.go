package middlewares

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sagar-rathod-tech53/provenloop/config"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/utils"
)

// DeserializeUser is a middleware to validate and fetch the user from the database based on the provided access token
func DeserializeUser(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		// 1. Try Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Fallback to cookie
		if token == "" {
			cookie, err := ctx.Cookie("token")
			if err == nil {
				token = cookie
			}
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		// Validate token
		config, _ := config.LoadConfig(".")
		sub, err := utils.ValidateToken(token, config.TokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		// Fetch user
		var user models.User
		query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
		row := db.QueryRow(query, sub)

		err = row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user data"})
			}
			ctx.Abort()
			return
		}

		// Attach user to context using "user" key
		ctx.Set("user", user)
		ctx.Next()
	}
}
