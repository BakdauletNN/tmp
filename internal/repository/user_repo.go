package repository

import (
	"context"
	"strings"
	"time"

	"tmp/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Base
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		Base: NewBase(db, "users"),
	}
}

func (r *UserRepository) WriteUser(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO "users" (name, email, who, pass_hash, is_registered)
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`
	return r.db.QueryRow(ctx, query, u.Name, u.Email, u.Who, u.PassHash, u.IsRegistered).Scan(&u.ID)
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, email, who, pass_hash, is_registered FROM "users" WHERE id = $1`
	var u models.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Who, &u.PassHash, &u.IsRegistered)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, mail string) (*models.User, error) {
	query := `SELECT id, name, email, who, pass_hash, is_registered FROM "users" WHERE lower(email) = lower($1)`
	var u models.User
	err := r.db.QueryRow(ctx, query, strings.TrimSpace(mail)).Scan(&u.ID, &u.Name, &u.Email, &u.Who, &u.PassHash, &u.IsRegistered)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) CreatePasswordReset(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `DELETE FROM password_reset_tokens WHERE user_id = $1`, userID); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM password_reset_tokens WHERE expires_at <= now()`); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *UserRepository) ResetPassword(ctx context.Context, tokenHash, passwordHash string, now time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var userID int
	err = tx.QueryRow(ctx, `
		SELECT user_id FROM password_reset_tokens
		WHERE token_hash = $1 AND expires_at > $2
		FOR UPDATE
	`, tokenHash, now).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return pgx.ErrNoRows
		}
		return err
	}
	if _, err = tx.Exec(ctx, `UPDATE users SET pass_hash = $1 WHERE id = $2`, passwordHash, userID); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM password_reset_tokens WHERE token_hash = $1`, tokenHash); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
