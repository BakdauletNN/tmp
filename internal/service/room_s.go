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

type RoomService struct {
	repo  *repository.RoomRepository
	cache *cache.Client
}

func NewRoomService(repo *repository.RoomRepository, cacheClient *cache.Client) *RoomService {
	return &RoomService{repo: repo, cache: cacheClient}
}

func (s *RoomService) GetRoom(ctx context.Context, id int, isAdmin bool) (models.Room, error) {
	log := logger.FromContext(ctx)

	if s.cache != nil {
		var cached models.Room
		if err := s.cache.Get(ctx, cache.Room(id), &cached); err == nil {
			log.Debug("room loaded from cache", slog.Int("room_id", id))
			if !isAdmin {
				cached.AccessCode = nil
			}
			return cached, nil
		} else if err != cache.ErrMiss {
			log.Warn("room cache read failed", slog.Int("room_id", id), slog.String("error", err.Error()))
		}
	}

	room, err := s.repo.GetInfoRoom(ctx, id)
	if err != nil {
		log.Warn("room lookup failed", slog.Int("room_id", id), slog.String("error", err.Error()))
		return room, err
	}
	if s.cache != nil {
		if err := s.cache.Set(ctx, cache.Room(id), room, 5*time.Minute); err != nil {
			log.Warn("room cache write failed", slog.Int("room_id", id), slog.String("error", err.Error()))
		}
	}
	if !isAdmin {
		room.AccessCode = nil
	}
	return room, nil
}

func (s *RoomService) SetAccessCode(ctx context.Context, roomID, code int) error {
	log := logger.FromContext(ctx)
	if err := s.repo.UpdateAccessCode(ctx, roomID, code); err != nil {
		log.Error("update room access code failed", slog.Int("room_id", roomID), slog.Int("code", code), slog.String("error", err.Error()))
		return err
	}
	if s.cache != nil {
		if err := s.cache.Delete(ctx, cache.Room(roomID)); err != nil {
			log.Warn("room cache invalidation failed", slog.Int("room_id", roomID), slog.String("error", err.Error()))
		}
	}
	return nil
}

func (s *RoomService) SearchRoom(ctx context.Context, filter *models.RoomFilter) ([]models.Room, error) {
	log := logger.FromContext(ctx)
	key, err := cache.Rooms(filter)
	if err != nil {
		log.Warn("room cache key creation failed", slog.String("error", err.Error()))
	} else if s.cache != nil {
		var cached []models.Room
		if err := s.cache.Get(ctx, key, &cached); err == nil {
			log.Debug("rooms loaded from cache", slog.Int("count", len(cached)))
			return cached, nil
		} else if err != cache.ErrMiss {
			log.Warn("rooms cache read failed", slog.String("error", err.Error()))
		}
	}

	rooms, err := s.repo.SearchRoom(ctx, filter)
	if err != nil {
		log.Error("search rooms failed", slog.String("error", err.Error()))
		return nil, err
	}
	if s.cache != nil {
		if err := s.cache.Set(ctx, key, rooms, time.Minute); err != nil {
			log.Warn("rooms cache write failed", slog.String("error", err.Error()))
		}
	}
	log.Debug("rooms searched", slog.Int("count", len(rooms)))
	return rooms, nil
}
