package tests

import (
	"testing"

	"tmp/internal/models"
)

func TestSignUpValidation(t *testing.T) {
	t.Run("valid signup passes", func(t *testing.T) {
		payload := models.SignUp{
			Name:     "Alice",
			Email:    "alice@example.com",
			Password: "secret123",
		}

		if err := payload.Validate(); err != nil {
			t.Fatalf("expected valid signup to pass validation, got error: %v", err)
		}
	})

	t.Run("invalid signup fails", func(t *testing.T) {
		cases := []models.SignUp{
			{Name: "A", Email: "alice@example.com", Password: "secret123"},
			{Name: "Alice", Email: "not-an-email", Password: "secret123"},
			{Name: "Alice", Email: "alice@example.com", Password: "short"},
		}

		for _, tc := range cases {
			if err := tc.Validate(); err == nil {
				t.Fatalf("expected validation error for signup: %+v", tc)
			}
		}
	})
}

func TestSignInValidation(t *testing.T) {
	t.Run("valid signin passes", func(t *testing.T) {
		payload := models.SignIn{
			Email:    "alice@example.com",
			Password: "secret123",
		}

		if err := payload.Validate(); err != nil {
			t.Fatalf("expected valid signin to pass validation, got error: %v", err)
		}
	})

	t.Run("invalid signin fails", func(t *testing.T) {
		cases := []models.SignIn{
			{Email: "not-an-email", Password: "secret123"},
			{Email: "alice@example.com", Password: "short"},
			{Email: "", Password: "secret123"},
		}

		for _, tc := range cases {
			if err := tc.Validate(); err == nil {
				t.Fatalf("expected validation error for signin: %+v", tc)
			}
		}
	})
}
