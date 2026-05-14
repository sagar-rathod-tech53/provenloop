package middlewares

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sagar-rathod-tech53/provenloop/config"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/utils"
)

func DeserializeUser(
	db *sql.DB,
) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		var token string

		authHeader := ctx.GetHeader(
			"Authorization",
		)

		if strings.HasPrefix(
			authHeader,
			"Bearer ",
		) {

			token = strings.TrimPrefix(
				authHeader,
				"Bearer ",
			)
		}

		if token == "" {

			ctx.AbortWithStatusJSON(
				401,
				gin.H{
					"status": false,
					"error":  "authorization token missing",
				},
			)

			return
		}

		// blacklist check

		blacklisted, _ := config.RDB.Get(
			ctx,
			"blacklist:"+token,
		).Result()

		if blacklisted == "true" {

			ctx.AbortWithStatusJSON(
				401,
				gin.H{
					"status": false,
					"error":  "user logged out",
				},
			)

			return
		}

		sub, err := utils.ValidateToken(
			token,
			config.AppConfig.TokenSecret,
		)

		if err != nil {

			ctx.AbortWithStatusJSON(
				401,
				gin.H{
					"status": false,
					"error":  "invalid token",
				},
			)

			return
		}

		userID := sub.(string)

		var user models.User

		query := `
		SELECT
			id,
			email,
			username,
			password_hash,
			is_verified,
			created_at,
			updated_at
		FROM users
		WHERE id=$1
		`

		err = db.QueryRowContext(
			ctx,
			query,
			userID,
		).Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.PasswordHash,
			&user.IsVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {

			ctx.AbortWithStatusJSON(
				401,
				gin.H{
					"status": false,
					"error":  "user not found",
				},
			)

			return
		}

		ctx.Set(
			"user",
			user,
		)

		ctx.Set(
			"user_id",
			user.ID,
		)

		ctx.Next()

	}
}
