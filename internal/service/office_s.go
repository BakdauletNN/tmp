package service

import (
	"context"

	"tmp/internal/models"
	"tmp/internal/repository"
)

type OfficeService struct {
	repo *repository.OfficeRepository
}

func NewOfficeService(repo *repository.OfficeRepository) *OfficeService {
	return &OfficeService{repo: repo}
}

func (s *OfficeService) GetOffice(ctx context.Context, id int) (*models.Office, error) {
	return s.repo.GetOfficeInfo(ctx, id)
}

func (s *OfficeService) SearchByAddress(ctx context.Context, f models.OfficeFilter) ([]models.Office, error) {
	return s.repo.SearchOffices(ctx, f)
}