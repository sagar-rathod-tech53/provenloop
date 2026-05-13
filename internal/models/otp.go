package models

import "time"

type OTP struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	OTP        string    `json:"otp"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}
