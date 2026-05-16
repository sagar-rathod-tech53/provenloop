package repositories

import (
	"context"
	"database/sql"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
)

type EducationRepository struct {
	DB *sql.DB
}

// CREATE
func (r *EducationRepository) Create(ctx context.Context, e models.Education) error {

	query := `
	INSERT INTO user_education (
		id, user_id, institution_name,
		degree, field_of_study,
		start_year, end_year,
		grade, description,
		is_verified, created_at, updated_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`

	_, err := r.DB.ExecContext(ctx, query,
		e.ID,
		e.UserID,
		e.InstitutionName,
		e.Degree,
		e.FieldOfStudy,
		e.StartYear,
		e.EndYear,
		e.Grade,
		e.Description,
		e.IsVerified,
		e.CreatedAt,
		e.UpdatedAt,
	)

	return err
}

// GET ALL
func (r *EducationRepository) GetAll(ctx context.Context, userID string) ([]models.Education, error) {

	query := `SELECT * FROM user_education WHERE user_id=$1 ORDER BY start_year DESC`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Education

	for rows.Next() {
		var e models.Education

		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.InstitutionName,
			&e.Degree,
			&e.FieldOfStudy,
			&e.StartYear,
			&e.EndYear,
			&e.Grade,
			&e.Description,
			&e.IsVerified,
			&e.CreatedAt,
			&e.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		list = append(list, e)
	}

	return list, nil
}

// UPDATE
func (r *EducationRepository) Update(ctx context.Context, e models.Education) error {

	query := `
	UPDATE user_education SET
		institution_name=$1,
		degree=$2,
		field_of_study=$3,
		start_year=$4,
		end_year=$5,
		grade=$6,
		description=$7,
		updated_at=NOW()
	WHERE id=$8 AND user_id=$9
	`

	_, err := r.DB.ExecContext(ctx, query,
		e.InstitutionName,
		e.Degree,
		e.FieldOfStudy,
		e.StartYear,
		e.EndYear,
		e.Grade,
		e.Description,
		e.ID,
		e.UserID,
	)

	return err
}

// DELETE
func (r *EducationRepository) Delete(ctx context.Context, id string, userID string) error {

	query := `DELETE FROM user_education WHERE id=$1 AND user_id=$2`

	_, err := r.DB.ExecContext(ctx, query, id, userID)

	return err
}
