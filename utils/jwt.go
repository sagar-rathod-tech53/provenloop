// ============================================
// utils/jwt.go
// ============================================

package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sagar-rathod-tech53/provenloop/config"
)

func GenerateAccessToken(
	userID string,
) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"exp": time.Now().Add(
			15 * time.Minute,
		).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(config.AppConfig.TokenSecret),
	)
}

func GenerateRefreshToken(
	userID string,
) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp": time.Now().Add(
			7 * 24 * time.Hour,
		).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(config.AppConfig.TokenSecret),
	)
}

func ValidateToken(
	tokenString string,
	secret string,
) (interface{}, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims["user_id"], nil
}

func ValidateRefreshToken(
	refreshToken string,
) (string, error) {

	token, err := jwt.Parse(
		refreshToken,
		func(token *jwt.Token) (interface{}, error) {

			return []byte(
				config.AppConfig.TokenSecret,
			), nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	// check token type
	tokenType, ok := claims["type"].(string)

	if !ok || tokenType != "refresh" {
		return "", errors.New(
			"invalid refresh token",
		)
	}

	// get user id
	userID, ok := claims["user_id"].(string)

	if !ok {
		return "", errors.New(
			"invalid user id",
		)
	}

	return userID, nil
}
