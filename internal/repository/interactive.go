package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"

	"go.uber.org/zap"
)

type InteractiveRepository interface {
	BatchIncrReadCnt(ctx context.Context, postIds []uint) error
	IncrReadCnt(ctx context.Context, postId uint) error
	IncrLike(ctx context.Context, postId uint, uid int64) error
	DecrLike(ctx context.Context, postId uint, uid int64) error
	IncrCollectionItem(ctx context.Context, postId uint, cid int64, uid int64) error
	DecrCollectionItem(ctx context.Context, postId uint, cid int64, uid int64) error
	Get(ctx context.Context, postId uint) (domain.Interactive, error)
	Liked(ctx context.Context, postId uint, uid int64) (bool, error)
	Collected(ctx context.Context, postId uint, uid int64) (bool, error)
	GetById(ctx context.Context, postIds []uint) ([]domain.Interactive, error)
}

type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	l     *zap.Logger
	cache cache.InteractiveCache
}

func NewInteractiveRepository(dao dao.InteractiveDAO, l *zap.Logger, cache cache.InteractiveCache) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao:   dao,
		l:     l,
		cache: cache,
	}
}

// BatchIncrReadCnt 批量增加阅读计数
func (c *CachedInteractiveRepository) BatchIncrReadCnt(ctx context.Context, postIds []uint) error {
	if err := c.dao.BatchIncrReadCnt(ctx, postIds); err != nil {
		return err
	}
	// 使用sync.WaitGroup来等待所有缓存更新操作完成
	var wg sync.WaitGroup
	wg.Add(len(postIds))
	for i := 0; i < len(postIds); i++ {
		go func(i int) {
			defer wg.Done()
			er := c.retryUpdateCache(ctx, postIds[i], 3)
			if er != nil {
				c.l.Error("post read count record failed", zap.Error(er))
			}
		}(i)
	}
	wg.Wait()
	return nil
}

// IncrReadCnt implements InteractiveRepository.
func (c *CachedInteractiveRepository) IncrReadCnt(ctx context.Context, postId uint) error {
	if err := c.dao.IncrReadCnt(ctx, postId); err != nil {
		return err
	}
	return c.cache.PostReadCountRecord(ctx, postId)
}

// IncrLike 增加点赞计数
func (c *CachedInteractiveRepository) IncrLike(ctx context.Context, postId uint, uid int64) error {
	if err := c.dao.InsertLikeInfo(ctx, dao.UserLikeBiz{
		BizID: postId,
		Uid:   uid,
	}); err != nil {
		return err
	}
	return c.cache.PostReadCountRecord(ctx, postId)
}

// DecrLike 减少点赞计数
func (c *CachedInteractiveRepository) DecrLike(ctx context.Context, postId uint, uid int64) error {
	if err := c.dao.DeleteLikeInfo(ctx, dao.UserLikeBiz{
		BizID: postId,
		Uid:   uid,
	}); err != nil {
		return err
	}
	return c.cache.DecrLikeCountRecord(ctx, postId)
}

// IncrCollectionItem 增加收藏计数
func (c *CachedInteractiveRepository) IncrCollectionItem(ctx context.Context, postId uint, cid int64, uid int64) error {
	if err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		BizID:        postId,
		CollectionId: cid,
		Uid:          uid,
	}); err != nil {
		return err
	}
	return c.cache.PostCollectCountRecord(ctx, postId)
}

// DecrCollectionItem 减少收藏计数
func (c *CachedInteractiveRepository) DecrCollectionItem(ctx context.Context, postId uint, cid int64, uid int64) error {
	if err := c.dao.DeleteCollectionBiz(ctx, dao.UserCollectionBiz{
		BizID:        postId,
		CollectionId: cid,
		Uid:          uid,
	}); err != nil {
		return err
	}
	return c.cache.DecrCollectCountRecord(ctx, postId)
}

// Get 获取互动信息
func (c *CachedInteractiveRepository) Get(ctx context.Context, postId uint) (domain.Interactive, error) {
	inc, err := c.cache.Get(ctx, postId)
	if err == nil {
		return inc, nil
	}
	ic, err := c.dao.Get(ctx, postId)
	if err != nil {
		c.l.Error(PostGetInteractiveERROR, zap.Error(err))
		return domain.Interactive{}, err
	}
	if er := c.cache.Set(ctx, postId, toDomain(ic)); er != nil {
		c.l.Error("set interactive cache failed", zap.Error(er))
	}
	return toDomain(ic), nil
}

// Liked 检查是否已点赞
func (c *CachedInteractiveRepository) Liked(ctx context.Context, postId uint, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, postId, uid)
	switch {
	case err == nil:
		// 如果没有错误，说明找到了点赞记录，返回true
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		// 如果错误是ErrRecordNotFound，说明没有找到点赞记录
		c.l.Error(PostGetLikedERROR, zap.Error(err))
		return false, nil
	default:
		// 如果是其他错误，返回false和错误信息
		return false, err
	}
}

// Collected 检查是否已收藏
func (c *CachedInteractiveRepository) Collected(ctx context.Context, postId uint, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, postId, uid)
	switch {
	case err == nil:
		// 如果没有错误，说明找到了收藏记录，返回true
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		// 如果错误是ErrRecordNotFound，说明没有找到收藏记录
		c.l.Error(PostGetCollectERROR, zap.Error(err))
		return false, nil
	default:
		// 如果是其他错误，返回false和错误信息
		return false, err
	}
}

// GetById 批量获取互动信息
func (c *CachedInteractiveRepository) GetById(ctx context.Context, postIds []uint) ([]domain.Interactive, error) {
	ics, err := c.dao.GetByIds(ctx, postIds)
	if err != nil {
		return make([]domain.Interactive, 0), err
	}
	result := make([]domain.Interactive, len(ics))
	for i, ic := range ics {
		result[i] = toDomain(ic)
	}

	return result, nil
}

// retryUpdateCache 重试更新缓存
func (c *CachedInteractiveRepository) retryUpdateCache(ctx context.Context, postId uint, retries int) error {
	for i := 0; i < retries; i++ {
		err := c.cache.PostReadCountRecord(ctx, postId)
		if err == nil {
			return nil // 更新成功，返回
		}
		// 如果更新失败，等待一段时间后再重试
		time.Sleep(time.Millisecond * 100) // 等待100毫秒
	}
	return errors.New("failed to update cache after retries")
}

func toDomain(ic dao.Interactive) domain.Interactive {
	return domain.Interactive{
		BizID:        ic.BizID,
		ReadCount:    ic.ReadCount,
		LikeCount:    ic.LikeCount,
		CollectCount: ic.CollectCount,
	}
}
