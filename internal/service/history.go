package service

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
)

type HistoryService interface {
	GetHistory(ctx context.Context, pagination domain.Pagination) ([]domain.History, error)
	DeleteOneHistory(ctx context.Context, postId uint, uid int64) error
	DeleteAllHistory(ctx context.Context, uid int64) error
}

type historyService struct {
	repo repository.HistoryRepository
	l    *zap.Logger
}

func NewHistoryService(repo repository.HistoryRepository, l *zap.Logger) HistoryService {
	return &historyService{
		repo: repo,
		l:    l,
	}
}

func (h *historyService) GetHistory(ctx context.Context, pagination domain.Pagination) ([]domain.History, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	history, err := h.repo.GetHistory(ctx, pagination)
	if err != nil {
		return nil, err
	}

	return history, nil
}
func (h *historyService) DeleteOneHistory(ctx context.Context, postId uint, uid int64) error {
	if err := h.repo.DeleteOneHistory(ctx, postId, uid); err != nil {
		return err
	}
	return nil
}

func (h *historyService) DeleteAllHistory(ctx context.Context, uid int64) error {
	if err := h.repo.DeleteAllHistory(ctx, uid); err != nil {
		return err
	}
	return nil
}
