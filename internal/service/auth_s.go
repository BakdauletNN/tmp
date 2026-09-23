package service

import (
	"context"
	"errors"
	"time"

	"tmp/internal/models"
	"tmp/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


type AuthService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}


func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name: name,
		Email:    email,
		PassHash: string(hash),
	}

	err = s.repo.WriteUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}


func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	claims := jwt.MapClaims{
	"user_id": user.ID,
	"role":    user.Who,
	"exp":     time.Now().Add(time.Hour * 24).Unix(),
}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
