package repository

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tmp/internal/models"
)

type BookingRepository struct {
	Base
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{
		Base: NewBase(db, "bookings"),
	}
}

func (r *BookingRepository) WriteBooking(ctx context.Context, b *models.Booking) error {
	query := `
		INSERT INTO "bookings" 
			(room_id, user_id, start_time, end_time, desk_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(
		ctx, query,
		b.RoomID, b.UserID, b.StartTime, b.EndTime, b.DeskID,
	).Scan(&b.ID)
}

func (r *BookingRepository) GetBooking(ctx context.Context, id int) (*models.Booking, error) {
	query := `
		SELECT id, room_id, user_id, desk_id, start_time, end_time
		FROM "bookings"
		WHERE id = $1
	`
	var b models.Booking
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.RoomID, &b.UserID, &b.DeskID, &b.StartTime, &b.EndTime,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) GetUserBookings(ctx context.Context, userID int) ([]models.Booking, error) {
	query := `
		SELECT id, room_id, user_id, desk_id, start_time, end_time
		FROM "bookings"
		WHERE user_id = $1
		ORDER BY start_time
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.RoomID, &b.UserID, &b.DeskID, &b.StartTime, &b.EndTime); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

func (r *BookingRepository) HasConflict(
	ctx context.Context,
	roomID int,
	startTime time.Time,
	endTime time.Time,
) (bool, error) {

	query := `
        SELECT EXISTS (
            SELECT 1
            FROM bookings
            WHERE room_id = $1
              AND start_time < $3
              AND end_time > $2
        )
    `

	var exists bool

	err := r.db.QueryRow(
		ctx,
		query,
		roomID,
		startTime,
		endTime,
	).Scan(&exists)

	return exists, err
}


func (r *BookingRepository) HasConflictDesk(
	ctx context.Context,
	DeskID int,
	startTime time.Time,
	endTime time.Time,
) (bool, error) {

	query := `
        SELECT EXISTS (
            SELECT 1
            FROM bookings
            WHERE desk_id = $1
              AND start_time < $3
              AND end_time > $2
        )
    `

	var exists bool

	err := r.db.QueryRow(
		ctx,
		query,
		DeskID,
		startTime,
		endTime,
	).Scan(&exists)

	return exists, err
}


func (r *BookingRepository) DeleteBooking(ctx context.Context, id, userID int) error {
	query := `DELETE FROM bookings WHERE id = $1 AND user_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows 
	}
	return nil
}