package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ErrDataNotFound 定义一个全局的记录未找到错误
var ErrDataNotFound = gorm.ErrRecordNotFound

// 评论数据访问对象
type commentDAO struct {
	db *gorm.DB    // 数据库连接实例
	l  *zap.Logger // 日志实例
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
	PID           sql.NullInt64 `gorm:"column:pid;index;default:1"`                                          // 父评论ID
	ParentComment *Comment      `gorm:"ForeignKey:PID;AssociationForeignKey:ID;constraint:OnDelete:CASCADE"` // 父评论
	CreatedAt     int64         `gorm:"autoCreateTime"`                                                      // 创建时间
	UpdatedAt     int64         `gorm:"autoUpdateTime"`                                                      // 更新时间
	Status        uint8         `gorm:"default:0"`                                                           // 评论状态 和domain/post.go中的Status对应
}

// CommentDAO 评论数据访问接口定义
type CommentDAO interface {
	CreateComment(ctx context.Context, comment Comment) error
	DeleteCommentById(ctx context.Context, commentId int64) error
	FindCommentsByPostId(ctx context.Context, postId int64, minId, limit int64) ([]Comment, error)
	GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]Comment, error)
	FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error)
	FindTopCommentsByPostId(ctx context.Context, postId int64) (Comment, error)
	FindCommentByCommentId(ctx context.Context, commentId int64) (Comment, error)
	UpdateComment(ctx context.Context, comment Comment) error
}

// NewCommentDAO 创建新的评论服务
func NewCommentDAO(db *gorm.DB, l *zap.Logger) CommentDAO {
	return &commentDAO{
		db: db,
		l:  l,
	}
}

// CreateComment 创建评论
func (c *commentDAO) CreateComment(ctx context.Context, comment Comment) error {
	if err := c.db.WithContext(ctx).Create(&comment).Error; err != nil {
		c.l.Error("创建评论失败", zap.Error(err))
		return err
	}
	return nil
}

// FindCommentByCommentId 根据commentId查找评论
func (c *commentDAO) FindCommentByCommentId(ctx context.Context, commentId int64) (Comment, error) {
	// 先查找评论是否存在
	var comment Comment
	if err := c.db.WithContext(ctx).First(&comment, "id = ?", commentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.l.Error("评论不存在", zap.Int64("commentId", commentId), zap.Error(err))
			return comment, fmt.Errorf("评论不存在")
		}
		c.l.Error("查找评论失败", zap.Error(err))
		return Comment{}, err
	}
	return Comment{}, nil
}
func (c *commentDAO) UpdateComment(ctx context.Context, comment Comment) error {
	// 确保更新的是指定 ID 的评论
	if err := c.db.WithContext(ctx).Model(&Comment{}).Where("id = ?", comment.Id).Updates(comment).Error; err != nil {
		c.l.Error("更新评论失败", zap.Error(err))
		return err
	}
	return nil
}

// DeleteCommentById 根据ID删除评论
func (c *commentDAO) DeleteCommentById(ctx context.Context, commentId int64) error {
	if err := c.db.WithContext(ctx).Delete(&Comment{Id: commentId}).Error; err != nil {
		c.l.Error("删除评论失败", zap.Error(err))
		return err
	}
	return nil
}

// GetMoreCommentsReply 获取更多评论回复
func (c *commentDAO) GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]Comment, error) {
	var comments []Comment

	query := c.db.WithContext(ctx).Where("root_id = ?", rootId)
	// 如果 maxId > 0，则只获取比 maxId 大的记录，避免重复加载
	if maxId > 0 {
		query = query.Where("id > ?", maxId)
	}

	// 按 ID 升序排列，获取 limit 条记录
	if err := query.Order("id ASC").Limit(int(limit)).Find(&comments).Error; err != nil {
		c.l.Error("获取评论回复失败", zap.Error(err))
		return nil, err
	}

	return comments, nil
}

// FindCommentsByPostId 根据 postID 查找评论
func (c *commentDAO) FindCommentsByPostId(ctx context.Context, postId int64, minId, limit int64) ([]Comment, error) {
	var comments []Comment

	query := c.db.WithContext(ctx).Where("post_id = ? AND pid IS NULL", postId)
	// 如果 minId > 0，则只获取比 minId 小的记录，避免重复加载
	if minId > 0 {
		query = query.Where("id < ?", minId)
	}

	// 获取 limit 条记录
	if err := query.Order("id DESC").Limit(int(limit)).Find(&comments).Error; err != nil {
		c.l.Error("获取评论失败", zap.Error(err))
		return nil, err
	}

	return comments, nil
}

func (c *commentDAO) FindTopCommentsByPostId(ctx context.Context, postId int64) (Comment, error) {
	var comment Comment
	pidValue := 1 // Note:固定值为 1
	query := c.db.WithContext(ctx).Where("post_id = ? AND pid = ?", postId, pidValue)
	// 获取 limit 条记录
	limit := 1 // Note:这里强制获取1条
	if err := query.Order("id DESC").Limit(int(limit)).Find(&comment).Error; err != nil {
		c.l.Error("获取评论失败", zap.Error(err))
		return Comment{}, err
	}

	return comment, nil
}

// FindRepliesByRid 根据根评论ID查找回复
func (c *commentDAO) FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error) {
	var replies []Comment

	// 按照 root_id 和 id > ? 过滤并按 id 升序排列，获取 limit 条记录
	if err := c.db.WithContext(ctx).
		Where("root_id = ? AND id > ?", rid, id).
		Order("id ASC").
		Limit(int(limit)).Find(&replies).Error; err != nil {
		c.l.Error("获取评论回复失败", zap.Error(err))
		return nil, err
	}

	return replies, nil
}
