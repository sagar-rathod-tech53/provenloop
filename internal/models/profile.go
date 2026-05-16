package models

import "time"

type UserProfile struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`

	Username string `json:"username"`
	Email    string `json:"email"`

	FullName            string  `json:"full_name"`
	Designation         *string `json:"designation"`
	Organization        *string `json:"organization"`
	ProfessionalSummary *string `json:"professional_summary"`
	Location            *string `json:"location"`

	ProfileImage *string `json:"profile_image"`
	CoverImage   *string `json:"cover_image"`

	ContactNumber *string `json:"contact_number"`

	CollegeName    *string `json:"college_name"`
	UniversityName *string `json:"university_name"`
	Degree         *string `json:"degree"`
	FieldOfStudy   *string `json:"field_of_study"`
	GraduationYear *int    `json:"graduation_year"`

	ProfileVideoURL *string `json:"profile_video_url"`

	LastActiveAt time.Time `json:"last_active_at"`
	IsVerified   bool      `json:"is_verified"`
	IsPublic     bool      `json:"is_public"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
