package service

import (
	"context"

	"tmp/internal/models"
	"tmp/internal/repository"
)

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) GetRoom(ctx context.Context, id int, isAdmin bool) (models.Room, error) {
	room, err := s.repo.GetInfoRoom(ctx, id)
	if err != nil {
		return room, err
	}
	if !isAdmin {
		room.AccessCode = nil
	}
	return room, nil
}

func (s *RoomService) SetAccessCode(ctx context.Context, roomID, code int) error {
	return s.repo.UpdateAccessCode(ctx, roomID, code)
}

func (s *RoomService) SearchRoom(ctx context.Context, filter *models.RoomFilter) ([]models.Room, error) {
	return s.repo.SearchRoom(ctx, filter)
}