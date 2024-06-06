package service

import (
	"LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"go.uber.org/zap"
)

type HistoryService interface {
	GetHistory(ctx context.Context, pagination domain.Pagination) ([]domain.History, error)
	DeleteOneHistory(ctx context.Context, postId int64, uid int64) error
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
		h.l.Error(constants.HistoryListError, zap.Error(err))
		return nil, err
	}

	return history, nil
}
func (h *historyService) DeleteOneHistory(ctx context.Context, postId int64, uid int64) error {
	if err := h.repo.DeleteOneHistory(ctx, postId, uid); err != nil {
		h.l.Error("delete one history failed", zap.Error(err))
		return err
	}
	return nil
}

func (h *historyService) DeleteAllHistory(ctx context.Context, uid int64) error {
	if err := h.repo.DeleteAllHistory(ctx, uid); err != nil {
		h.l.Error("delete all history failed", zap.Error(err))
		return err
	}
	return nil
}
