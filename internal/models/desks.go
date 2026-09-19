package models

type Desk struct {
	ID     int    `json:"id" db:"id"`
	RoomID int    `json:"room_id" db:"room_id"`
	Label  string `json:"label" db:"label"`
}