package models

import "time"

type Skill struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`

	SkillName string `json:"skill_name"`

	ProficiencyLevel string `json:"proficiency_level"`

	EndorsementsCount int `json:"endorsements_count"`

	IsVerified bool `json:"is_verified"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
