package repository

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"go.uber.org/zap"
)

type HistoryRepository interface {
	GetHistory(ctx context.Context, pagination domain.Pagination) ([]domain.History, error)
	SetHistory(ctx context.Context, post domain.Post) error
	DeleteOneHistory(ctx context.Context, postId uint, uid int64) error
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

// GetHistory 获取历史记录
func (h *historyRepository) GetHistory(ctx context.Context, pagination domain.Pagination) ([]domain.History, error) {
	histories, err := h.cache.GetCache(ctx, pagination)
	if err != nil {
		return nil, err
	}

	return histories, nil
}

// SetHistory 设置历史记录
func (h *historyRepository) SetHistory(ctx context.Context, post domain.Post) error {
	history := toDomainHistory(post)

	return h.cache.SetCache(ctx, history)
}

// DeleteOneHistory 删除一条历史记录
func (h *historyRepository) DeleteOneHistory(ctx context.Context, postId uint, uid int64) error {
	return h.cache.DeleteOneCache(ctx, postId, uid)
}

func (h *historyRepository) DeleteAllHistory(ctx context.Context, uid int64) error {
	return h.cache.DeleteAllHistory(ctx, uid)
}

// createContentSummary 创建内容摘要,限制为28个汉字
func createContentSummary(content string) string {
	const limit = 28

	runes := []rune(content)
	if len(runes) <= limit {
		return content
	}

	return string(runes[:limit])
}

// toDomainHistory 将帖子转换为历史记录
func toDomainHistory(post domain.Post) domain.History {
	return domain.History{
		Content: createContentSummary(post.Content),
		Uid:     post.Uid,
		Tags:    post.Tags,
		PostID:  post.ID,
		Title:   post.Title,
	}
}
