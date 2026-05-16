package models

import "time"

type Education struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`

	InstitutionName string `json:"institution_name"`
	Degree          string `json:"degree"`
	FieldOfStudy    string `json:"field_of_study"`

	StartYear int `json:"start_year"`
	EndYear   int `json:"end_year"`

	Grade       string `json:"grade"`
	Description string `json:"description"`

	IsVerified bool `json:"is_verified"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
