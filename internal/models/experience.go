package models

import "time"

type Experience struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`

	CompanyName    string `json:"company_name"`
	Role           string `json:"role"`
	EmploymentType string `json:"employment_type"`
	Location       string `json:"location"`

	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`

	IsCurrent  *bool `json:"is_current"`
	IsVerified bool  `json:"is_verified"`

	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
