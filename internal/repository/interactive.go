package repository

import (
	. "LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

type InteractiveRepository interface {
	// BatchIncrReadCnt biz 和 bizId 长度必须一致
	BatchIncrReadCnt(ctx context.Context, biz []string, ids []int64) error                    // 批量更新阅读计数(与kafka配合使用)
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error                      // 增加阅读计数
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error                      // 减少阅读计数
	IncrCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error // 收藏
	DecrCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error // 取消收藏
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)                // 获取互动信息
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)                 // 检查是否已点赞
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)             // 检查是否被收藏
	GetById(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)       // 批量获取互动信息
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

func (c *CachedInteractiveRepository) BatchIncrReadCnt(ctx context.Context, biz []string, ids []int64) error {
	if err := c.dao.BatchIncrReadCnt(ctx, biz, ids); err != nil {
		c.l.Error("batch incr read count failed", zap.Error(err))
		return err
	}
	// 使用sync.WaitGroup来等待所有缓存更新操作完成
	var wg sync.WaitGroup
	wg.Add(len(biz))
	for i := 0; i < len(biz); i++ {
		go func(i int) {
			defer wg.Done()
			er := c.retryUpdateCache(ctx, biz[i], ids[i], 3)
			if er != nil {
				c.l.Error("post read count record failed", zap.Error(er))
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func (c *CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	if err := c.dao.InsertLikeInfo(ctx, models.UserLikeBiz{
		BizName: biz,
		BizID:   id,
		Uid:     uid,
	}); err != nil {
		c.l.Error("incr like failed", zap.Error(err))
		return err
	}
	return c.cache.PostReadCountRecord(ctx, biz, id)
}

func (c *CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	if err := c.dao.DeleteLikeInfo(ctx, models.UserLikeBiz{
		BizName: biz,
		BizID:   id,
		Uid:     uid,
	}); err != nil {
		c.l.Error("decr like failed", zap.Error(err))
		return err
	}
	return c.cache.DecrLikeCountRecord(ctx, biz, id)
}

func (c *CachedInteractiveRepository) IncrCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	if err := c.dao.InsertCollectionBiz(ctx, models.UserCollectionBiz{
		BizName:      biz,
		BizID:        id,
		CollectionId: cid,
		Uid:          uid,
	}); err != nil {
		c.l.Error("incr collection item failed", zap.Error(err))
		return err
	}
	return c.cache.PostCollectCountRecord(ctx, biz, id)
}
func (c *CachedInteractiveRepository) DecrCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	if err := c.dao.DeleteCollectionBiz(ctx, models.UserCollectionBiz{
		BizName:      biz,
		BizID:        id,
		CollectionId: cid,
		Uid:          uid,
	}); err != nil {
		c.l.Error("decr collection item failed", zap.Error(err))
		return err
	}
	return c.cache.DecrCollectCountRecord(ctx, biz, id)
}

func (c *CachedInteractiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	inc, err := c.cache.Get(ctx, biz, id)
	if err == nil {
		return inc, nil
	}
	ic, err := c.dao.Get(ctx, biz, id)
	if err != nil {
		c.l.Error(PostGetInteractiveERROR, zap.Error(err))
		return domain.Interactive{}, err
	}
	if er := c.cache.Set(ctx, biz, id, toDomain(ic)); er != nil {
		c.l.Error("set interactive cache failed", zap.Error(er))
	}
	return toDomain(ic), nil
}

func (c *CachedInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
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

func (c *CachedInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, id, uid)
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

func (c *CachedInteractiveRepository) GetById(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	ics, err := c.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		c.l.Error("get interactions failed", zap.Error(err))
		return make([]domain.Interactive, 0), err
	}
	result := make([]domain.Interactive, len(ics))
	for i, ic := range ics {
		result[i] = toDomain(ic)
	}

	return result, nil
}

func (c *CachedInteractiveRepository) retryUpdateCache(ctx context.Context, biz string, id int64, retries int) error {
	for i := 0; i < retries; i++ {
		err := c.cache.PostReadCountRecord(ctx, biz, id)
		if err == nil {
			return nil // 更新成功，返回
		}
		// 如果更新失败，等待一段时间后再重试
		time.Sleep(time.Millisecond * 100) // 等待100毫秒
	}
	return errors.New("failed to update cache after retries")
}

func toDomain(ic models.Interactive) domain.Interactive {
	return domain.Interactive{
		BizID:        ic.BizID,
		ReadCount:    ic.ReadCount,
		LikeCount:    ic.LikeCount,
		CollectCount: ic.CollectCount,
	}
}
