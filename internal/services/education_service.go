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

type EducationService struct {
	Repo *repositories.EducationRepository
}

var eduMutex sync.Mutex

// CREATE
func (s *EducationService) Create(ctx context.Context, e models.Education) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	eduMutex.Lock()
	defer eduMutex.Unlock()

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

	_ = config.RDB.Del(ctx, "edu:"+e.UserID)

	return nil
}

// GET ALL (CACHE)
func (s *EducationService) GetAll(ctx context.Context, userID string) ([]models.Education, error) {

	cacheKey := "edu:" + userID

	if cached, err := config.RDB.Get(ctx, cacheKey).Result(); err == nil {
		var data []models.Education
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
func (s *EducationService) Update(ctx context.Context, e models.Education) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := s.Repo.Update(ctx, e)
	if err != nil {
		return err
	}

	_ = config.RDB.Del(ctx, "edu:"+e.UserID)

	return nil
}

// DELETE
func (s *EducationService) Delete(ctx context.Context, id string, userID string) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := s.Repo.Delete(ctx, id, userID)
	if err != nil {
		return err
	}

	_ = config.RDB.Del(ctx, "edu:"+userID)

	return nil
}
