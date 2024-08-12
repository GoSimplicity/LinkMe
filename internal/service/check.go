package service

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type CheckService interface {
	SubmitCheck(ctx context.Context, check domain.Check) (int64, error)                   // 提交审核
	ApproveCheck(ctx context.Context, checkID int64, remark string, uid int64) error      // 审核通过
	RejectCheck(ctx context.Context, checkID int64, remark string, uid int64) error       // 审核拒绝
	ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) // 获取审核列表
	CheckDetail(ctx context.Context, checkID int64) (domain.Check, error)                 // 获取审核详情
	GetCheckCount(ctx context.Context) (int64, error)
}

type checkService struct {
	repo         repository.CheckRepository
	ActivityRepo repository.ActivityRepository
	postRepo     repository.PostRepository
	historyRepo  repository.HistoryRepository
	searchRepo   repository.SearchRepository
	l            *zap.Logger
}

func NewCheckService(repo repository.CheckRepository, postRepo repository.PostRepository, historyRepo repository.HistoryRepository, searchRepo repository.SearchRepository, l *zap.Logger, ActivityRepo repository.ActivityRepository) CheckService {
	return &checkService{
		repo:         repo,
		postRepo:     postRepo,
		historyRepo:  historyRepo,
		ActivityRepo: ActivityRepo,
		searchRepo:   searchRepo,
		l:            l,
	}
}

func (s *checkService) SubmitCheck(ctx context.Context, check domain.Check) (int64, error) {
	// 设置状态为审核中
	check.Status = domain.UnderReview
	id, err := s.repo.Create(ctx, check)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *checkService) ApproveCheck(ctx context.Context, checkID int64, remark string, uid int64) error {
	// 获取审核详情
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		return fmt.Errorf("获取审核详情失败: %w", err)
	}

	// 检查是否已审核
	if check.Status == domain.UnApproved || check.Status == domain.Approved {
		return fmt.Errorf("请勿重复提交：%v", checkID)
	}

	// 更新审核状态为通过
	if err := s.repo.UpdateStatus(ctx, domain.Check{
		ID:     checkID,
		Remark: remark,
		Status: domain.Approved,
	}); err != nil {
		s.l.Error("更新审核状态失败", zap.Int64("CheckID", checkID), zap.String("Remark", remark), zap.Error(err))
		return fmt.Errorf("更新审核状态失败: %w", err)
	}

	// 获取相关的帖子
	post, err := s.postRepo.GetPostById(ctx, check.PostID, check.UserID)
	if err != nil {
		s.l.Error("获取帖子失败", zap.Uint("PostID", check.PostID), zap.Int64("UserID", check.UserID), zap.Error(err))
		return fmt.Errorf("获取帖子失败: %w", err)
	}

	// 更新帖子状态为已发布并同步
	post.Status = domain.Published
	if err := s.postRepo.UpdateStatus(ctx, post); err != nil {
		s.l.Error("更新帖子状态失败", zap.Uint("PostID", post.ID), zap.Error(err))
		return fmt.Errorf("更新帖子状态失败: %w", err)
	}

	// 存入历史记录
	if err := s.historyRepo.SetHistory(ctx, post); err != nil {
		s.l.Error("存入历史记录失败", zap.Uint("PostID", post.ID), zap.Error(err))
	}

	// 添加搜索索引
	if err := s.searchRepo.InputPost(ctx, domain.PostSearch{
		Id:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Status:  post.Status,
	}); err != nil {
		s.l.Error("添加搜索索引失败", zap.Uint("PostID", post.ID), zap.Error(err))
	}

	// 异步记录最近活动
	go func() {
		if err := s.ActivityRepo.SetRecentActivity(context.Background(), domain.RecentActivity{
			UserID:      uid,
			Description: "审核通过",
			Time:        strconv.FormatInt(time.Now().Unix(), 10),
		}); err != nil {
			s.l.Error("记录最近活动失败", zap.Int64("UserID", uid), zap.Error(err))
		}
	}()

	s.l.Info("审核通过并发布帖子", zap.Uint("PostID", post.ID))

	return nil
}

func (s *checkService) RejectCheck(ctx context.Context, checkID int64, remark string, uid int64) error {
	// 获取审核详情
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		return fmt.Errorf("get check detail failed: %w", err)
	}
	// 检查状态

	if check.Status == domain.UnApproved || check.Status == domain.Approved {
		return fmt.Errorf("请勿重复提交：%v", checkID)
	}

	// 更新审核状态为拒绝
	err = s.repo.UpdateStatus(ctx, domain.Check{
		ID:     checkID,
		Remark: remark,
		Status: domain.UnApproved,
	})
	if err != nil {
		return fmt.Errorf("update check status failed: %w", err)
	}

	s.l.Info("Rejected check", zap.Int64("CheckID", checkID), zap.String("Remark", remark))

	go func() {
		er := s.ActivityRepo.SetRecentActivity(context.Background(), domain.RecentActivity{
			UserID:      uid,
			Description: "审核拒绝",
			Time:        strconv.FormatInt(time.Now().Unix(), 10),
		})
		if er != nil {
			s.l.Error("Failed to set recent activity", zap.Int64("UserID", uid), zap.Error(er))
		}
	}()

	return nil
}

func (s *checkService) ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.Check, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset

	checks, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	s.l.Info("Listed checks", zap.Int("Count", len(checks)))

	return checks, nil
}

func (s *checkService) CheckDetail(ctx context.Context, checkID int64) (domain.Check, error) {
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		return domain.Check{}, err
	}

	s.l.Info("Fetched check detail", zap.Int64("CheckID", checkID))

	return check, nil
}

func (s *checkService) GetCheckCount(ctx context.Context) (int64, error) {
	count, err := s.repo.GetCheckCount(ctx)
	if err != nil {
		return -1, err
	}

	return count, nil
}
