package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"tmp/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	Base
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{
		Base: NewBase(db, "bookings"),
	}
}

// CreateWithLock checks for a conflicting booking and inserts the new one
// inside a single transaction, serialized per room via an advisory lock.
// This closes the check-then-insert race that existed between separate
// HasConflict/WriteBooking calls: two concurrent requests for the same
// room now queue on the lock instead of both passing the conflict check.
func (r *BookingRepository) CreateWithLock(ctx context.Context, b *models.Booking) (conflict bool, err error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx) // no-op once committed

	if _, err = tx.Exec(ctx, `SELECT pg_advisory_xact_lock($1)`, int64(b.RoomID)); err != nil {
		return false, err
	}

	var exists bool
	if b.DeskID != nil {
		var deskBelongsToRoom bool
		err = tx.QueryRow(ctx, `
			SELECT EXISTS (SELECT 1 FROM desks WHERE id = $1 AND room_id = $2)
		`, *b.DeskID, b.RoomID).Scan(&deskBelongsToRoom)
		if err != nil {
			return false, err
		}
		if !deskBelongsToRoom {
			return false, errors.New("desk does not belong to room")
		}
		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM bookings
				WHERE room_id = $2
				  AND (desk_id = $1 OR desk_id IS NULL)
				  AND start_time < $4 AND end_time > $3
			)`, *b.DeskID, b.RoomID, b.StartTime, b.EndTime).Scan(&exists)
	} else {
		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM bookings
				WHERE room_id = $1 AND start_time < $3 AND end_time > $2
			)`, b.RoomID, b.StartTime, b.EndTime).Scan(&exists)
	}
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	if err = tx.QueryRow(ctx, `
		INSERT INTO "bookings" (room_id, user_id, start_time, end_time, desk_id, public_code)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, b.RoomID, b.UserID, b.StartTime, b.EndTime, b.DeskID, publicCode()).Scan(&b.ID); err != nil {
		return false, err
	}

	return false, tx.Commit(ctx)
}

// publicCode generates a short random token for bots/links to reference a
// booking by, instead of the sequential id (which is guessable/enumerable).
func publicCode() string {
	buf := make([]byte, 6) // 12 hex chars, ~2^48 space — not brute-forceable
	if _, err := rand.Read(buf); err != nil {
		// crypto/rand failing means the OS entropy source is broken; there's
		// nothing sane to fall back to, so surface it loudly instead of
		// silently handing out a predictable code.
		panic("public code generation: " + err.Error())
	}
	return hex.EncodeToString(buf)
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
		SELECT id, room_id, user_id, desk_id, start_time, end_time, COALESCE(public_code, '')
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
		if err := rows.Scan(&b.ID, &b.RoomID, &b.UserID, &b.DeskID, &b.StartTime, &b.EndTime, &b.PublicCode); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
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
