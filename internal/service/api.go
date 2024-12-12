package service

import (
	"context"
	"errors"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
)

type ApiService interface {
	CreateApi(ctx context.Context, api *domain.Api) error
	GetApiById(ctx context.Context, id int) (*domain.Api, error)
	UpdateApi(ctx context.Context, api *domain.Api) error
	DeleteApi(ctx context.Context, id int) error
	ListApis(ctx context.Context, page, pageSize int) ([]*domain.Api, int, error)
}

type apiService struct {
	l    *zap.Logger
	repo repository.ApiRepository
}

func NewApiService(l *zap.Logger, repo repository.ApiRepository) ApiService {
	return &apiService{
		l:    l,
		repo: repo,
	}
}

// CreateApi 创建新的API
func (a *apiService) CreateApi(ctx context.Context, api *domain.Api) error {
	if api == nil {
		a.l.Warn("API不能为空")
		return errors.New("api不能为空")
	}

	return a.repo.CreateApi(ctx, api)
}

// GetApiById 根据ID获取API
func (a *apiService) GetApiById(ctx context.Context, id int) (*domain.Api, error) {
	if id <= 0 {
		a.l.Warn("API ID无效", zap.Int("ID", id))
		return nil, errors.New("api id无效")
	}

	return a.repo.GetApiById(ctx, id)
}

// UpdateApi 更新API信息
func (a *apiService) UpdateApi(ctx context.Context, api *domain.Api) error {
	if api == nil {
		a.l.Warn("API不能为空")
		return errors.New("api不能为空")
	}

	return a.repo.UpdateApi(ctx, api)
}

// DeleteApi 删除指定ID的API
func (a *apiService) DeleteApi(ctx context.Context, id int) error {
	if id <= 0 {
		a.l.Warn("API ID无效", zap.Int("ID", id))
		return errors.New("api id无效")
	}

	return a.repo.DeleteApi(ctx, id)
}

// ListApis 分页获取API列表
func (a *apiService) ListApis(ctx context.Context, page, pageSize int) ([]*domain.Api, int, error) {
	if page < 1 || pageSize < 1 {
		a.l.Warn("分页参数无效", zap.Int("页码", page), zap.Int("每页数量", pageSize))
		return nil, 0, errors.New("分页参数无效")
	}

	return a.repo.ListApis(ctx, page, pageSize)
}
