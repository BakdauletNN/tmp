package models

type Office struct {
	ID            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Address       string `json:"address" db:"address"`
	Star          int    `json:"star" db:"star"`
	HasKitchen    bool   `json:"has_kitchen" db:"has_kitchen"`
	TimeRangeWork string `json:"time_range_work" db:"time_range_work"`
	MetroNear     bool   `json:"metro_near" db:"metro_near"`
}
