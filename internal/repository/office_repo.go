package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"tmp/internal/models"
)

type OfficeRepository struct {
	Base
}

func NewOfficeRepository(db *pgxpool.Pool) *OfficeRepository {
	return &OfficeRepository{
		Base: NewBase(db, "offices"),
	}
}

func (r *OfficeRepository) GetOfficeInfo(ctx context.Context, id int) (*models.Office, error) {
	query := `
		SELECT 
			id, 
			name, 
			address, 
			star, 
			has_kitchen, 
			time_range_work, 
			metro_near
		FROM "offices" 
		WHERE id = $1
	`

	var o models.Office

	err := r.db.QueryRow(ctx, query, id).Scan(
		&o.ID,
		&o.Name,
		&o.Address,
		&o.Star,
		&o.HasKitchen,
		&o.TimeRangeWork,
		&o.MetroNear,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get office info: %w", err)
	}

	return &o, nil
}
