package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"tmp/internal/logger"
	"tmp/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPasswordReset = errors.New("password reset token is invalid or expired")

type authRepository interface {
	WriteUser(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	CreatePasswordReset(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	ResetPassword(ctx context.Context, tokenHash, passwordHash string, now time.Time) error
}

type PasswordResetSender interface {
	SendPasswordReset(to, name, resetURL string) error
}

type AuthService struct {
	repo                 authRepository
	jwtSecret            string
	passwordResetSender  PasswordResetSender
	passwordResetBaseURL string
}

func NewAuthService(repo authRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) ConfigurePasswordReset(sender PasswordResetSender, webBaseURL string) {
	s.passwordResetSender = sender
	s.passwordResetBaseURL = strings.TrimRight(webBaseURL, "/")
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	log := logger.FromContext(ctx)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("password hash failed", slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		PassHash: string(hash),
	}

	if err := s.repo.WriteUser(ctx, user); err != nil {
		log.Error("write user failed", slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("user registered in service", slog.Int("user_id", user.ID), slog.String("email", email))
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	log := logger.FromContext(ctx)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		log.Warn("login user not found", slog.String("email", email), slog.String("error", err.Error()))
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		log.Warn("invalid password attempt", slog.Int("user_id", user.ID), slog.String("email", email))
		return "", errors.New("invalid password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Who,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Error("jwt signing failed", slog.Int("user_id", user.ID), slog.String("error", err.Error()))
		return "", err
	}

	log.Info("token generated", slog.Int("user_id", user.ID), slog.String("role", user.Who))
	return tokenStr, nil
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	log := logger.FromContext(ctx)
	user, err := s.repo.GetByEmail(ctx, strings.TrimSpace(email))
	if errors.Is(err, pgx.ErrNoRows) {
		log.Info("password reset skipped: account not found")
		return nil
	}
	if err != nil {
		return err
	}
	if s.passwordResetSender == nil || s.passwordResetBaseURL == "" {
		return errors.New("password reset sender is not configured")
	}

	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return fmt.Errorf("generate password reset token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(rawToken)
	tokenHash := hashResetToken(token)
	expiresAt := time.Now().Add(30 * time.Minute)
	if err := s.repo.CreatePasswordReset(ctx, user.ID, tokenHash, expiresAt); err != nil {
		return err
	}

	resetURL := s.passwordResetBaseURL + "/#/reset-password?token=" + url.QueryEscape(token)
	if err := s.passwordResetSender.SendPasswordReset(user.Email, user.Name, resetURL); err != nil {
		log.Error("password reset message delivery failed", slog.String("error", err.Error()))
		return err
	}
	log.Info("password reset email accepted by SMTP")
	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, password string) error {
	if len(token) < 40 || len(token) > 64 {
		return ErrInvalidPasswordReset
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.repo.ResetPassword(ctx, hashResetToken(token), string(passwordHash), time.Now())
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidPasswordReset
	}
	return err
}

func hashResetToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
