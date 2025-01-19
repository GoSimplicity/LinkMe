package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/publish"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
)

type CheckService interface {
	ApproveCheck(ctx context.Context, checkID int64, remark string, uid int64) error      // 审核通过
	RejectCheck(ctx context.Context, checkID int64, remark string, uid int64) error       // 审核拒绝
	ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) // 获取审核列表
	CheckDetail(ctx context.Context, checkID int64) (domain.Check, error)                 // 获取审核详情
	GetCheckCount(ctx context.Context) (int64, error)
}

type checkService struct {
	repo         repository.CheckRepository
	ActivityRepo repository.ActivityRepository
	producer     publish.Producer
	searchRepo   repository.SearchRepository
	l            *zap.Logger
}

func NewCheckService(repo repository.CheckRepository, searchRepo repository.SearchRepository, l *zap.Logger, ActivityRepo repository.ActivityRepository, producer publish.Producer) CheckService {
	return &checkService{
		repo:         repo,
		ActivityRepo: ActivityRepo,
		searchRepo:   searchRepo,
		l:            l,
		producer:     producer,
	}
}

// ApproveCheck 审核通过
func (s *checkService) ApproveCheck(ctx context.Context, checkID int64, remark string, uid int64) error {
	// 获取审核详情
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		return fmt.Errorf("获取审核详情失败: %w", err)
	}

	// 检查是否已审核
	if check.Status != domain.UnderReview {
		return fmt.Errorf("当前审核状态不允许操作,状态:%v", check.Status)
	}

	// 更新审核状态为通过
	if err := s.repo.UpdateStatus(ctx, domain.Check{
		ID:     checkID,
		Remark: remark,
		Status: domain.Approved,
	}); err != nil {
		s.l.Error("更新审核状态失败",
			zap.Int64("check_id", checkID),
			zap.String("remark", remark),
			zap.Error(err))
		return fmt.Errorf("更新审核状态失败: %w", err)
	}

	// 设置超时上下文
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// 异步处理非关键任务
	go func() {
		defer cancel()
		// 发送发布事件
		if err := s.producer.ProducePublishEvent(publish.PublishEvent{
			PostId: check.PostID,
			Uid:    check.Uid,
			Status: domain.Published,
		}); err != nil {
			s.l.Error("发送发布事件失败",
				zap.Uint("post_id", check.PostID),
				zap.Int64("uid", check.Uid),
				zap.Error(err))
		}

		// 记录最近活动
		if err := s.ActivityRepo.SetRecentActivity(ctxTimeout, domain.RecentActivity{
			UserID:      uid,
			Description: "审核通过",
			Time:        strconv.FormatInt(time.Now().Unix(), 10),
		}); err != nil {
			s.l.Error("记录最近活动失败",
				zap.Int64("user_id", uid),
				zap.Error(err))
		}
	}()

	return nil
}

// RejectCheck 审核拒绝
func (s *checkService) RejectCheck(ctx context.Context, checkID int64, remark string, uid int64) error {
	// 获取审核详情
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		return fmt.Errorf("获取审核详情失败: %w", err)
	}

	// 检查状态
	if check.Status != domain.UnderReview {
		return fmt.Errorf("当前审核状态不允许操作,状态:%v", check.Status)
	}

	// 使用事务处理状态更新
	err = s.repo.WithTx(ctx, func(txCtx context.Context) error {
		// 更新审核状态为拒绝
		if err := s.repo.UpdateStatus(txCtx, domain.Check{
			ID:     checkID,
			Remark: remark,
			Status: domain.UnApproved,
		}); err != nil {
			s.l.Error("更新审核状态失败",
				zap.Int64("check_id", checkID),
				zap.String("remark", remark),
				zap.Error(err))
			return fmt.Errorf("更新审核状态失败: %w", err)
		}

		// 异步发送审核拒绝事件
		if err := s.producer.ProducePublishEvent(publish.PublishEvent{
			PostId: check.PostID,
			Uid:    check.Uid,
			Status: domain.Draft, // 审核拒绝后,帖子状态为草稿
		}); err != nil {
			s.l.Error("发送审核拒绝事件失败", zap.Error(err))
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 设置超时上下文
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// 异步处理非关键任务
	go func() {
		defer cancel()
		// 记录最近活动
		if err := s.ActivityRepo.SetRecentActivity(ctxTimeout, domain.RecentActivity{
			UserID:      uid,
			Description: "审核拒绝",
			Time:        strconv.FormatInt(time.Now().Unix(), 10),
		}); err != nil {
			s.l.Error("记录最近活动失败",
				zap.Int64("user_id", uid),
				zap.Error(err))
		}
	}()

	return nil
}

// ListChecks 获取审核列表
func (s *checkService) ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset

	checks, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	s.l.Info("获取审核列表成功", zap.Int("数量", len(checks)))

	return checks, nil
}

// CheckDetail 获取审核详情
func (s *checkService) CheckDetail(ctx context.Context, checkID int64) (domain.Check, error) {
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		return domain.Check{}, err
	}

	s.l.Info("获取审核详情成功", zap.Int64("审核ID", checkID))

	return check, nil
}

// GetCheckCount 获取审核数量
func (s *checkService) GetCheckCount(ctx context.Context) (int64, error) {
	count, err := s.repo.GetCheckCount(ctx)
	if err != nil {
		return -1, err
	}

	return count, nil
}
