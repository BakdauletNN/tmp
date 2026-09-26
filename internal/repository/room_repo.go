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

func (r *RoomRepository) SearchRoom(ctx context.Context, f *models.RoomFilter) ([]models.Room, error) {
	query := `
		SELECT id, office_id, type, price_hour, qty_desks, qty_person, access_code, has_air_conditioner, has_prayer_room
		FROM "rooms"
		WHERE 1=1
	`
	var args []interface{}
	argN := 1

	if f.OfficeID > 0 {
		query += fmt.Sprintf(" AND office_id = $%d", argN)
		args = append(args, f.OfficeID)
		argN++
	}
	if f.Type != "" {
		query += fmt.Sprintf(" AND type ILIKE '%%' || $%d || '%%'", argN)
		args = append(args, f.Type)
		argN++
	}
	if f.MaxPriceHour > 0 {
		query += fmt.Sprintf(" AND price_hour <= $%d", argN)
		args = append(args, f.MaxPriceHour)
		argN++
	}
	if f.MinQtyPerson > 0 {
		query += fmt.Sprintf(" AND qty_person >= $%d", argN)
		args = append(args, f.MinQtyPerson)
		argN++
	}
	if f.HasAirConditioner != nil {
		query += fmt.Sprintf(" AND has_air_conditioner = $%d", argN)
		args = append(args, *f.HasAirConditioner)
		argN++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var rm models.Room
		if err := rows.Scan(&rm.ID, &rm.OfficeID, &rm.Type, &rm.PriceHour, &rm.QtyDesks, &rm.QtyPerson, &rm.AccessCode, &rm.HasAirConditioner, &rm.HasPrayerRoom); err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}
	return rooms, rows.Err()
}