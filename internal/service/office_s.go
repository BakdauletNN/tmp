package service

import (
	"context"
	"log/slog"
	"time"

	"tmp/internal/cache"
	"tmp/internal/logger"
	"tmp/internal/models"
	"tmp/internal/repository"
)

type OfficeService struct {
	repo  *repository.OfficeRepository
	cache *cache.Client
}

func NewOfficeService(repo *repository.OfficeRepository, cacheClient *cache.Client) *OfficeService {
	return &OfficeService{repo: repo, cache: cacheClient}
}

func (s *OfficeService) GetOffice(ctx context.Context, id int, isAdmin bool) (*models.Office, error) {
	log := logger.FromContext(ctx)

	if s.cache != nil {
		var cached models.Office
		if err := s.cache.Get(ctx, cache.Office(id), &cached); err == nil {
			log.Debug("office loaded from cache", slog.Int("office_id", id))
			if !isAdmin {
				cached.OwnerID = 0
			}
			return &cached, nil
		} else if err != cache.ErrMiss {
			log.Warn("office cache read failed", slog.Int("office_id", id), slog.String("error", err.Error()))
		}
	}

	office, err := s.repo.GetOfficeInfo(ctx, id)
	if err != nil {
		log.Warn("office lookup failed", slog.Int("office_id", id), slog.String("error", err.Error()))
		return nil, err
	}
	if s.cache != nil {
		if err := s.cache.Set(ctx, cache.Office(id), office, 5*time.Minute); err != nil {
			log.Warn("office cache write failed", slog.Int("office_id", id), slog.String("error", err.Error()))
		}
	}
	if !isAdmin {
		office.OwnerID = 0
	}
	return office, nil
}

func (s *OfficeService) SearchByAddress(ctx context.Context, f models.OfficeFilter) ([]models.Office, error) {
	log := logger.FromContext(ctx)
	key, err := cache.Offices(f)
	if err != nil {
		log.Warn("office cache key creation failed", slog.String("error", err.Error()))
	} else if s.cache != nil {
		var cached []models.Office
		if err := s.cache.Get(ctx, key, &cached); err == nil {
			log.Debug("offices loaded from cache", slog.Int("count", len(cached)))
			return cached, nil
		} else if err != cache.ErrMiss {
			log.Warn("offices cache read failed", slog.String("error", err.Error()))
		}
	}
	offices, err := s.repo.SearchOffices(ctx, f)
	if err != nil {
		log.Error("search offices failed", slog.String("error", err.Error()))
		return nil, err
	}
	if s.cache != nil {
		if err := s.cache.Set(ctx, key, offices, time.Minute); err != nil {
			log.Warn("offices cache write failed", slog.String("error", err.Error()))
		}
	}
	log.Debug("offices searched", slog.Int("count", len(offices)))
	return offices, nil
}
