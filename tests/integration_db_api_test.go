package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"tmp/internal/database"
	handlerHTTP "tmp/internal/handler/http"
	"tmp/internal/models"
	"tmp/internal/repository"
	"tmp/internal/service"
)

func setupIntegrationDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "postgresql://cowork:cowork@localhost:5432/cowork?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := database.Connect(dsn)
	if err != nil {
		t.Skipf("postgres not available for integration tests: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		t.Skipf("postgres is not ready for integration tests: %v", err)
	}

	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			who VARCHAR(255) NOT NULL,
			tg_chat_id BIGINT,
			pass_hash VARCHAR(255) NOT NULL,
			is_registered BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)
	if err != nil {
		db.Close()
		t.Fatalf("failed to initialize users table: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestUserRepository_Integration(t *testing.T) {
	db := setupIntegrationDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	email := fmt.Sprintf("integration-%d@example.com", time.Now().UnixNano())
	user := &models.User{
		Name:         "Integration User",
		Email:        email,
		Who:          "customer",
		PassHash:     "hashed-password",
		IsRegistered: true,
	}

	if err := repo.WriteUser(ctx, user); err != nil {
		t.Fatalf("expected DB insert to succeed, got: %v", err)
	}

	if user.ID == 0 {
		t.Fatal("expected user ID to be generated after insert")
	}

	stored, err := repo.GetByEmail(ctx, email)
	if err != nil {
		t.Fatalf("expected user to be fetched from database, got: %v", err)
	}

	if stored.Email != email {
		t.Fatalf("expected email %q, got %q", email, stored.Email)
	}

	if stored.Name != user.Name {
		t.Fatalf("expected name %q, got %q", user.Name, stored.Name)
	}

	_, err = db.Exec(ctx, `DELETE FROM users WHERE email = $1`, email)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}

func TestAuthAPI_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupIntegrationDB(t)
	repo := repository.NewUserRepository(db)
	authService := service.NewAuthService(repo, "super-secret-key")
	handler := handlerHTTP.NewAuthHandler(authService)

	router := gin.New()
	handler.RegisterRoutes(router)

	email := fmt.Sprintf("api-%d@example.com", time.Now().UnixNano())
	registerBody := fmt.Sprintf(`{"name":"API User","email":"%s","password":"supersecure123"}`, email)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(registerBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("register endpoint expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var registerResp struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &registerResp); err != nil {
		t.Fatalf("invalid register response: %v", err)
	}

	if registerResp.Email != email {
		t.Fatalf("expected registered email %q, got %q", email, registerResp.Email)
	}

	loginBody := fmt.Sprintf(`{"email":"%s","password":"supersecure123"}`, email)
	loginRec := httptest.NewRecorder()
	loginReq := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(loginRec, loginReq)

	if loginRec.Code != http.StatusOK {
		t.Fatalf("login endpoint expected 200, got %d: %s", loginRec.Code, loginRec.Body.String())
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("invalid login response: %v", err)
	}

	if strings.TrimSpace(loginResp.Token) == "" {
		t.Fatal("expected JWT token to be returned after login")
	}

	_, err := db.Exec(context.Background(), `DELETE FROM users WHERE email = $1`, email)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}
