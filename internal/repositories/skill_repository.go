package repositories

import (
	"context"
	"database/sql"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
)

type SkillRepository struct {
	DB *sql.DB
}

func (r *SkillRepository) Create(
	ctx context.Context,
	s models.Skill,
) error {

	query := `
	INSERT INTO user_skills(
	id,
	user_id,
	skill_name,
	proficiency_level,
	endorsements_count,
	is_verified,
	created_at,
	updated_at
	)
	VALUES(
	$1,$2,$3,$4,$5,$6,$7,$8
	)
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		s.ID,
		s.UserID,
		s.SkillName,
		s.ProficiencyLevel,
		s.EndorsementsCount,
		s.IsVerified,
		s.CreatedAt,
		s.UpdatedAt,
	)

	return err
}

func (r *SkillRepository) GetAll(
	ctx context.Context,
	userID string,
) ([]models.Skill, error) {

	query := `
	SELECT
	id,
	user_id,
	skill_name,
	proficiency_level,
	endorsements_count,
	is_verified,
	created_at,
	updated_at
	FROM user_skills
	WHERE user_id=$1
	ORDER BY endorsements_count DESC
	`

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var skills []models.Skill

	for rows.Next() {

		var s models.Skill

		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.SkillName,
			&s.ProficiencyLevel,
			&s.EndorsementsCount,
			&s.IsVerified,
			&s.CreatedAt,
			&s.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		skills = append(
			skills,
			s,
		)

	}

	return skills, nil
}

func (r *SkillRepository) Update(
	ctx context.Context,
	s models.Skill,
) error {

	query := `
	UPDATE user_skills
	SET
	skill_name=$1,
	proficiency_level=$2,
	updated_at=NOW()
	WHERE id=$3
	AND user_id=$4
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		s.SkillName,
		s.ProficiencyLevel,
		s.ID,
		s.UserID,
	)

	return err
}

func (r *SkillRepository) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {

	_, err := r.DB.ExecContext(
		ctx,
		`
	DELETE FROM user_skills
	WHERE id=$1
	AND user_id=$2
	`,
		id,
		userID,
	)

	return err
}
