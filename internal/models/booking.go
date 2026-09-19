package models

import (
	"time"
	
)

type Booking struct {
	ID        int       `json:"id" db:"id"`
	RoomID    int       `json:"room_id" db:"room_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	DeskID    *int      `json:"desk_id,omitempty" db:"desk_id"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
}

type CreateBookingInput struct {
	DeskID    *int      `json:"desk_id,omitempty" example:"3"`
	RoomID    int       `json:"room_id" validate:"required,gt=0" example:"1"`
	StartTime time.Time `json:"start_time" validate:"required" example:"2026-09-20T10:00:00Z"`
	EndTime   time.Time `json:"end_time" validate:"required" example:"2026-09-20T12:00:00Z"`
}
func (i CreateBookingInput) Validate() error {
	return validate.Struct(i)
}