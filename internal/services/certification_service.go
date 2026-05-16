package services

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/sagar-rathod-tech53/provenloop/config"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/internal/repositories"
)

type CertificationService struct {
	Repo *repositories.CertificationRepository
}

var certMutex sync.Mutex

func (s *CertificationService) Create(
	ctx context.Context,
	c models.Certification,
) error {

	ctx, cancel := context.WithTimeout(
		ctx,
		10*time.Second,
	)

	defer cancel()

	certMutex.Lock()
	defer certMutex.Unlock()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {

		defer wg.Done()

		c.ID = uuid.New().String()

	}()

	go func() {

		defer wg.Done()

		now := time.Now()

		c.CreatedAt = now
		c.UpdatedAt = now

	}()

	wg.Wait()

	err := s.Repo.Create(
		ctx,
		c,
	)

	if err != nil {
		return err
	}

	_ = config.RDB.Del(
		ctx,
		"certifications:"+c.UserID,
	)

	return nil
}

func (s *CertificationService) GetAll(
	ctx context.Context,
	userID string,
) ([]models.Certification, error) {

	cacheKey := "certifications:" + userID

	cache, err := config.RDB.Get(
		ctx,
		cacheKey,
	).Result()

	if err == nil {

		var data []models.Certification

		_ = json.Unmarshal(
			[]byte(cache),
			&data,
		)

		return data, nil
	}

	data, err := s.Repo.GetAll(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(data)

	_ = config.RDB.Set(
		ctx,
		cacheKey,
		b,
		10*time.Minute,
	).Err()

	return data, nil
}

func (s *CertificationService) Update(
	ctx context.Context,
	c models.Certification,
) error {

	err := s.Repo.Update(
		ctx,
		c,
	)

	if err != nil {
		return err
	}

	_ = config.RDB.Del(
		ctx,
		"certifications:"+c.UserID,
	)

	return nil
}

func (s *CertificationService) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {

	err := s.Repo.Delete(
		ctx,
		id,
		userID,
	)

	if err != nil {
		return err
	}

	_ = config.RDB.Del(
		ctx,
		"certifications:"+userID,
	)

	return nil
}
