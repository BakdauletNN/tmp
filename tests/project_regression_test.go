package tests

import (
	"testing"

	"tmp/internal/cache"
	"tmp/internal/models"
)

func TestProjectAuthRegression(t *testing.T) {
	t.Run("signup validation accepts real user data", func(t *testing.T) {
		payload := models.SignUp{
			Name:     "Alice Johnson",
			Email:    "alice@coworkgo.com",
			Password: "supersecure123",
		}

		if err := payload.Validate(); err != nil {
			t.Fatalf("expected valid signup to pass: %v", err)
		}
	})

	t.Run("signup rejects weak or malformed values", func(t *testing.T) {
		cases := []models.SignUp{
			{Name: "A", Email: "alice@coworkgo.com", Password: "supersecure123"},
			{Name: "Alice Johnson", Email: "not-an-email", Password: "supersecure123"},
			{Name: "Alice Johnson", Email: "alice@coworkgo.com", Password: "short"},
		}

		for _, tc := range cases {
			if err := tc.Validate(); err == nil {
				t.Fatalf("expected signup validation to fail for %+v", tc)
			}
		}
	})

	t.Run("signin validation accepts correct credentials shape", func(t *testing.T) {
		payload := models.SignIn{
			Email:    "alice@coworkgo.com",
			Password: "supersecure123",
		}

		if err := payload.Validate(); err != nil {
			t.Fatalf("expected valid signin to pass: %v", err)
		}
	})

	t.Run("signin rejects malformed credentials", func(t *testing.T) {
		cases := []models.SignIn{
			{Email: "bad-email", Password: "supersecure123"},
			{Email: "alice@coworkgo.com", Password: "short"},
			{Email: "", Password: "supersecure123"},
		}

		for _, tc := range cases {
			if err := tc.Validate(); err == nil {
				t.Fatalf("expected signin validation to fail for %+v", tc)
			}
		}
	})
}

func TestProjectCacheRegression(t *testing.T) {
	t.Run("room search keys are deterministic and scoped by filter", func(t *testing.T) {
		first, err := cache.Rooms(&models.RoomFilter{OfficeID: 7, Type: "private", MaxPriceHour: 30, MinQtyPerson: 2})
		if err != nil {
			t.Fatalf("expected room filter key: %v", err)
		}

		second, err := cache.Rooms(&models.RoomFilter{OfficeID: 7, Type: "private", MaxPriceHour: 30, MinQtyPerson: 2})
		if err != nil {
			t.Fatalf("expected same room filter key: %v", err)
		}

		if first != second {
			t.Fatalf("room keys must be stable for same filter: %q != %q", first, second)
		}

		third, err := cache.Rooms(&models.RoomFilter{OfficeID: 7, Type: "open", MaxPriceHour: 30, MinQtyPerson: 2})
		if err != nil {
			t.Fatalf("expected alternate room filter key: %v", err)
		}

		if first == third {
			t.Fatal("room keys must change when filter changes")
		}
	})

	t.Run("office search keys are deterministic and ignore pointer identity", func(t *testing.T) {
		hasKitchen := true
		metroNear := true
		workHours := "09:00-21:00"

		first, err := cache.Offices(models.OfficeFilter{
			Address:       "Almaty, Central Ave",
			MinStar:       4,
			HasKitchen:    &hasKitchen,
			MetroNear:     &metroNear,
			TimeRangeWork: &workHours,
		})
		if err != nil {
			t.Fatalf("expected office filter key: %v", err)
		}

		second, err := cache.Offices(models.OfficeFilter{
			Address:       "Almaty, Central Ave",
			MinStar:       4,
			HasKitchen:    boolPtr(true),
			MetroNear:     boolPtr(true),
			TimeRangeWork: stringPtr("09:00-21:00"),
		})
		if err != nil {
			t.Fatalf("expected same office filter key: %v", err)
		}

		if first != second {
			t.Fatalf("office keys must be stable for same filter: %q != %q", first, second)
		}

		third, err := cache.Offices(models.OfficeFilter{
			Address:       "Almaty, Central Ave",
			MinStar:       5,
			HasKitchen:    &hasKitchen,
			MetroNear:     &metroNear,
			TimeRangeWork: &workHours,
		})
		if err != nil {
			t.Fatalf("expected changed office filter key: %v", err)
		}

		if first == third {
			t.Fatal("office keys must change when search parameters change")
		}
	})
}

func boolPtr(v bool) *bool       { return &v }
func stringPtr(v string) *string { return &v }
