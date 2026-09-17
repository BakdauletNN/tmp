package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tmp/internal/models"
)

type RoomRepository struct {
	Base
}

func NewRoomRepo(db *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{
		Base: NewBase(db, "rooms"),
	}
}

func (r *RoomRepository) UpdateAccessCode(ctx context.Context, roomID int, code int) error {
	query := `
		UPDATE rooms 
		SET access_code = $1 
		WHERE id = $2
	`

	cmdTag, err := r.db.Exec(ctx, query, code, roomID)
	if err != nil {
		return fmt.Errorf("failed to update access code for room %d: %w", roomID, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *RoomRepository) GetInfoRoom(ctx context.Context, id int) (models.Room, error) {
	query := `
		SELECT 
			id, 
			office_id, 
			type, 
			price_hour, 
			qty_desks, 
			qty_person, 
			access_code, 
			has_air_conditioner, 
			has_prayer_room 
		FROM rooms 
		WHERE id = $1
	`

	var room models.Room

	err := r.db.QueryRow(ctx, query, id).Scan(
		&room.ID,
		&room.OfficeID,
		&room.Type,
		&room.PriceHour,
		&room.QtyDesks,
		&room.QtyPerson,
		&room.AccessCode,
		&room.HasAirConditioner,
		&room.HasPrayerRoom,
	)
	if err != nil {
		return models.Room{}, err
	}

	return room, nil
}
