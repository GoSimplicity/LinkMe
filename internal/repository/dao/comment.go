package dao

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ErrDataNotFound 错误：数据未找到
var ErrDataNotFound = gorm.ErrRecordNotFound

// 评论数据访问对象
type commentDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

// Comment 评论模型定义
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

// CommentDAO 评论数据访问接口
type CommentDAO interface {
	CreateComment(ctx context.Context, comment Comment) error
	DeleteCommentById(ctx context.Context, commentId int64) error
	FindCommentsByPostId(ctx context.Context, postId int64, minId, limit int64) ([]Comment, error)
	GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]Comment, error)
	FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error)
}

// NewCommentService 创建新的评论服务
func NewCommentService(db *gorm.DB, l *zap.Logger) CommentDAO {
	return &commentDAO{
		db: db,
		l:  l,
	}
}

// CreateComment 创建评论
func (c *commentDAO) CreateComment(ctx context.Context, comment Comment) error {
	if err := c.db.WithContext(ctx).Create(&comment).Error; err != nil {
		c.l.Error("create comment failed", zap.Error(err))
		return err
	}
	return nil
}

// DeleteCommentById 根据ID删除评论
func (c *commentDAO) DeleteCommentById(ctx context.Context, commentId int64) error {
	if err := c.db.WithContext(ctx).Delete(&Comment{
		Id: commentId,
	}).Error; err != nil {
		c.l.Error("delete comment failed", zap.Error(err))
		return err
	}
	return nil
}

// GetMoreCommentsReply 获取更多评论回复
func (c *commentDAO) GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]Comment, error) {
	var comments []Comment
	query := c.db.WithContext(ctx).Where("root_id = ?", rootId)
	// 如果maxId > 0，添加id > maxId的条件
	// 当页面初次加载时，不传 maxId 或 maxId 为 0，查询条件是 root_id = ?，这时会加载最早的 limit 条回复
	// 当用户需要加载更多回复时，会将当前已经加载的最大回复 ID 作为 maxId 传递给接口，这样就避免了加载已经获取过的回复
	if maxId > 0 {
		query = query.Where("id > ?", maxId)
	}
	if err := query.Order("id ASC").Limit(int(limit)).Find(&comments).Error; err != nil {
		c.l.Error("list replies failed", zap.Error(err))
		return nil, err
	}
	return comments, nil
}

// FindCommentsByPostId 根据postID查找评论
func (c *commentDAO) FindCommentsByPostId(ctx context.Context, postId int64, minId, limit int64) ([]Comment, error) {
	var comments []Comment
	query := c.db.WithContext(ctx).Where("post_id = ? AND pid IS NULL", postId)
	// 如果minId > 0，添加id < minId的条件
	// 当页面初次加载时，不传 minId 或 minId 为 0，查询条件是 post_id = ? AND pid IS NULL，这时会加载最新的 limit 条评论
	// 当用户需要加载更多评论时，会将当前已经加载的最小评论 ID 作为 minId 传递给接口，这样就避免了加载已经获取过的评论
	if minId > 0 {
		query = query.Where("id < ?", minId)
	}
	if err := query.Limit(int(limit)).Find(&comments).Error; err != nil {
		c.l.Error("list comments failed", zap.Error(err))
		return nil, err
	}
	return comments, nil
}

// FindRepliesByRid 根据根评论ID查找回复
func (c *commentDAO) FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error) {
	var res []Comment
	if err := c.db.WithContext(ctx).
		Where("root_id = ? AND id > ?", rid, id).
		Order("id ASC").
		Limit(int(limit)).Find(&res).Error; err != nil {
		c.l.Error("list comments failed", zap.Error(err))
		return nil, err
	}
	return res, nil
}
