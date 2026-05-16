package repositories

import (
	"context"
	"database/sql"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
)

type ExperienceRepository struct {
	DB *sql.DB
}

// CREATE
func (r *ExperienceRepository) Create(ctx context.Context, e models.Experience) error {

	query := `
	INSERT INTO user_experience (
		id, user_id,
		company_name, role, employment_type,
		location, start_date, end_date,
		is_current, description,
		is_verified, created_at, updated_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`

	_, err := r.DB.ExecContext(ctx, query,
		e.ID,
		e.UserID,
		e.CompanyName,
		e.Role,
		e.EmploymentType,
		e.Location,
		e.StartDate,
		e.EndDate,
		e.IsCurrent,
		e.Description,
		e.IsVerified,
		e.CreatedAt,
		e.UpdatedAt,
	)

	return err
}

// GET ALL
func (r *ExperienceRepository) GetAll(ctx context.Context, userID string) ([]models.Experience, error) {

	query := `SELECT * FROM user_experience WHERE user_id=$1 ORDER BY start_date DESC`

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Experience

	for rows.Next() {
		var e models.Experience

		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.CompanyName,
			&e.Role,
			&e.EmploymentType,
			&e.Location,
			&e.StartDate,
			&e.EndDate,
			&e.IsCurrent,
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
func (r *ExperienceRepository) Update(ctx context.Context, e models.Experience) error {

	query := `
	UPDATE user_experience SET
		company_name=$1,
		role=$2,
		employment_type=$3,
		location=$4,
		start_date=$5,
		end_date=$6,
		is_current=$7,
		description=$8,
		updated_at=NOW()
	WHERE id=$9 AND user_id=$10
	`

	_, err := r.DB.ExecContext(ctx, query,
		e.CompanyName,
		e.Role,
		e.EmploymentType,
		e.Location,
		e.StartDate,
		e.EndDate,
		e.IsCurrent,
		e.Description,
		e.ID,
		e.UserID,
	)

	return err
}

// DELETE
func (r *ExperienceRepository) Delete(ctx context.Context, id string, userID string) error {

	query := `DELETE FROM user_experience WHERE id=$1 AND user_id=$2`

	_, err := r.DB.ExecContext(ctx, query, id, userID)

	return err
}
