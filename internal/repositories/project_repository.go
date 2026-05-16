package repositories

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
)

type ProjectRepository struct {
	DB *sql.DB
}

func (r *ProjectRepository) Create(
	ctx context.Context,
	p models.Project,
) error {

	query := `
	INSERT INTO user_projects(
	id,
	user_id,
	title,
	description,
	tech_stack,
	github_url,
	live_url,
	project_video_url,
	start_date,
	end_date,
	is_verified,
	created_at,
	updated_at
	)
	VALUES(
	$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
	)
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		p.ID,
		p.UserID,
		p.Title,
		p.Description,
		pq.Array(p.TechStack),
		p.GithubURL,
		p.LiveURL,
		p.ProjectVideoURL,
		p.StartDate,
		p.EndDate,
		p.IsVerified,
		p.CreatedAt,
		p.UpdatedAt,
	)

	return err
}

func (r *ProjectRepository) GetAll(
	ctx context.Context,
	userID string,
) ([]models.Project, error) {

	query := `
	SELECT
	id,
	user_id,
	title,
	description,
	tech_stack,
	github_url,
	live_url,
	project_video_url,
	start_date,
	end_date,
	is_verified,
	created_at,
	updated_at
	FROM user_projects
	WHERE user_id=$1
	ORDER BY created_at DESC
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

	var projects []models.Project

	for rows.Next() {

		var p models.Project

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Description,
			pq.Array(&p.TechStack),
			&p.GithubURL,
			&p.LiveURL,
			&p.ProjectVideoURL,
			&p.StartDate,
			&p.EndDate,
			&p.IsVerified,
			&p.CreatedAt,
			&p.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		projects = append(
			projects,
			p,
		)

	}

	return projects, nil
}

func (r *ProjectRepository) Update(
	ctx context.Context,
	p models.Project,
) error {

	query := `
	UPDATE user_projects
	SET
	title=$1,
	description=$2,
	tech_stack=$3,
	github_url=$4,
	live_url=$5,
	project_video_url=$6,
	start_date=$7,
	end_date=$8,
	updated_at=NOW()
	WHERE id=$9
	AND user_id=$10
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		p.Title,
		p.Description,
		pq.Array(p.TechStack),
		p.GithubURL,
		p.LiveURL,
		p.ProjectVideoURL,
		p.StartDate,
		p.EndDate,
		p.ID,
		p.UserID,
	)

	return err
}

func (r *ProjectRepository) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {

	_, err := r.DB.ExecContext(
		ctx,
		`DELETE FROM user_projects
	WHERE id=$1
	AND user_id=$2`,
		id,
		userID,
	)

	return err
}
