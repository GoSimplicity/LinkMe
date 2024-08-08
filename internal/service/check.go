package service

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/constants"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type CheckService interface {
	SubmitCheck(ctx context.Context, check domain.Check) (int64, error)                       // 提交审核
	ApproveCheck(ctx context.Context, checkID int64, remark string, uid int64) error          // 审核通过
	RejectCheck(ctx context.Context, checkID int64, remark string, uid int64) error           // 审核拒绝
	ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.CheckList, error) // 获取审核列表
	CheckDetail(ctx context.Context, checkID int64) (domain.Check, error)                     // 获取审核详情
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
	id, err := s.repo.Create(ctx, check)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *checkService) ApproveCheck(ctx context.Context, checkID int64, remark string, uid int64) error {
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		return fmt.Errorf("get check detail failed: %w", err)
	}
	if check.Status == constants.PostUnApproved || check.Status == constants.PostApproved {
		return fmt.Errorf("请勿重复提交：%v", checkID)
	}
	// 更新审核状态为通过
	err = s.repo.UpdateStatus(ctx, domain.Check{
		ID:     checkID,
		Remark: remark,
		Status: constants.PostApproved,
	})
	if err != nil {
		s.l.Error("Failed to update check status", zap.Int64("CheckID", checkID), zap.String("Remark", remark), zap.Error(err))
		return fmt.Errorf("update check status failed: %w", err)
	}
	s.l.Info("Approved check", zap.Int64("CheckID", checkID), zap.String("Remark", remark))
	// 获取审核详情
	check, err = s.repo.FindByID(ctx, checkID)
	if err != nil {
		s.l.Error("Failed to get check detail", zap.Int64("CheckID", checkID), zap.Error(err))
		return fmt.Errorf("get check detail failed: %w", err)
	}
	// 获取相关的帖子
	post, err := s.postRepo.GetPostById(ctx, check.PostID, check.UserID)
	if err != nil {
		s.l.Error("Failed to get post", zap.Uint("PostID", check.PostID), zap.Int64("UserID", check.UserID), zap.Error(err))
		return fmt.Errorf("get post failed: %w", err)
	}
	// 更新帖子状态为已发布并同步
	post.Status = domain.Published
	if er := s.postRepo.UpdateStatus(ctx, post); er != nil {
		s.l.Error("Failed to update post status", zap.Uint("PostID", post.ID), zap.Error(er))
		return fmt.Errorf("update post status failed: %w", er)
	}
	// 存入历史记录
	if er := s.historyRepo.SetHistory(ctx, post); er != nil {
		s.l.Error("Set history failed", zap.Uint("PostID", post.ID), zap.Error(er))
	}
	s.l.Info("Post has been published", zap.Uint("PostID", post.ID))
	// 添加搜索索引
	err = s.searchRepo.InputPost(ctx, domain.PostSearch{
		Id:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Status:  post.Status,
	})
	if err != nil {
		s.l.Error("Add search index failed", zap.Uint("PostID", post.ID), zap.Error(err))
	}
	go func() {
		er := s.ActivityRepo.SetRecentActivity(context.Background(), domain.RecentActivity{
			UserID:      uid,
			Description: "审核通过",
			Time:        strconv.FormatInt(time.Now().Unix(), 10),
		})
		if er != nil {
			s.l.Error("Failed to set recent activity", zap.Int64("UserID", uid), zap.Error(er))
		}
	}()
	return nil
}

func (s *checkService) RejectCheck(ctx context.Context, checkID int64, remark string, uid int64) error {
	// 获取审核详情
	check, err := s.repo.FindByID(ctx, checkID)
	if err != nil {
		return fmt.Errorf("get check detail failed: %w", err)
	}
	// 检查状态
	if check.Status == constants.PostUnApproved || check.Status == constants.PostApproved {
		return fmt.Errorf("请勿重复提交：%v", checkID)
	}
	// 更新审核状态为拒绝
	err = s.repo.UpdateStatus(ctx, domain.Check{
		ID:     checkID,
		Remark: remark,
		Status: constants.PostUnApproved,
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

func (s *checkService) ListChecks(ctx context.Context, pagination domain.Pagination) ([]domain.CheckList, error) {
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
