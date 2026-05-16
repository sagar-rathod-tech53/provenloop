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

type ProjectService struct {
	Repo *repositories.ProjectRepository
}

var projectMutex sync.Mutex

func (s *ProjectService) Create(
	ctx context.Context,
	p models.Project,
) error {

	ctx, cancel := context.WithTimeout(
		ctx,
		10*time.Second,
	)

	defer cancel()

	projectMutex.Lock()
	defer projectMutex.Unlock()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		p.ID = uuid.New().String()

	}()

	go func() {

		defer wg.Done()

		now := time.Now()

		p.CreatedAt = now
		p.UpdatedAt = now

	}()

	wg.Wait()

	err := s.Repo.Create(
		ctx,
		p,
	)

	if err != nil {
		return err
	}

	_ = config.RDB.Del(
		ctx,
		"projects:"+p.UserID,
	)

	return nil
}

func (s *ProjectService) GetAll(
	ctx context.Context,
	userID string,
) ([]models.Project, error) {

	cacheKey := "projects:" + userID

	cache, err := config.RDB.Get(
		ctx,
		cacheKey,
	).Result()

	if err == nil {

		var p []models.Project

		json.Unmarshal(
			[]byte(cache),
			&p,
		)

		return p, nil
	}

	data, err := s.Repo.GetAll(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(data)

	config.RDB.Set(
		ctx,
		cacheKey,
		b,
		10*time.Minute,
	)

	return data, nil
}

func (s *ProjectService) Update(
	ctx context.Context,
	p models.Project,
) error {

	err := s.Repo.Update(
		ctx,
		p,
	)

	if err != nil {
		return err
	}

	_ = config.RDB.Del(
		ctx,
		"projects:"+p.UserID,
	)

	return nil
}

func (s *ProjectService) Delete(
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
		"projects:"+userID,
	)

	return nil
}
