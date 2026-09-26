package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Base struct {
	db        *pgxpool.Pool
	tableName string
}

func NewBase(db *pgxpool.Pool, tableName string) Base {
	return Base{
		db:        db,
		tableName: tableName,
	}
}

func (r *Base) Delete(ctx context.Context, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", pgx.Identifier{r.tableName}.Sanitize())

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete from %s: %w", r.tableName, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *Base) Exists(ctx context.Context, id int) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", pgx.Identifier{r.tableName}.Sanitize())

	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existence in %s: %w", r.tableName, err)
	}

	return exists, nil
}

// Count returns the total number of records
func (r *Base) Count(ctx context.Context) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", pgx.Identifier{r.tableName}.Sanitize())

	var count int
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count rows in %s: %w", r.tableName, err)
	}

	return count, nil
}
