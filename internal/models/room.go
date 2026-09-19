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

type RoomFilter struct {
	OfficeID          int     `form:"office_id"`
	Type              string  `form:"type"`
	MaxPriceHour      float64 `form:"max_price_hour"`
	MinQtyPerson      int     `form:"min_qty_person"`
	HasAirConditioner *bool   `form:"has_air_conditioner"`
}
