package service

import (
	"context"
	"errors"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"

	"go.uber.org/zap"
)

type InteractiveService interface {
	Like(ctx context.Context, postId uint, uid int64) error                            // 点赞
	CancelLike(ctx context.Context, postId uint, uid int64) error                      // 取消点赞
	Collect(ctx context.Context, postId uint, uid int64) error                         // 收藏
	CancelCollect(ctx context.Context, postId uint, uid int64) error                   // 取消收藏
	Get(ctx context.Context, postId uint) (domain.Interactive, error)                  // 获取互动信息
	GetByIds(ctx context.Context, postIds []uint) (map[uint]domain.Interactive, error) // 批量获取互动信息(热榜算法需要)
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

func (i *interactiveService) Like(ctx context.Context, postId uint, uid int64) error {
	// 检查是否已点赞,避免重复操作
	exist, err := i.repo.Liked(ctx, postId, uid)
	if err != nil && !errors.Is(err, dao.ErrRecordNotFound) {
		i.l.Error("查询点赞状态失败", zap.Error(err))
		return err
	}

	if exist {
		return errors.New("请勿重复点赞")
	}

	// 增加点赞计数
	if err := i.repo.IncrLike(ctx, postId, uid); err != nil {
		i.l.Error("点赞失败", zap.Error(err))
		return err
	}

	return nil
}

func (i *interactiveService) CancelLike(ctx context.Context, postId uint, uid int64) error {
	return i.repo.DecrLike(ctx, postId, uid)
}

func (i *interactiveService) Collect(ctx context.Context, postId uint, uid int64) error {
	collected, err := i.repo.Collected(ctx, postId, uid)
	if err != nil && !errors.Is(err, dao.ErrRecordNotFound) {
		i.l.Error("查询收藏状态失败", zap.Error(err))
		return err
	}

	if collected {
		return errors.New("请勿重复收藏")
	}

	return i.repo.IncrCollectionItem(ctx, postId, uid)
}

func (i *interactiveService) CancelCollect(ctx context.Context, postId uint, uid int64) error {
	return i.repo.DecrCollectionItem(ctx, postId, uid)
}

func (i *interactiveService) Get(ctx context.Context, postId uint) (domain.Interactive, error) {
	di, err := i.repo.Get(ctx, postId)
	if err != nil {
		return domain.Interactive{}, err
	}

	return di, err
}

func (i *interactiveService) GetByIds(ctx context.Context, postIds []uint) (map[uint]domain.Interactive, error) {
	dis, err := i.repo.GetById(ctx, postIds)
	if err != nil {
		return nil, err
	}
	resultDis := make(map[uint]domain.Interactive)

	for _, interactive := range dis {
		resultDis[interactive.BizID] = interactive
	}

	return resultDis, err
}
