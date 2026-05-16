package models

import "time"

type Certification struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`

	Title string `json:"title"`

	IssuingOrganization string `json:"issuing_organization"`

	CredentialID string `json:"credential_id"`

	CredentialURL string `json:"credential_url"`

	IssueDate  *string `json:"issue_date"`
	ExpiryDate *string `json:"expiry_date"`

	IsVerified bool `json:"is_verified"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
