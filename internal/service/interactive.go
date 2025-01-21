package service

import (
	"context"
	"errors"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"

	"go.uber.org/zap"
)

// InteractiveService 定义互动相关的业务接口
type InteractiveService interface {
	Like(ctx context.Context, postId uint, uid int64) error
	CancelLike(ctx context.Context, postId uint, uid int64) error
	Collect(ctx context.Context, postId uint, uid int64) error
	CancelCollect(ctx context.Context, postId uint, uid int64) error
	Get(ctx context.Context, postId uint) (domain.Interactive, error)
	GetByIds(ctx context.Context, postIds []uint) (map[uint]domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
	l    *zap.Logger
}

func NewInteractiveService(repo repository.InteractiveRepository, l *zap.Logger) InteractiveService {
	return &interactiveService{
		repo: repo,
		l:    l,
	}
}

// Like 处理点赞逻辑
func (i *interactiveService) Like(ctx context.Context, postId uint, uid int64) error {
	if postId == 0 || uid <= 0 {
		return errors.New("invalid parameters")
	}

	exist, err := i.repo.Liked(ctx, postId, uid)
	if err != nil && !errors.Is(err, dao.ErrRecordNotFound) {
		i.l.Error("点赞状态查询失败", zap.Error(err), zap.Uint("postId", postId), zap.Int64("uid", uid))
		return err
	}

	if exist {
		return errors.New("已点赞")
	}

	return i.repo.IncrLike(ctx, postId, uid)
}

// CancelLike 处理取消点赞逻辑
func (i *interactiveService) CancelLike(ctx context.Context, postId uint, uid int64) error {
	if postId == 0 || uid <= 0 {
		return errors.New("invalid parameters")
	}
	return i.repo.DecrLike(ctx, postId, uid)
}

// Collect 处理收藏逻辑
func (i *interactiveService) Collect(ctx context.Context, postId uint, uid int64) error {
	if postId == 0 || uid <= 0 {
		return errors.New("invalid parameters")
	}

	collected, err := i.repo.Collected(ctx, postId, uid)
	if err != nil && !errors.Is(err, dao.ErrRecordNotFound) {
		i.l.Error("收藏状态查询失败", zap.Error(err), zap.Uint("postId", postId), zap.Int64("uid", uid))
		return err
	}

	if collected {
		return errors.New("已收藏")
	}

	return i.repo.IncrCollectionItem(ctx, postId, uid)
}

// CancelCollect 处理取消收藏逻辑
func (i *interactiveService) CancelCollect(ctx context.Context, postId uint, uid int64) error {
	if postId == 0 || uid <= 0 {
		return errors.New("invalid parameters")
	}
	return i.repo.DecrCollectionItem(ctx, postId, uid)
}

// Get 获取单个互动信息
func (i *interactiveService) Get(ctx context.Context, postId uint) (domain.Interactive, error) {
	if postId == 0 {
		return domain.Interactive{}, errors.New("invalid parameters")
	}
	return i.repo.Get(ctx, postId)
}

// GetByIds 批量获取互动信息
func (i *interactiveService) GetByIds(ctx context.Context, postIds []uint) (map[uint]domain.Interactive, error) {
	if len(postIds) == 0 {
		return nil, errors.New("invalid parameters")
	}

	interactions, err := i.repo.GetById(ctx, postIds)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]domain.Interactive, len(interactions))
	for _, interaction := range interactions {
		result[interaction.BizID] = interaction
	}

	return result, nil
}
