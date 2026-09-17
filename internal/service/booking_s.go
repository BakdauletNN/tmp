package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"tmp/internal/models"
	"tmp/internal/repository"
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) Create_Booking(ctx context.Context, ) (*models.User, error) {
	//TODO
}