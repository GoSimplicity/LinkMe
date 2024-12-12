package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
)

type ApiRepository interface {
	CreateApi(ctx context.Context, api *domain.Api) error
	GetApiById(ctx context.Context, id int) (*domain.Api, error)
	UpdateApi(ctx context.Context, api *domain.Api) error
	DeleteApi(ctx context.Context, id int) error
	ListApis(ctx context.Context, page, pageSize int) ([]*domain.Api, int, error)
}

type apiRepository struct {
	l   *zap.Logger
	dao dao.ApiDAO
}

func NewApiRepository(l *zap.Logger, dao dao.ApiDAO) ApiRepository {
	return &apiRepository{
		l:   l,
		dao: dao,
	}
}

// CreateApi 创建新的API
func (a *apiRepository) CreateApi(ctx context.Context, api *domain.Api) error {
	if api == nil {
		return dao.ErrInvalidMenu // 添加空指针检查
	}
	return a.dao.CreateApi(ctx, a.apiToDAO(api))
}

// GetApiById 根据ID获取API
func (a *apiRepository) GetApiById(ctx context.Context, id int) (*domain.Api, error) {
	if id <= 0 {
		return nil, dao.ErrInvalidMenu // 添加ID有效性检查
	}
	api, err := a.dao.GetApiById(ctx, id)
	if err != nil {
		return nil, err
	}
	return a.apiFromDAO(api), nil
}

// UpdateApi 更新API信息
func (a *apiRepository) UpdateApi(ctx context.Context, api *domain.Api) error {
	if api == nil {
		return dao.ErrInvalidMenu
	}
	return a.dao.UpdateApi(ctx, a.apiToDAO(api))
}

// DeleteApi 删除API
func (a *apiRepository) DeleteApi(ctx context.Context, id int) error {
	if id <= 0 {
		return dao.ErrInvalidMenu
	}
	return a.dao.DeleteApi(ctx, id)
}

// ListApis 分页获取API列表
func (a *apiRepository) ListApis(ctx context.Context, page, pageSize int) ([]*domain.Api, int, error) {
	// 参数校验
	if page <= 0 || pageSize <= 0 {
		return nil, 0, dao.ErrInvalidMenu
	}

	apis, total, err := a.dao.ListApis(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return a.apisFromDAO(apis), total, nil
}

// apiToDAO 将domain层的API对象转换为DAO层对象
func (a *apiRepository) apiToDAO(api *domain.Api) *dao.Api {
	if api == nil {
		return nil
	}
	return &dao.Api{
		ID:          api.ID,
		Name:        api.Name,
		Path:        api.Path,
		Method:      api.Method,
		Description: api.Description,
		Version:     api.Version,
		Category:    api.Category,
		IsPublic:    api.IsPublic,
		CreateTime:  api.CreateTime,
		UpdateTime:  api.UpdateTime,
		IsDeleted:   api.IsDeleted,
	}
}

// apiFromDAO 将DAO层的API对象转换为domain层对象
func (a *apiRepository) apiFromDAO(api *dao.Api) *domain.Api {
	if api == nil {
		return nil
	}
	return &domain.Api{
		ID:          api.ID,
		Name:        api.Name,
		Path:        api.Path,
		Method:      api.Method,
		Description: api.Description,
		Version:     api.Version,
		Category:    api.Category,
		IsPublic:    api.IsPublic,
		CreateTime:  api.CreateTime,
		UpdateTime:  api.UpdateTime,
		IsDeleted:   api.IsDeleted,
	}
}

// apisFromDAO 将DAO层的API对象列表转换为domain层对象列表
func (a *apiRepository) apisFromDAO(apis []*dao.Api) []*domain.Api {
	if apis == nil {
		return nil
	}
	result := make([]*domain.Api, 0, len(apis))
	for _, api := range apis {
		if api != nil {
			result = append(result, a.apiFromDAO(api))
		}
	}
	return result
}
