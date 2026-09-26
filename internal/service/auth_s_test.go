package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	"tmp/internal/models"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type fakeAuthRepository struct {
	user         *models.User
	lookupEmail  string
	resetHash    string
	resetExpires time.Time
	passwordHash string
	resetError   error
}

func (f *fakeAuthRepository) WriteUser(context.Context, *models.User) error { return nil }

func (f *fakeAuthRepository) GetByEmail(_ context.Context, email string) (*models.User, error) {
	f.lookupEmail = email
	if f.user == nil {
		return nil, pgx.ErrNoRows
	}
	return f.user, nil
}

func (f *fakeAuthRepository) CreatePasswordReset(_ context.Context, _ int, hash string, expires time.Time) error {
	f.resetHash = hash
	f.resetExpires = expires
	return f.resetError
}

func (f *fakeAuthRepository) ResetPassword(_ context.Context, hash, passwordHash string, _ time.Time) error {
	if f.resetError != nil {
		return f.resetError
	}
	f.resetHash = hash
	f.passwordHash = passwordHash
	return nil
}

type fakePasswordResetSender struct {
	to       string
	name     string
	resetURL string
	err      error
}

func (f *fakePasswordResetSender) SendPasswordReset(to, name, resetURL string) error {
	f.to, f.name, f.resetURL = to, name, resetURL
	return f.err
}

func TestRequestPasswordResetSendsHashedSingleUseToken(t *testing.T) {
	repo := &fakeAuthRepository{user: &models.User{ID: 7, Name: "Demo", Email: "demo@example.com"}}
	sender := &fakePasswordResetSender{}
	service := NewAuthService(repo, "test-secret")
	service.ConfigurePasswordReset(sender, "https://cowork.example")
	startedAt := time.Now()

	if err := service.RequestPasswordReset(context.Background(), repo.user.Email); err != nil {
		t.Fatalf("request reset: %v", err)
	}
	if sender.to != repo.user.Email || sender.name != repo.user.Name {
		t.Fatalf("unexpected recipient: %q %q", sender.to, sender.name)
	}
	parsed, err := url.Parse(sender.resetURL)
	if err != nil {
		t.Fatalf("parse reset URL: %v", err)
	}
	query, err := url.ParseQuery(strings.TrimPrefix(parsed.Fragment, "/reset-password?"))
	if err != nil {
		t.Fatalf("parse reset token: %v", err)
	}
	token := query.Get("token")
	if token == "" {
		t.Fatal("expected reset token in link")
	}
	hash := sha256.Sum256([]byte(token))
	if repo.resetHash != hex.EncodeToString(hash[:]) {
		t.Fatal("expected only the token hash to be stored")
	}
	if repo.resetExpires.Before(startedAt.Add(29*time.Minute)) || repo.resetExpires.After(startedAt.Add(31*time.Minute)) {
		t.Fatalf("unexpected token expiration: %s", repo.resetExpires)
	}
}

func TestRequestPasswordResetDoesNotRevealUnknownEmail(t *testing.T) {
	sender := &fakePasswordResetSender{}
	service := NewAuthService(&fakeAuthRepository{}, "test-secret")
	service.ConfigurePasswordReset(sender, "https://cowork.example")

	if err := service.RequestPasswordReset(context.Background(), "unknown@example.com"); err != nil {
		t.Fatalf("unknown account should receive a neutral success: %v", err)
	}
	if sender.resetURL != "" {
		t.Fatal("must not send a reset link for an unknown account")
	}
}

func TestRequestPasswordResetTrimsEmail(t *testing.T) {
	repo := &fakeAuthRepository{user: &models.User{ID: 7, Name: "Demo", Email: "demo@example.com"}}
	service := NewAuthService(repo, "test-secret")
	service.ConfigurePasswordReset(&fakePasswordResetSender{}, "https://cowork.example")

	if err := service.RequestPasswordReset(context.Background(), "  demo@example.com  "); err != nil {
		t.Fatalf("request reset: %v", err)
	}
	if repo.lookupEmail != "demo@example.com" {
		t.Fatalf("expected trimmed email lookup, got %q", repo.lookupEmail)
	}
}

func TestResetPasswordStoresNewHashAndRejectsInvalidToken(t *testing.T) {
	repo := &fakeAuthRepository{}
	service := NewAuthService(repo, "test-secret")
	token := strings.Repeat("x", 43)

	if err := service.ResetPassword(context.Background(), token, "new-password"); err != nil {
		t.Fatalf("reset password: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.passwordHash), []byte("new-password")); err != nil {
		t.Fatalf("new password was not hashed: %v", err)
	}
	if !errors.Is(service.ResetPassword(context.Background(), "bad", "new-password"), ErrInvalidPasswordReset) {
		t.Fatal("expected malformed token to be rejected")
	}
}
