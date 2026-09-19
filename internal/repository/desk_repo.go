package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"tmp/internal/models"
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