package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"tmp/internal/models"
	"tmp/internal/repository"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		PassHash: string(hash),
	}

	err = s.repo.WriteUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
