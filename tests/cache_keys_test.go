package tests

import (
	"testing"

	"tmp/internal/cache"
	"tmp/internal/models"
)

func TestRoomCacheKeyIsStableAndUnique(t *testing.T) {
	filter := &models.RoomFilter{
		OfficeID:     12,
		Type:         "private",
		MaxPriceHour: 25,
		MinQtyPerson: 3,
	}

	first, err := cache.Rooms(filter)
	if err != nil {
		t.Fatalf("expected room cache key, got error: %v", err)
	}

	second, err := cache.Rooms(&models.RoomFilter{
		OfficeID:     12,
		Type:         "private",
		MaxPriceHour: 25,
		MinQtyPerson: 3,
	})
	if err != nil {
		t.Fatalf("expected same room filter to generate key, got error: %v", err)
	}

	if first != second {
		t.Fatalf("expected stable key for same filter, got %q and %q", first, second)
	}

	third, err := cache.Rooms(&models.RoomFilter{
		OfficeID:     12,
		Type:         "shared",
		MaxPriceHour: 25,
		MinQtyPerson: 3,
	})
	if err != nil {
		t.Fatalf("expected alternate room filter key, got error: %v", err)
	}

	if first == third {
		t.Fatal("expected different room filters to produce different cache keys")
	}
}

func TestOfficeCacheKeyIsStableAndUnique(t *testing.T) {
	address := "Almaty, Central Avenue"
	first, err := cache.Offices(models.OfficeFilter{
		Address:       address,
		MinStar:       4,
		HasKitchen:    boolPtr(true),
		MetroNear:     boolPtr(true),
		TimeRangeWork: stringPtr("09:00-21:00"),
	})
	if err != nil {
		t.Fatalf("expected office cache key, got error: %v", err)
	}

	second, err := cache.Offices(models.OfficeFilter{
		Address:       address,
		MinStar:       4,
		HasKitchen:    boolPtr(true),
		MetroNear:     boolPtr(true),
		TimeRangeWork: stringPtr("09:00-21:00"),
	})
	if err != nil {
		t.Fatalf("expected same office filter to generate key, got error: %v", err)
	}

	if first != second {
		t.Fatalf("expected stable key for same office filter, got %q and %q", first, second)
	}

	third, err := cache.Offices(models.OfficeFilter{
		Address:       address,
		MinStar:       5,
		HasKitchen:    boolPtr(true),
		MetroNear:     boolPtr(true),
		TimeRangeWork: stringPtr("09:00-21:00"),
	})
	if err != nil {
		t.Fatalf("expected alternate office filter key, got error: %v", err)
	}

	if first == third {
		t.Fatal("expected different office filters to produce different cache keys")
	}
}
