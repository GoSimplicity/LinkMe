package repository

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"go.uber.org/zap"
	"time"
)

type CheckRepository interface {
	Create(ctx context.Context, check domain.Check) (int64, error)                     // 创建审核记录
	UpdateStatus(ctx context.Context, check domain.Check) error                        // 更新审核状态
	FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) // 获取审核列表
	FindByID(ctx context.Context, checkID int64) (domain.Check, error)                 // 获取审核详情
	FindByPostId(ctx context.Context, postID uint) (domain.Check, error)               // 根据帖子ID获取审核信息
	GetCheckCount(ctx context.Context) (int64, error)                                  // 获取审核数量
}

type checkRepository struct {
	dao   dao.CheckDAO
	cache cache.CheckCache
	l     *zap.Logger
}

func NewCheckRepository(dao dao.CheckDAO, cache cache.CheckCache, l *zap.Logger) CheckRepository {
	return &checkRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (r *checkRepository) FindByPostId(ctx context.Context, postID uint) (domain.Check, error) {
	// 生成缓存键
	cacheKey := r.cache.GenerateCacheKey(postID)

	// 尝试从缓存中获取
	if cachedCheck, err := r.cache.GetCache(ctx, cacheKey); err == nil && cachedCheck != nil {
		r.l.Info("Cache hit for FindByPostId", zap.String("key", cacheKey))
		return *cachedCheck, nil
	}

	// 缓存未命中，从数据库获取
	check, err := r.dao.FindByPostId(ctx, postID)
	if err != nil {
		return domain.Check{}, err
	}

	// 将结果存入缓存
	r.cache.SetCache(ctx, cacheKey, toDomainCheck(check), 5*time.Minute)

	return toDomainCheck(check), nil
}

func (r *checkRepository) Create(ctx context.Context, check domain.Check) (int64, error) {
	// 先查找是否存在该帖子审核信息
	dc, err := r.dao.FindByPostId(ctx, check.PostID)
	if dc.PostID != 0 && err == nil {
		return -1, nil
	}

	// 创建新的审核信息
	id, err := r.dao.Create(ctx, toDAOCheck(check))
	if err != nil {
		return -1, err
	}

	// 清除与该帖子相关的缓存
	cacheKey := r.cache.GenerateCacheKey(check.PostID)
	r.cache.ClearCache(ctx, cacheKey)

	return id, nil
}

func (r *checkRepository) UpdateStatus(ctx context.Context, check domain.Check) error {
	// 更新数据库中的审核状态
	if err := r.dao.UpdateStatus(ctx, toDAOCheck(check)); err != nil {
		return err
	}

	// 清除与该帖子相关的缓存
	cacheKey := r.cache.GenerateCacheKey(check.PostID)
	r.cache.ClearCache(ctx, cacheKey)

	return nil
}

func (r *checkRepository) FindAll(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	// 可以为分页查询添加缓存，生成缓存键
	cacheKey := r.cache.GeneratePaginationCacheKey(pagination)

	// 尝试从缓存中获取
	if cachedChecks, err := r.cache.GetCacheList(ctx, cacheKey); err == nil && cachedChecks != nil {
		r.l.Info("Cache hit for FindAll", zap.String("key", cacheKey))
		return cachedChecks, nil
	}

	// 缓存未命中，从数据库获取
	checks, err := r.dao.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	// 将结果存入缓存
	r.cache.SetCacheList(ctx, cacheKey, toDomainChecks(checks), 5*time.Minute)

	return toDomainChecks(checks), nil
}

func (r *checkRepository) FindByID(ctx context.Context, checkID int64) (domain.Check, error) {
	// 生成缓存键
	cacheKey := r.cache.GenerateCacheKey(checkID)

	// 尝试从缓存中获取
	if cachedCheck, err := r.cache.GetCache(ctx, cacheKey); err == nil && cachedCheck != nil {
		r.l.Info("Cache hit for FindByID", zap.String("key", cacheKey))
		return *cachedCheck, nil
	}

	// 缓存未命中，从数据库获取
	check, err := r.dao.FindByID(ctx, checkID)
	if err != nil {
		return domain.Check{}, err
	}

	// 将结果存入缓存
	r.cache.SetCache(ctx, cacheKey, toDomainCheck(check), 5*time.Minute)

	return toDomainCheck(check), nil
}

func (r *checkRepository) GetCheckCount(ctx context.Context) (int64, error) {
	// 生成缓存键
	cacheKey := r.cache.GenerateCountCacheKey()

	// 尝试从缓存中获取
	if cachedCount, err := r.cache.GetCountCache(ctx, cacheKey); err == nil {
		r.l.Info("Cache hit for GetCheckCount", zap.String("key", cacheKey))
		return cachedCount, nil
	}

	// 缓存未命中，从数据库获取
	count, err := r.dao.GetCheckCount(ctx)
	if err != nil {
		return -1, err
	}

	// 将结果存入缓存
	r.cache.SetCountCache(ctx, cacheKey, count, 5*time.Minute)

	return count, nil
}

// toDAOCheck 将 domain.Check 转换为 dao.Check
func toDAOCheck(domainCheck domain.Check) dao.Check {
	return dao.Check{
		ID:        domainCheck.ID,
		PostID:    domainCheck.PostID,
		Content:   domainCheck.Content,
		Title:     domainCheck.Title,
		Author:    domainCheck.UserID,
		Status:    domainCheck.Status,
		Remark:    domainCheck.Remark,
		CreatedAt: domainCheck.CreatedAt,
		UpdatedAt: domainCheck.UpdatedAt,
	}
}

// toDomainCheck 将 dao.Check 转换为 domain.Check
func toDomainCheck(daoCheck dao.Check) domain.Check {
	return domain.Check{
		ID:        daoCheck.ID,
		PostID:    daoCheck.PostID,
		Content:   daoCheck.Content,
		Title:     daoCheck.Title,
		UserID:    daoCheck.Author,
		Status:    daoCheck.Status,
		Remark:    daoCheck.Remark,
		CreatedAt: daoCheck.CreatedAt,
		UpdatedAt: daoCheck.UpdatedAt,
	}
}

// toDomainChecks 将 []dao.Check 转换为 []domain.Check
func toDomainChecks(daoChecks []dao.Check) []domain.Check {
	domainChecks := make([]domain.Check, len(daoChecks))
	for i, daoCheck := range daoChecks {
		domainChecks[i] = toDomainCheck(daoCheck)
	}
	return domainChecks
}
