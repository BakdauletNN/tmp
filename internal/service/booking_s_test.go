package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"tmp/internal/models"
)

type fakeBookingRepo struct {
	conflict bool
	err      error
	created  *models.Booking
}

func (f *fakeBookingRepo) CreateWithLock(
	ctx context.Context,
	b *models.Booking,
) (bool, error) {

	if f.err != nil {
		return false, f.err
	}

	if f.conflict {
		return true, nil
	}

	b.ID = 1
	b.PublicCode = "test123456"
	f.created = b

	return false, nil
}

func (f *fakeBookingRepo) GetUserBookings(
	ctx context.Context,
	userID int,
) ([]models.Booking, error) {
	return nil, nil
}

func (f *fakeBookingRepo) DeleteBooking(
	ctx context.Context,
	id int,
	userID int,
) error {
	return nil
}

type fakeRoomRepo struct {
	notFound bool
}

func (f *fakeRoomRepo) GetInfoRoom(
	ctx context.Context,
	id int,
) (models.Room, error) {

	if f.notFound {
		return models.Room{}, errors.New("room not found")
	}

	return models.Room{
		ID: id,
	}, nil
}

func TestCreateBooking(t *testing.T) {
	start := time.Date(
		2026,
		10,
		1,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	end := start.Add(time.Hour)

	tests := []struct {
		name        string
		booking     *fakeBookingRepo
		room        *fakeRoomRepo
		start       time.Time
		end         time.Time
		wantErr     error
		wantBooking bool
	}{
		{
			name:        "invalid time range",
			booking:     &fakeBookingRepo{},
			room:        &fakeRoomRepo{},
			start:       end,
			end:         start,
			wantBooking: false,
		},
		{
			name:        "room not found",
			booking:     &fakeBookingRepo{},
			room:        &fakeRoomRepo{notFound: true},
			start:       start,
			end:         end,
			wantBooking: false,
		},
		{
			name:        "booking conflict",
			booking:     &fakeBookingRepo{conflict: true},
			room:        &fakeRoomRepo{},
			start:       start,
			end:         end,
			wantErr:     ErrRoomBooked,
			wantBooking: false,
		},
		{
			name:        "booking created successfully",
			booking:     &fakeBookingRepo{},
			room:        &fakeRoomRepo{},
			start:       start,
			end:         end,
			wantBooking: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			service := &BookingService{
				repo:     tt.booking,
				roomRepo: tt.room,
			}

			booking, err := service.CreateBooking(
				context.Background(),
				1,
				42,
				tt.start,
				tt.end,
				nil,
			)

			if tt.wantBooking {
				if err != nil {
					t.Fatalf(
						"expected booking to be created, got error: %v",
						err,
					)
				}

				if booking == nil {
					t.Fatal("expected booking, got nil")
				}

				if booking.ID == 0 {
					t.Fatal("expected booking ID to be set")
				}

				if booking.PublicCode == "" {
					t.Fatal("expected public code to be set")
				}

				if tt.booking.created == nil {
					t.Fatal("expected repository CreateWithLock to be called")
				}

				return
			}

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.wantErr,
					err,
				)
			}
		})
	}
}
