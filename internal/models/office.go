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


type OfficeFilter struct {
	Address    string `form:"address"`
	MinStar    int    `form:"min_star"`
	HasKitchen *bool  `form:"has_kitchen"`
	MetroNear  *bool  `form:"metro_near"`
	TimeRangeWork *string `form:"time_range_work"`
}