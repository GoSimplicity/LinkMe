package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"go.uber.org/zap"
)

type PlateService interface {
	CreatePlate(ctx context.Context, plate domain.Plate) error
	ListPlate(ctx context.Context, pagination domain.Pagination) ([]domain.Plate, error)
	UpdatePlate(ctx context.Context, plate domain.Plate) error
	DeletePlate(ctx context.Context, plateId int64, uid int64) error
}

type plateService struct {
	l    *zap.Logger
	repo repository.PlateRepository
}

func NewPlateService(l *zap.Logger, repo repository.PlateRepository) PlateService {
	return &plateService{
		l:    l,
		repo: repo,
	}
}

func (p *plateService) CreatePlate(ctx context.Context, plate domain.Plate) error {
	return p.repo.CreatePlate(ctx, plate)
}

func (p *plateService) ListPlate(ctx context.Context, pagination domain.Pagination) ([]domain.Plate, error) {
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	plates, err := p.repo.ListPlate(ctx, pagination)
	if err != nil {
		p.l.Error("failed to list plate", zap.Error(err))
		return nil, err
	}
	return plates, err
}

func (p *plateService) UpdatePlate(ctx context.Context, plate domain.Plate) error {
	return p.repo.UpdatePlate(ctx, plate)
}

func (p *plateService) DeletePlate(ctx context.Context, plateId int64, uid int64) error {
	return p.repo.DeletePlate(ctx, plateId, uid)
}
