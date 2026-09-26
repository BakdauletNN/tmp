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


func (r *OfficeRepository) SearchOffices(ctx context.Context, f models.OfficeFilter) ([]models.Office, error) {
	query := `
		SELECT id, name, address, star, has_kitchen, time_range_work, metro_near
		FROM "offices"
		WHERE 1=1
	`
	var args []interface{}
	argN := 1

	if f.Address != "" {
		query += fmt.Sprintf(" AND address ILIKE '%%' || $%d || '%%'", argN)
		args = append(args, f.Address)
		argN++
	}
	if f.MinStar > 0 {
		query += fmt.Sprintf(" AND star >= $%d", argN)
		args = append(args, f.MinStar)
		argN++
	}
	if f.HasKitchen != nil {
		query += fmt.Sprintf(" AND has_kitchen = $%d", argN)
		args = append(args, *f.HasKitchen)
		argN++
	}
	if f.MetroNear != nil {
		query += fmt.Sprintf(" AND metro_near = $%d", argN)
		args = append(args, *f.MetroNear)
		argN++
	}
	if f.TimeRangeWork != nil {
		query += fmt.Sprintf(" AND time_range_work = $%d", argN)
		args = append(args, *f.TimeRangeWork)
		argN++
}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offices []models.Office
	for rows.Next() {
		var o models.Office
		if err := rows.Scan(&o.ID, &o.Name, &o.Address, &o.Star, &o.HasKitchen, &o.TimeRangeWork, &o.MetroNear); err != nil {
			return nil, err
		}
		offices = append(offices, o)
	}
	return offices, rows.Err()
}