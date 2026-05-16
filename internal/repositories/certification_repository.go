package repositories

import (
	"context"
	"database/sql"

	"github.com/sagar-rathod-tech53/provenloop/internal/models"
)

type CertificationRepository struct {
	DB *sql.DB
}

func (r *CertificationRepository) Create(
	ctx context.Context,
	c models.Certification,
) error {

	query := `
	INSERT INTO user_certifications(
	id,
	user_id,
	title,
	issuing_organization,
	credential_id,
	credential_url,
	issue_date,
	expiry_date,
	is_verified,
	created_at,
	updated_at
	)
	VALUES(
	$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
	)
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		c.ID,
		c.UserID,
		c.Title,
		c.IssuingOrganization,
		c.CredentialID,
		c.CredentialURL,
		c.IssueDate,
		c.ExpiryDate,
		c.IsVerified,
		c.CreatedAt,
		c.UpdatedAt,
	)

	return err
}

func (r *CertificationRepository) GetAll(
	ctx context.Context,
	userID string,
) ([]models.Certification, error) {

	query := `
	SELECT
	id,
	user_id,
	title,
	issuing_organization,
	credential_id,
	credential_url,
	issue_date,
	expiry_date,
	is_verified,
	created_at,
	updated_at
	FROM user_certifications
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

	var list []models.Certification

	for rows.Next() {

		var c models.Certification

		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Title,
			&c.IssuingOrganization,
			&c.CredentialID,
			&c.CredentialURL,
			&c.IssueDate,
			&c.ExpiryDate,
			&c.IsVerified,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		list = append(
			list,
			c,
		)
	}

	return list, nil
}

func (r *CertificationRepository) Update(
	ctx context.Context,
	c models.Certification,
) error {

	query := `
	UPDATE user_certifications
	SET
	title=$1,
	issuing_organization=$2,
	credential_id=$3,
	credential_url=$4,
	issue_date=$5,
	expiry_date=$6,
	updated_at=NOW()
	WHERE id=$7
	AND user_id=$8
	`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		c.Title,
		c.IssuingOrganization,
		c.CredentialID,
		c.CredentialURL,
		c.IssueDate,
		c.ExpiryDate,
		c.ID,
		c.UserID,
	)

	return err
}

func (r *CertificationRepository) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {

	_, err := r.DB.ExecContext(
		ctx,
		`
	DELETE FROM user_certifications
	WHERE id=$1
	AND user_id=$2
	`,
		id,
		userID,
	)

	return err
}
