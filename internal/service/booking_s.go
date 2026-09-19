package service

import (
	"context"
	"errors"
	"time"

	"tmp/internal/models"
	"tmp/internal/repository"
)

var ErrRoomBooked = errors.New("room is already booked for this time")

type BookingService struct {
	repo     *repository.BookingRepository
	roomRepo *repository.RoomRepository
}

func NewBookingService(repo *repository.BookingRepository, roomRepo *repository.RoomRepository) *BookingService {
	return &BookingService{repo: repo, roomRepo: roomRepo}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID, roomID int, start, end time.Time, deskID *int) (*models.Booking, error) {
	if !end.After(start) {
		return nil, errors.New("end_time must be after start_time")
	}

	if _, err := s.roomRepo.GetInfoRoom(ctx, roomID); err != nil {
		return nil, errors.New("room not found")
	}

	var conflict bool
	var err error

	if deskID != nil {
		conflict, err = s.repo.HasConflictDesk(ctx, *deskID, start, end)
	} else {
		conflict, err = s.repo.HasConflict(ctx, roomID, start, end)
	}
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, ErrRoomBooked
	}

	booking := &models.Booking{
		RoomID:    roomID,
		UserID:    userID,
		DeskID:    deskID,
		StartTime: start,
		EndTime:   end,
	}
	if err := s.repo.WriteBooking(ctx, booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *BookingService) GetUserBookings(ctx context.Context, userID int) ([]models.Booking, error) {
	return s.repo.GetUserBookings(ctx, userID)
}