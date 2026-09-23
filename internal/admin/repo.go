package admin

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"tmp/internal/models"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) AllOffices(ctx context.Context) ([]models.Office, error) {
	query := `SELECT id, name, address, star, has_kitchen, time_range_work, metro_near, owner_id FROM offices`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offices []models.Office
	for rows.Next() {
		var o models.Office
		if err := rows.Scan(&o.ID, &o.Name, &o.Address, &o.Star, &o.HasKitchen, &o.TimeRangeWork, &o.MetroNear, &o.OwnerID); err != nil {
			return nil, err
		}
		offices = append(offices, o)
	}
	return offices, rows.Err()
}

func (r *AdminRepository) AllBookings(ctx context.Context) ([]models.Booking, error) {
	query := `SELECT id, room_id, user_id, desk_id, start_time, end_time FROM bookings ORDER BY start_time`
	rows, err := r.db.Query(ctx, query)
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

func (r *AdminRepository) AllUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, email, who, tg_chat_id, is_registered FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Who, &u.TgChatID, &u.IsRegistered); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}


