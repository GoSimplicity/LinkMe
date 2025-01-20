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
	ApproveCheck(ctx context.Context, checkID int64, remark string, uid int64) error
	RejectCheck(ctx context.Context, checkID int64, remark string, uid int64) error
	ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error)
	CheckDetail(ctx context.Context, checkID int64) (domain.Check, error)
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
		s.l.Error("获取审核详情失败", zap.Error(err))
		return fmt.Errorf("获取审核详情失败: %w", err)
	}

	// 检查是否已审核
	if check.Status != domain.UnderReview {
		s.l.Error("当前审核状态不允许操作", zap.Uint8("status", check.Status))
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

	// 使用errgroup并发处理异步任务
	go func() {
		// 创建带超时的context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 并发执行发布事件和记录活动
		done := make(chan error, 2)
		go func() {
			done <- s.producer.ProducePublishEvent(publish.PublishEvent{
				PostId: check.PostID,
				Uid:    check.Uid,
				Status: domain.Published,
			})
		}()

		go func() {
			done <- s.recordActivity(uid, "审核通过")
		}()

		// 等待所有goroutine完成或超时
		for i := 0; i < 2; i++ {
			select {
			case err := <-done:
				if err != nil {
					s.l.Error("异步任务执行失败", zap.Error(err))
				}
			case <-ctx.Done():
				s.l.Error("异步任务执行超时")
				return
			}
		}
	}()

	return nil
}

// RejectCheck 审核拒绝
func (s *checkService) RejectCheck(ctx context.Context, checkID int64, remark string, uid int64) error {
	// 获取审核详情
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		s.l.Error("获取审核详情失败", zap.Error(err))
		return fmt.Errorf("获取审核详情失败: %w", err)
	}

	// 检查状态
	if check.Status != domain.UnderReview {
		s.l.Error("当前审核状态不允许操作", zap.Uint8("status", check.Status))
		return fmt.Errorf("当前审核状态不允许操作,状态:%v", check.Status)
	}

	if err := s.repo.UpdateStatus(ctx, domain.Check{
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

	// 使用errgroup并发处理异步任务
	go func() {
		// 创建带超时的context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 并发执行发布事件和记录活动
		done := make(chan error, 2)
		go func() {
			done <- s.producer.ProducePublishEvent(publish.PublishEvent{
				PostId: check.PostID,
				Uid:    check.Uid,
				Status: domain.Draft,
			})
		}()

		go func() {
			done <- s.recordActivity(uid, "审核拒绝")
		}()

		// 等待所有goroutine完成或超时
		for i := 0; i < 2; i++ {
			select {
			case err := <-done:
				if err != nil {
					s.l.Error("异步任务执行失败", zap.Error(err))
				}
			case <-ctx.Done():
				s.l.Error("异步任务执行超时")
				return
			}
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
		s.l.Error("获取审核列表失败", zap.Error(err))
		return nil, err
	}

	return checks, nil
}

// CheckDetail 获取审核详情
func (s *checkService) CheckDetail(ctx context.Context, checkID int64) (domain.Check, error) {
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		s.l.Error("获取审核详情失败", zap.Error(err))
		return domain.Check{}, err
	}

	return check, nil
}

// recordActivity 记录活动
func (s *checkService) recordActivity(uid int64, desc string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.ActivityRepo.SetRecentActivity(ctx, domain.RecentActivity{
		UserID:      uid,
		Description: desc,
		Time:        strconv.FormatInt(time.Now().Unix(), 10),
	}); err != nil {
		s.l.Error("记录最近活动失败",
			zap.Int64("user_id", uid),
			zap.Error(err))
		return err
	}

	return nil
}
