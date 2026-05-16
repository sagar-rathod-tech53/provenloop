package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sagar-rathod-tech53/provenloop/config"
	"github.com/sagar-rathod-tech53/provenloop/internal/models"
	"github.com/sagar-rathod-tech53/provenloop/internal/repositories"
)

type UserProfileService struct {
	Repo *repositories.UserProfileRepository
}

var profileMutex sync.Mutex

// CREATE PROFILE
func (s *UserProfileService) CreateProfile(ctx context.Context, p models.UserProfile) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	exists, err := s.Repo.ProfileExists(ctx, p.UserID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("profile already exists")
	}

	// concurrency
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
		p.LastActiveAt = now
	}()

	wg.Wait()

	err = s.Repo.CreateProfile(ctx, p)
	if err != nil {
		return err
	}

	_ = config.RDB.Del(ctx, "profile:"+p.UserID)

	return nil
}

// GET PROFILE (CACHE + DB)
func (s *UserProfileService) GetProfile(ctx context.Context, userID string) (models.UserProfile, error) {

	cacheKey := "profile:" + userID

	// =========================
	// 🔥 REDIS CACHE CHECK
	// =========================
	if cached, err := config.RDB.Get(ctx, cacheKey).Result(); err == nil {
		var p models.UserProfile
		_ = json.Unmarshal([]byte(cached), &p)
		return p, nil
	}

	// =========================
	// 🗄️ DB CALL (JOIN QUERY)
	// =========================
	profile, err := s.Repo.GetProfile(ctx, userID)
	if err != nil {
		return profile, err
	}

	// =========================
	// ⚡ CACHE STORE
	// =========================
	bytes, _ := json.Marshal(profile)
	_ = config.RDB.Set(ctx, cacheKey, bytes, 10*time.Minute).Err()

	return profile, nil
}

// =========================
// UPDATE PROFILE
// =========================
func (s *UserProfileService) UpdateProfile(
	ctx context.Context,
	userID string,
	p models.UserProfile,
) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 🔒 concurrency lock (avoid race update issues)
	profileMutex.Lock()
	defer profileMutex.Unlock()

	var wg sync.WaitGroup
	wg.Add(2)

	// async timestamp update
	go func() {
		defer wg.Done()
		p.UpdatedAt = time.Now()
		p.LastActiveAt = time.Now()
	}()

	// async validation hook (future extension)
	go func() {
		defer wg.Done()
		// placeholder for validation or enrichment
	}()

	wg.Wait()

	err := s.Repo.UpdateProfile(ctx, userID, p)
	if err != nil {
		return err
	}

	// invalidate cache
	_ = config.RDB.Del(ctx, "profile:"+userID)

	return nil
}

// =========================
// DELETE PROFILE
// =========================
func (s *UserProfileService) DeleteProfile(
	ctx context.Context,
	userID string,
) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// concurrency safe delete
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)

	go func() {
		defer wg.Done()

		err := s.Repo.DeleteProfile(ctx, userID)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
	}

	// clear cache
	_ = config.RDB.Del(ctx, "profile:"+userID)

	return nil
}
