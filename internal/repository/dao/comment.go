package dao

import (
	"LinkMe/internal/domain"
	"context"
	"database/sql"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrDataNotFound = gorm.ErrRecordNotFound

type commentDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

type Comment struct {
	Id            int64         `gorm:"autoIncrement;primaryKey"`                      // 评论Id
	UserId        int64         `gorm:"index:idx_user_id"`                             // 发表评论的用户ID
	Biz           string        `gorm:"index:idx_biz_type_id"`                         // 业务类型
	BizId         int64         `gorm:"index:idx_biz_type_id"`                         // 业务ID
	Content       string        `gorm:"column:content;type:text"`                      // 评论内容
	PostId        int64         `gorm:"index:idx_post_id"`                             // 帖子ID，用于多级评论
	RootID        sql.NullInt64 `gorm:"column:root_id;index"`                          // 根评论ID
	PID           sql.NullInt64 `gorm:"column:pid;index"`                              // 父评论ID
	RootComment   *Comment      `gorm:"foreignKey:RootID;constraint:OnDelete:CASCADE"` // 根评论
	ParentComment *Comment      `gorm:"foreignKey:PID;constraint:OnDelete:CASCADE"`    // 父评论
	Children      []Comment     `gorm:"foreignKey:PID;constraint:OnDelete:CASCADE"`    // 子评论
	CreatedAt     int64         `gorm:"autoCreateTime"`                                // 创建时间
	UpdatedAt     int64         `gorm:"autoUpdateTime"`                                // 更新时间
}

type CommentDAO interface {
	CreateComment(ctx context.Context, comment Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
	ListComment(ctx context.Context, Pagination domain.Pagination) ([]domain.Comment, error)
	GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error)
}

func NewCommentService(db *gorm.DB, l *zap.Logger) CommentDAO {
	return &commentDAO{
		db: db,
		l:  l,
	}
}

// CreateComment implements CommentDAO.
func (c *commentDAO) CreateComment(ctx context.Context, comment Comment) error {
	if err := c.db.WithContext(ctx).Create(&comment).Error; err != nil {
		c.l.Error("create comment failed", zap.Error(err))
		return err
	}
	return nil
}

// DeleteComment implements CommentDAO.
func (c *commentDAO) DeleteComment(ctx context.Context, commentId int64) error {
	panic("unimplemented")
}

// GetMoreCommentReply implements CommentDAO.
func (c *commentDAO) GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error) {
	panic("unimplemented")
}

// ListComment implements CommentDAO.
func (c *commentDAO) ListComment(ctx context.Context, Pagination domain.Pagination) ([]domain.Comment, error) {
	panic("unimplemented")
}
