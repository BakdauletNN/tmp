package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"tmp/internal/logger"
	"tmp/internal/models"
	"tmp/internal/repository"
)

var ErrRoomBooked = errors.New("room is already booked for this time")

// bookingRepo and roomRepo are satisfied by *repository.BookingRepository and
// *repository.RoomRepository respectively (structural typing — no change
// needed where BookingService is constructed). They exist so tests can pass
// a fake instead of hitting Postgres.
type bookingRepo interface {
	CreateWithLock(ctx context.Context, b *models.Booking) (conflict bool, err error)
	GetUserBookings(ctx context.Context, userID int) ([]models.Booking, error)
	DeleteBooking(ctx context.Context, id, userID int) error
}

type roomRepo interface {
	GetInfoRoom(ctx context.Context, id int) (models.Room, error)
}

type BookingService struct {
	repo     bookingRepo
	roomRepo roomRepo
}

func NewBookingService(repo *repository.BookingRepository, roomRepo *repository.RoomRepository) *BookingService {
	return &BookingService{repo: repo, roomRepo: roomRepo}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID, roomID int, start, end time.Time, deskID *int) (*models.Booking, error) {
	log := logger.FromContext(ctx)

	if !end.After(start) {
		log.Warn("invalid booking time range", slog.Time("start_time", start), slog.Time("end_time", end))
		return nil, errors.New("end_time must be after start_time")
	}

	if _, err := s.roomRepo.GetInfoRoom(ctx, roomID); err != nil {
		log.Warn("booking attempt for unknown room", slog.Int("room_id", roomID), slog.String("error", err.Error()))
		return nil, errors.New("room not found")
	}

	booking := &models.Booking{
		RoomID:    roomID,
		UserID:    userID,
		DeskID:    deskID,
		StartTime: start,
		EndTime:   end,
	}

	conflict, err := s.repo.CreateWithLock(ctx, booking)
	if err != nil {
		log.Error("booking check failed", slog.Int("room_id", roomID), slog.String("error", err.Error()))
		return nil, err
	}
	if conflict {
		log.Warn("booking conflict detected", slog.Int("room_id", roomID), slog.Int("user_id", userID), slog.Time("start_time", start), slog.Time("end_time", end))
		return nil, ErrRoomBooked
	}

	log.Info("booking created",
		slog.Int("booking_id", booking.ID),
		slog.Int("room_id", roomID),
		slog.Int("user_id", userID),
	)
	return booking, nil
}

func (s *BookingService) GetUserBookings(ctx context.Context, userID int) ([]models.Booking, error) {
	log := logger.FromContext(ctx)
	bookings, err := s.repo.GetUserBookings(ctx, userID)
	if err != nil {
		log.Error("failed to load user bookings", slog.Int("user_id", userID), slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("user bookings loaded", slog.Int("user_id", userID), slog.Int("count", len(bookings)))
	return bookings, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID, userID int) error {
	log := logger.FromContext(ctx)
	if err := s.repo.DeleteBooking(ctx, bookingID, userID); err != nil {
		log.Warn("booking cancel failed", slog.Int("booking_id", bookingID), slog.Int("user_id", userID), slog.String("error", err.Error()))
		return err
	}
	log.Info("booking cancelled", slog.Int("booking_id", bookingID), slog.Int("user_id", userID))
	return nil
}
