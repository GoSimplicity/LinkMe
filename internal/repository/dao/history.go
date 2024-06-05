package dao

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/models"
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// HistoryDAO 接口定义
type HistoryDAO interface {
	AddHistoryRecord(ctx context.Context, post domain.Post, actionType string) error
	GetHistoryRecord(ctx context.Context, pagination domain.Pagination) ([]models.HistoryRecord, error)
}

// historyDAO 结构体
type historyDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

// NewHistoryDAO 创建新的 HistoryDAO 实例
func NewHistoryDAO(db *gorm.DB, l *zap.Logger) HistoryDAO {
	return &historyDAO{
		db: db,
		l:  l,
	}
}

// AddHistoryRecord 添加历史记录
func (h *historyDAO) AddHistoryRecord(ctx context.Context, post domain.Post, actionType string) error {
	historyRecord := models.HistoryRecord{
		PostID:     post.ID,
		Title:      post.Title,
		Content:    post.Content,
		ActionType: actionType,
		ActionTime: time.Now().UnixMilli(),
		AuthorID:   post.Author.Id,
		Status:     post.Status,
		Slug:       post.Slug,
		CategoryID: post.CategoryID,
		Tags:       post.Tags,
	}
	if err := h.db.WithContext(ctx).Create(&historyRecord).Error; err != nil {
		h.l.Error("failed to add history record", zap.Error(err))
		return err
	}
	return nil
}

// GetHistoryRecord 获取历史记录
func (h *historyDAO) GetHistoryRecord(ctx context.Context, pagination domain.Pagination) ([]models.HistoryRecord, error) {
	var historyRecords []models.HistoryRecord
	status := domain.Published
	actionType := "created"
	if err := h.db.WithContext(ctx).Where("author_id = ? AND action_type = ? AND status = ?", pagination.Uid, actionType, status).Find(&historyRecords).Error; err != nil {
		h.l.Error("failed to get history records", zap.Error(err))
		return nil, err
	}
	return historyRecords, nil
}
