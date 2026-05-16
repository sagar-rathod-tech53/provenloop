package models

import "time"

type Project struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`

	Title       string `json:"title"`
	Description string `json:"description"`

	TechStack []string `json:"tech_stack"`

	GithubURL string `json:"github_url"`
	LiveURL   string `json:"live_url"`

	ProjectVideoURL string `json:"project_video_url"`

	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`

	IsVerified bool `json:"is_verified"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
