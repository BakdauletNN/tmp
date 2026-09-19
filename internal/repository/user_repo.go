package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"tmp/internal/models"
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
		INSERT INTO "users" (name, email, who, tg_chat_id, pass_hash, is_registered)
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id
	`
	return r.db.QueryRow(ctx, query, u.Name, u.Email, u.Who, u.TgChatID, u.PassHash, u.IsRegistered).Scan(&u.ID)
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, email, who, tg_chat_id, pass_hash, is_registered FROM "users" WHERE id = $1`
	var u models.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID,&u.Name, &u.Email, &u.Who, &u.TgChatID, &u.PassHash, &u.IsRegistered)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, mail string) (*models.User, error) {
	query := `SELECT id, name, email, who, tg_chat_id, pass_hash, is_registered FROM "users" WHERE email = $1`
	var u models.User
	err := r.db.QueryRow(ctx, query, mail).Scan(&u.ID, &u.Name, &u.Email, &u.Who, &u.TgChatID, &u.PassHash, &u.IsRegistered)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
