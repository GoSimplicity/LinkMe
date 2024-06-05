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
	SetHistory(ctx context.Context, post domain.Post, actionType string) error
	DeleteHistory(ctx context.Context, id int64, uid int64) error
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

func (h *historyService) SetHistory(ctx context.Context, post domain.Post, actionType string) error {
	return h.repo.SetHistory(ctx, post, actionType)
}

func (h *historyService) DeleteHistory(ctx context.Context, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}
