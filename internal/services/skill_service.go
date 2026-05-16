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

type SkillService struct {
	Repo *repositories.SkillRepository
}

var skillMutex sync.Mutex

func (s *SkillService) Create(
	ctx context.Context,
	sk models.Skill,
) error {

	ctx, cancel := context.WithTimeout(
		ctx,
		10*time.Second,
	)

	defer cancel()

	skillMutex.Lock()
	defer skillMutex.Unlock()

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {

		defer wg.Done()

		sk.ID = uuid.New().String()

	}()

	go func() {

		defer wg.Done()

		now := time.Now()

		sk.CreatedAt = now
		sk.UpdatedAt = now

	}()

	wg.Wait()

	err := s.Repo.Create(
		ctx,
		sk,
	)

	if err != nil {
		return err
	}

	_ = config.RDB.Del(
		ctx,
		"skills:"+sk.UserID,
	)

	return nil
}

func (s *SkillService) GetAll(
	ctx context.Context,
	userID string,
) ([]models.Skill, error) {

	cacheKey := "skills:" + userID

	cache, err := config.RDB.Get(
		ctx,
		cacheKey,
	).Result()

	if err == nil {

		var data []models.Skill

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

func (s *SkillService) Update(
	ctx context.Context,
	sk models.Skill,
) error {

	err := s.Repo.Update(
		ctx,
		sk,
	)

	if err != nil {
		return err
	}

	_ = config.RDB.Del(
		ctx,
		"skills:"+sk.UserID,
	)

	return nil
}

func (s *SkillService) Delete(
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
		"skills:"+userID,
	)

	return nil
}
