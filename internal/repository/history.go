package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"go.uber.org/zap"
)

type HistoryRepository interface {
	GetHistory(ctx context.Context, pagination domain.Pagination) ([]domain.History, error)
	SetHistory(ctx context.Context, post domain.Post, actionType string) error
}

type historyRepository struct {
	l   *zap.Logger
	dao dao.HistoryDAO
}

func NewHistoryRepository(l *zap.Logger, dao dao.HistoryDAO) HistoryRepository {
	return &historyRepository{
		l:   l,
		dao: dao,
	}
}

func (h *historyRepository) GetHistory(ctx context.Context, pagination domain.Pagination) ([]domain.History, error) {
	record, err := h.dao.GetHistoryRecord(ctx, pagination)
	if err != nil {
		h.l.Error("get history record failed", zap.Error(err))
		return nil, err
	}
	return toHistoryDomain(record), nil
}

func (h *historyRepository) SetHistory(ctx context.Context, post domain.Post, actionType string) error {
	err := h.dao.AddHistoryRecord(ctx, post, actionType)
	if err != nil {
		h.l.Error("add history record failed", zap.Error(err))
		return err
	}
	return nil
}

// 转换为领域层
func toHistoryDomain(mh []models.HistoryRecord) []domain.History {
	domainHistory := make([]domain.History, len(mh))
	for i, repoHistory := range mh {
		domainHistory[i] = domain.History{
			PostID:     repoHistory.PostID,
			Title:      repoHistory.Title,
			Content:    repoHistory.Content,
			ActionType: repoHistory.ActionType,
			AuthorID:   repoHistory.AuthorID,
			Status:     repoHistory.Status,
			Slug:       repoHistory.Slug,
			CategoryID: repoHistory.CategoryID,
			Tags:       repoHistory.Tags,
		}
	}
	return domainHistory
}

// 转换为dao层
func fromHistoryDomain(history domain.History) models.HistoryRecord {
	return models.HistoryRecord{
		PostID:     history.PostID,
		Title:      history.Title,
		Content:    history.Content,
		ActionType: history.ActionType,
		AuthorID:   history.AuthorID,
		Status:     history.Status,
		Slug:       history.Slug,
		CategoryID: history.CategoryID,
		Tags:       history.Tags,
	}
}
