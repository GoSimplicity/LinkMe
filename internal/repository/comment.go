package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/redis/go-redis/v9"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
)

type commentRepository struct {
	dao   dao.CommentDAO
	cache cache.CommentCache
}

type CommentRepository interface {
	CreateComment(ctx context.Context, comment domain.Comment) (int64, error)
	DeleteComment(ctx context.Context, commentId int64) error
	ListComments(ctx context.Context, postId, minID, limit int64) ([]domain.Comment, error)
	GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]domain.Comment, error)
	GetTopCommentsReply(ctx context.Context, postId int64) (domain.Comment, error)
	FindCommentByCommentId(ctx context.Context, commentId int64) (domain.Comment, error)
	UpdateComment(ctx context.Context, comment domain.Comment) error
}

func NewCommentRepository(dao dao.CommentDAO, cache cache.CommentCache) CommentRepository {
	return &commentRepository{
		dao:   dao,
		cache: cache,
	}
}

// CreateComment 创建评论
func (c *commentRepository) CreateComment(ctx context.Context, comment domain.Comment) (int64, error) {
	return c.dao.CreateComment(ctx, c.toDAOComment(comment))
}

// FindCommentByCommentId 根据评论ID查找评论
func (c *commentRepository) FindCommentByCommentId(ctx context.Context, commentId int64) (domain.Comment, error) {
	comment, err := c.dao.FindCommentByCommentId(ctx, commentId)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("查找评论失败: %w", err)
	}
	return c.toDomainComment(comment), nil
}

// UpdateComment 更新评论
func (c *commentRepository) UpdateComment(ctx context.Context, comment domain.Comment) error {
	return c.dao.UpdateComment(ctx, c.toDAOComment(comment))
}

// DeleteComment 删除评论
func (c *commentRepository) DeleteComment(ctx context.Context, commentId int64) error {
	return c.dao.DeleteCommentById(ctx, commentId)
}

// GetMoreCommentsReply 获取更多评论回复
func (c *commentRepository) GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]domain.Comment, error) {
	comments, err := c.dao.GetMoreCommentsReply(ctx, rootId, maxId, limit)
	if err != nil {
		return nil, fmt.Errorf("获取评论回复失败: %w", err)
	}
	return c.toDomainSliceComments(comments), nil
}

// GetTopCommentsReply 加载顶级评论
func (c *commentRepository) GetTopCommentsReply(ctx context.Context, postId int64) (domain.Comment, error) {
	// 优先从缓存获取
	if comments, err := c.cache.Get(ctx, postId); err == nil {
		return comments, nil
	} else if err != redis.Nil {
		return domain.Comment{}, fmt.Errorf("获取缓存失败: %w", err)
	}

	// 缓存未命中,从数据库加载
	daoComments, err := c.dao.FindTopCommentsByPostId(ctx, postId)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("从数据库加载评论失败: %w", err)
	}

	// 转换并写入缓存
	domainComment := c.toDomainComment(daoComments)
	if err := c.cache.Set(ctx, domainComment); err != nil {
		return domain.Comment{}, fmt.Errorf("写入缓存失败: %w", err)
	}

	return domainComment, nil
}

// ListComments 列出评论
func (c *commentRepository) ListComments(ctx context.Context, postId int64, minId, limit int64) ([]domain.Comment, error) {
	// 获取评论列表
	daoComments, err := c.dao.FindCommentsByPostId(ctx, postId, minId, limit)
	if err != nil {
		return nil, fmt.Errorf("获取评论列表失败: %w", err)
	}

	return c.toDomainSliceComments(daoComments), nil
}

// toDAOComment 将领域模型评论转换为DAO评论
func (c *commentRepository) toDAOComment(comment domain.Comment) dao.Comment {
	now := time.Now().UnixMilli()
	daoComment := dao.Comment{
		Id:        comment.Id,
		UserId:    comment.UserId,
		Biz:       comment.Biz,
		BizId:     comment.BizId,
		PostId:    comment.PostId,
		Content:   comment.Content,
		CreatedAt: now,
		UpdatedAt: now,
		Status:    comment.Status,
	}

	if comment.ParentComment != nil {
		daoComment.PID = sql.NullInt64{
			Int64: comment.ParentComment.Id,
			Valid: true,
		}
	}

	if comment.RootComment != nil {
		daoComment.RootId = sql.NullInt64{
			Int64: comment.RootComment.Id,
			Valid: true,
		}
	}

	return daoComment
}

// toDomainComment 将DAO评论转换为领域模型评论
func (c *commentRepository) toDomainComment(daoComment dao.Comment) domain.Comment {
	domainComment := domain.Comment{
		Id:        daoComment.Id,
		UserId:    daoComment.UserId,
		Biz:       daoComment.Biz,
		BizId:     daoComment.BizId,
		PostId:    daoComment.PostId,
		Content:   daoComment.Content,
		CreatedAt: daoComment.CreatedAt,
		UpdatedAt: daoComment.UpdatedAt,
		Status:    daoComment.Status,
	}

	if daoComment.PID.Valid {
		domainComment.ParentComment = &domain.Comment{Id: daoComment.PID.Int64}
	}

	if daoComment.RootId.Valid {
		domainComment.RootComment = &domain.Comment{Id: daoComment.RootId.Int64}
	}

	return domainComment
}

// toDomainSliceComments 将DAO评论切片转换为领域模型评论切片
func (c *commentRepository) toDomainSliceComments(daoComments []dao.Comment) []domain.Comment {
	domainComments := make([]domain.Comment, 0, len(daoComments))
	for _, daoComment := range daoComments {
		domainComments = append(domainComments, c.toDomainComment(daoComment))
	}
	return domainComments
}
