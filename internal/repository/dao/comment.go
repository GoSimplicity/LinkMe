package dao

import (
	"context"
	"database/sql"
	"github.com/GoSimplicity/LinkMe/internal/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 错误：数据未找到
var ErrDataNotFound = gorm.ErrRecordNotFound

// 评论数据访问对象
type commentDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

// 评论模型定义
type Comment struct {
	Id            int64         `gorm:"autoIncrement;primaryKey"`                                            // 评论ID
	UserId        int64         `gorm:"index:idx_user_id"`                                                   // 发表评论的用户ID
	Biz           string        `gorm:"index:idx_biz_type_id"`                                               // 业务类型
	BizId         int64         `gorm:"index:idx_biz_type_id"`                                               // 业务ID
	Content       string        `gorm:"column:content;type:text"`                                            // 评论内容
	PostId        int64         `gorm:"index:idx_post_id"`                                                   // 帖子ID，用于多级评论
	RootId        sql.NullInt64 `gorm:"column:root_id;index"`                                                // 根评论ID
	PID           sql.NullInt64 `gorm:"column:pid;index"`                                                    // 父评论ID
	ParentComment *Comment      `gorm:"ForeignKey:PID;AssociationForeignKey:ID;constraint:OnDelete:CASCADE"` // 父评论
	CreatedAt     int64         `gorm:"autoCreateTime"`                                                      // 创建时间
	UpdatedAt     int64         `gorm:"autoUpdateTime"`                                                      // 更新时间
}

// 评论数据访问接口
type CommentDAO interface {
	CreateComment(ctx context.Context, comment Comment) error
	DeleteCommentById(ctx context.Context, commentId int64) error
	FindCommentsByBiz(ctx context.Context, biz string, bizId, minID, limit int64) ([]Comment, error)
	GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error)
	FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error)
}

// 创建新的评论服务
func NewCommentService(db *gorm.DB, l *zap.Logger) CommentDAO {
	return &commentDAO{
		db: db,
		l:  l,
	}
}

// 创建评论
func (c *commentDAO) CreateComment(ctx context.Context, comment Comment) error {
	if err := c.db.WithContext(ctx).Create(&comment).Error; err != nil {
		c.l.Error("create comment failed", zap.Error(err))
		return err
	}
	return nil
}

// 根据ID删除评论
func (c *commentDAO) DeleteCommentById(ctx context.Context, commentId int64) error {
	if err := c.db.WithContext(ctx).Delete(&Comment{
		Id: commentId,
	}).Error; err != nil {
		c.l.Error("delete comment failed", zap.Error(err))
		return err
	}
	return nil
}

// 获取更多评论回复
func (c *commentDAO) GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error) {
	panic("unimplemented")
}

// 根据业务类型和业务ID查找评论
func (c *commentDAO) FindCommentsByBiz(ctx context.Context, biz string, bizId, minID, limit int64) ([]Comment, error) {
	var comments []Comment
	if err := c.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND id < ? AND pid IS NULL", biz, bizId, minID).
		Limit(int(limit)).Find(&comments).Error; err != nil {
		c.l.Error("list comments failed", zap.Error(err))
		return nil, err
	}
	return comments, nil
}

// 根据根评论ID查找回复
func (c *commentDAO) FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).
		Where("root_id = ? AND id > ?", rid, id).
		Order("id ASC").
		Limit(int(limit)).Find(&res).Error
	return res, err
}
