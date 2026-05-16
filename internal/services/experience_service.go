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

type ExperienceService struct {
	Repo *repositories.ExperienceRepository
}

var expMutex sync.Mutex

// CREATE
func (s *ExperienceService) Create(ctx context.Context, e models.Experience) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	expMutex.Lock()
	defer expMutex.Unlock()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		e.ID = uuid.New().String()
	}()

	go func() {
		defer wg.Done()
		now := time.Now()
		e.CreatedAt = now
		e.UpdatedAt = now
	}()

	wg.Wait()

	err := s.Repo.Create(ctx, e)
	if err != nil {
		return err
	}

	_ = config.RDB.Del(ctx, "exp:"+e.UserID)

	return nil
}

// GET ALL (CACHE)
func (s *ExperienceService) GetAll(ctx context.Context, userID string) ([]models.Experience, error) {

	cacheKey := "exp:" + userID

	if cached, err := config.RDB.Get(ctx, cacheKey).Result(); err == nil {
		var data []models.Experience
		_ = json.Unmarshal([]byte(cached), &data)
		return data, nil
	}

	data, err := s.Repo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	bytes, _ := json.Marshal(data)
	_ = config.RDB.Set(ctx, cacheKey, bytes, 10*time.Minute)

	return data, nil
}

// UPDATE
func (s *ExperienceService) Update(ctx context.Context, e models.Experience) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := s.Repo.Update(ctx, e)
	if err != nil {
		return err
	}

	_ = config.RDB.Del(ctx, "exp:"+e.UserID)

	return nil
}

// DELETE
func (s *ExperienceService) Delete(ctx context.Context, id string, userID string) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := s.Repo.Delete(ctx, id, userID)
	if err != nil {
		return err
	}

	_ = config.RDB.Del(ctx, "exp:"+userID)

	return nil
}
