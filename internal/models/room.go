package models

type Room struct {
	ID                int     `json:"id" db:"id"`
	OfficeID          int     `json:"office_id" db:"office_id"`
	Type              string  `json:"type" db:"type"`
	PriceHour         float64 `json:"price_hour" db:"price_hour"`
	QtyDesks          int     `json:"qty_desks" db:"qty_desks"`
	QtyPerson         int     `json:"qty_person" db:"qty_person"`
	AccessCode        *int    `json:"access_code" db:"access_code"`
	HasAirConditioner bool    `json:"has_air_conditioner" db:"has_air_conditioner"`
	HasPrayerRoom     bool    `json:"has_prayer_room" db:"has_prayer_room"`
}
