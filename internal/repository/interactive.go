package repository

import (
	"context"
	"errors"

	. "github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"

	"go.uber.org/zap"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, postId uint) error
	IncrLike(ctx context.Context, postId uint, uid int64) error
	DecrLike(ctx context.Context, postId uint, uid int64) error
	IncrCollectionItem(ctx context.Context, postId uint, uid int64) error
	DecrCollectionItem(ctx context.Context, postId uint, uid int64) error
	Get(ctx context.Context, postId uint) (domain.Interactive, error)
	Liked(ctx context.Context, postId uint, uid int64) (bool, error)
	Collected(ctx context.Context, postId uint, uid int64) (bool, error)
	GetById(ctx context.Context, postIds []uint) ([]domain.Interactive, error)
}

type InteractiveRepositoryImpl struct {
	dao dao.InteractiveDAO
	l   *zap.Logger
}

func NewInteractiveRepository(dao dao.InteractiveDAO, l *zap.Logger) InteractiveRepository {
	return &InteractiveRepositoryImpl{
		dao: dao,
		l:   l,
	}
}

// IncrReadCnt 增加阅读计数
func (i *InteractiveRepositoryImpl) IncrReadCnt(ctx context.Context, postId uint) error {
	return i.dao.IncrReadCnt(ctx, postId)
}

// IncrLike 增加点赞计数
func (i *InteractiveRepositoryImpl) IncrLike(ctx context.Context, postId uint, uid int64) error {
	return i.dao.InsertLikeInfo(ctx, dao.UserLikeBiz{
		BizID: postId,
		Uid:   uid,
	})
}

// DecrLike 减少点赞计数
func (i *InteractiveRepositoryImpl) DecrLike(ctx context.Context, postId uint, uid int64) error {
	return i.dao.DeleteLikeInfo(ctx, dao.UserLikeBiz{
		BizID: postId,
		Uid:   uid,
	})
}

// IncrCollectionItem 增加收藏计数
func (i *InteractiveRepositoryImpl) IncrCollectionItem(ctx context.Context, postId uint, uid int64) error {
	return i.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		BizID: postId,
		Uid:   uid,
	})
}

// DecrCollectionItem 减少收藏计数
func (i *InteractiveRepositoryImpl) DecrCollectionItem(ctx context.Context, postId uint, uid int64) error {
	return i.dao.DeleteCollectionBiz(ctx, dao.UserCollectionBiz{
		BizID: postId,
		Uid:   uid,
	})
}

// Get 获取互动信息
func (i *InteractiveRepositoryImpl) Get(ctx context.Context, postId uint) (domain.Interactive, error) {
	ic, err := i.dao.Get(ctx, postId)
	if err != nil {
		i.l.Error(PostGetInteractiveERROR, zap.Error(err))
		return domain.Interactive{}, err
	}
	return toDomain(ic), nil
}

// Liked 检查是否已点赞
func (i *InteractiveRepositoryImpl) Liked(ctx context.Context, postId uint, uid int64) (bool, error) {
	_, err := i.dao.GetLikeInfo(ctx, postId, uid)
	return i.checkExistence(err, PostGetLikedERROR)
}

// Collected 检查是否已收藏
func (i *InteractiveRepositoryImpl) Collected(ctx context.Context, postId uint, uid int64) (bool, error) {
	_, err := i.dao.GetCollectInfo(ctx, postId, uid)
	return i.checkExistence(err, PostGetCollectERROR)
}

// GetById 批量获取互动信息
func (i *InteractiveRepositoryImpl) GetById(ctx context.Context, postIds []uint) ([]domain.Interactive, error) {
	ics, err := i.dao.GetByIds(ctx, postIds)
	if err != nil {
		return make([]domain.Interactive, 0), err
	}

	result := make([]domain.Interactive, len(ics))
	for idx, ic := range ics {
		result[idx] = toDomain(ic)
	}
	return result, nil
}

// checkExistence 检查是否存在
func (i *InteractiveRepositoryImpl) checkExistence(err error, logMsg string) (bool, error) {
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		i.l.Error(logMsg, zap.Error(err))
		return false, nil
	default:
		return false, err
	}
}

func toDomain(ic dao.Interactive) domain.Interactive {
	return domain.Interactive{
		BizID:        ic.BizID,
		ReadCount:    ic.ReadCount,
		LikeCount:    ic.LikeCount,
		CollectCount: ic.CollectCount,
	}
}
