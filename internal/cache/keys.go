package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"tmp/internal/models"
)

func Room(id int) string {
	return fmt.Sprintf("cowork:room:%d", id)
}

func Office(id int) string {
	return fmt.Sprintf("cowork:office:%d", id)
}

func Rooms(filter *models.RoomFilter) (string, error) {
	return filterKey("rooms", filter)
}

func Offices(filter models.OfficeFilter) (string, error) {
	return filterKey("offices", filter)
}

func filterKey(prefix string, value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("cowork:%s:%x", prefix, sha256.Sum256(data)), nil
}
