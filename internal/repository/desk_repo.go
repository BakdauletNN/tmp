package repository

import (
	"context"
	"time"

	"tmp/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DeskRepository struct {
	Base
}

func NewDeskRepo(db *pgxpool.Pool) *DeskRepository {
	return &DeskRepository{
		Base: NewBase(db, "desks"),
	}
}

func (r *DeskRepository) ListByRoom(ctx context.Context, roomID int) ([]*models.Desk, error) {
	query := `SELECT id, room_id, label FROM desks WHERE room_id = $1`

	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Desk
	for rows.Next() {
		var d models.Desk
		if err := rows.Scan(&d.ID, &d.RoomID, &d.Label); err != nil {
			return nil, err
		}
		result = append(result, &d)
	}
	return result, rows.Err()
}

func (r *DeskRepository) ListAvailability(ctx context.Context, roomID int, start, end time.Time) ([]models.Desk, error) {
	query := `
		SELECT d.id, d.room_id, d.label,
			NOT EXISTS (
				SELECT 1 FROM bookings b
				WHERE b.room_id = d.room_id
				  AND (b.desk_id = d.id OR b.desk_id IS NULL)
				  AND b.start_time < $3
				  AND b.end_time > $2
			) AS available
		FROM desks d
		WHERE d.room_id = $1
		ORDER BY d.id
	`
	rows, err := r.db.Query(ctx, query, roomID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var desks []models.Desk
	for rows.Next() {
		var desk models.Desk
		if err := rows.Scan(&desk.ID, &desk.RoomID, &desk.Label, &desk.Available); err != nil {
			return nil, err
		}
		desks = append(desks, desk)
	}
	return desks, rows.Err()
}
