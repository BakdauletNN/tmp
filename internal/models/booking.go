package models

import "time"

type Booking struct {
	ID        int       `json:"id" db:"id"`
	RoomID    int       `json:"room_id" db:"room_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
}
