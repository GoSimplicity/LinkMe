package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PlateService interface {
	CreatePlate(ctx *gin.Context, plate domain.Plate) error
	ListPlate(ctx *gin.Context, pagination domain.Pagination) ([]domain.Plate, error)
	UpdatePlate(ctx *gin.Context, plateId int64, uid int64) error
	DeletePlate(ctx *gin.Context, plateId int64, uid int64) error
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

func (p *plateService) CreatePlate(ctx *gin.Context, plate domain.Plate) error {
	return p.repo.CreatePlate(ctx, plate)
}

func (p *plateService) ListPlate(ctx *gin.Context, pagination domain.Pagination) ([]domain.Plate, error) {
	//TODO implement me
	panic("implement me")
}

func (p *plateService) UpdatePlate(ctx *gin.Context, plateId int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *plateService) DeletePlate(ctx *gin.Context, plateId int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}
