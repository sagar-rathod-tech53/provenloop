package repositories

import (
	"context"
	"database/sql"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
)

type UserProfileRepository struct {
	DB *sql.DB
}

// CHECK EXISTS
func (r *UserProfileRepository) ProfileExists(ctx context.Context, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_profile WHERE user_id=$1)`

	var exists bool
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&exists)
	return exists, err
}

// CREATE PROFILE
func (r *UserProfileRepository) CreateProfile(ctx context.Context, p models.UserProfile) error {

	query := `
	INSERT INTO user_profile (
		id, user_id,
		profile_image, full_name,
		designation, organization, professional_summary,
		location, contact_number,
		college_name, university_name, degree,
		field_of_study, graduation_year,
		profile_video_url,
		last_active_at, is_verified, is_public,
		created_at, updated_at
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,
		$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20
	)
	`

	_, err := r.DB.ExecContext(ctx, query,
		p.ID,
		p.UserID,
		p.ProfileImage,
		p.FullName,
		p.Designation,
		p.Organization,
		p.ProfessionalSummary,
		p.Location,
		p.ContactNumber,
		p.CollegeName,
		p.UniversityName,
		p.Degree,
		p.FieldOfStudy,
		p.GraduationYear,
		p.ProfileVideoURL,
		p.LastActiveAt,
		p.IsVerified,
		p.IsPublic,
		p.CreatedAt,
		p.UpdatedAt,
	)

	return err
}

// GET PROFILE WITH JOIN (email + username)
func (r *UserProfileRepository) GetProfile(ctx context.Context, userID string) (models.UserProfile, error) {

	query := `
	SELECT 
		p.id,
		p.user_id,

		u.username,
		u.email,

		p.profile_image,
		p.full_name,
		p.designation,
		p.organization,
		p.professional_summary,
		p.location,
		p.contact_number,

		p.college_name,
		p.university_name,
		p.degree,
		p.field_of_study,
		p.graduation_year,

		p.profile_video_url,
		p.last_active_at,
		p.is_verified,
		p.is_public,

		p.created_at,
		p.updated_at

	FROM user_profile p
	JOIN users u ON u.id = p.user_id
	WHERE p.user_id = $1
	`

	var p models.UserProfile

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&p.ID,
		&p.UserID,

		&p.Username,
		&p.Email,

		&p.ProfileImage,
		&p.FullName,
		&p.Designation,
		&p.Organization,
		&p.ProfessionalSummary,
		&p.Location,
		&p.ContactNumber,

		&p.CollegeName,
		&p.UniversityName,
		&p.Degree,
		&p.FieldOfStudy,
		&p.GraduationYear,

		&p.ProfileVideoURL,
		&p.LastActiveAt,
		&p.IsVerified,
		&p.IsPublic,

		&p.CreatedAt,
		&p.UpdatedAt,
	)

	return p, err
}

func (r *UserProfileRepository) UpdateProfile(
	ctx context.Context,
	userID string,
	p models.UserProfile,
) error {

	query := `
	UPDATE user_profile SET
		profile_image = COALESCE($1, profile_image),
		full_name = COALESCE($2, full_name),
		designation = COALESCE($3, designation),
		organization = COALESCE($4, organization),
		professional_summary = COALESCE($5, professional_summary),
		location = COALESCE($6, location),
		contact_number = COALESCE($7, contact_number),

		college_name = COALESCE($8, college_name),
		university_name = COALESCE($9, university_name),
		degree = COALESCE($10, degree),
		field_of_study = COALESCE($11, field_of_study),
		graduation_year = COALESCE($12, graduation_year),

		profile_video_url = COALESCE($13, profile_video_url),

		last_active_at = NOW(),
		updated_at = NOW()
	WHERE user_id = $14
	`

	_, err := r.DB.ExecContext(ctx, query,
		p.ProfileImage,
		p.FullName,
		p.Designation,
		p.Organization,
		p.ProfessionalSummary,
		p.Location,
		p.ContactNumber,
		p.CollegeName,
		p.UniversityName,
		p.Degree,
		p.FieldOfStudy,
		p.GraduationYear,
		p.ProfileVideoURL,
		userID,
	)

	return err
}

// =========================
// DELETE PROFILE
// =========================
func (r *UserProfileRepository) DeleteProfile(
	ctx context.Context,
	userID string,
) error {

	query := `DELETE FROM user_profile WHERE user_id = $1`

	_, err := r.DB.ExecContext(ctx, query, userID)

	return err
}
