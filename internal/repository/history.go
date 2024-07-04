package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/cache"
	"context"
	"go.uber.org/zap"
)

type HistoryRepository interface {
	GetHistory(ctx context.Context, pagination domain.Pagination) ([]domain.History, error)
	SetHistory(ctx context.Context, post domain.Post) error
	DeleteOneHistory(ctx context.Context, postId int64, uid int64) error
	DeleteAllHistory(ctx context.Context, uid int64) error
}

type historyRepository struct {
	l     *zap.Logger
	cache cache.HistoryCache
}

func NewHistoryRepository(l *zap.Logger, cache cache.HistoryCache) HistoryRepository {
	return &historyRepository{
		l:     l,
		cache: cache,
	}
}

func (h *historyRepository) GetHistory(ctx context.Context, pagination domain.Pagination) ([]domain.History, error) {
	record, err := h.cache.GetCache(ctx, pagination)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (h *historyRepository) SetHistory(ctx context.Context, post domain.Post) error {
	history := toDomainHistory(post)
	err := h.cache.SetCache(ctx, history)
	if err != nil {
		return err
	}
	return nil
}

func (h *historyRepository) DeleteOneHistory(ctx context.Context, postId int64, uid int64) error {
	err := h.cache.DeleteOneCache(ctx, postId, uid)
	if err != nil {
		return err
	}
	return nil
}

func (h *historyRepository) DeleteAllHistory(ctx context.Context, uid int64) error {
	err := h.cache.DeleteAllHistory(ctx, uid)
	if err != nil {
		return err
	}
	return nil
}

// createContentSummary 创建内容摘要，限制为28个汉字
func createContentSummary(content string) string {
	const limit = 28
	runes := []rune(content)
	if len(runes) > limit {
		return string(runes[:limit])
	}
	return content
}

func toDomainHistory(post domain.Post) domain.History {
	return domain.History{
		Content:  createContentSummary(post.Content),
		Deleted:  post.Deleted,
		AuthorID: post.Author.Id,
		Tags:     post.Tags,
		PostID:   post.ID,
		Title:    post.Title,
	}
}
