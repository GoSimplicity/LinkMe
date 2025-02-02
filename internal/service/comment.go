package service

import (
	"context"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository"
)

// 评论服务结构体
type commentService struct {
	repo repository.CommentRepository
}

// CommentService 评论服务接口
type CommentService interface {
	CreateComment(ctx context.Context, comment domain.Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
	ListComments(ctx context.Context, postId, minID, limit int64) ([]domain.Comment, error)
	GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]domain.Comment, error)
	GetTopCommentsReply(ctx context.Context, postId int64) (domain.Comment, error)
}

// NewCommentService 创建新的评论服务
func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}

// CreateComment 创建评论的实现
func (c *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	// 实现创建评论的逻辑
	return c.repo.CreateComment(ctx, comment)
}

// DeleteComment 删除评论的实现
func (c *commentService) DeleteComment(ctx context.Context, commentId int64) error {
	return c.repo.DeleteComment(ctx, commentId)
}

// GetMoreCommentsReply 获取更多评论回复的实现
func (c *commentService) GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]domain.Comment, error) {
	// 实现获取更多评论回复的逻辑
	return c.repo.GetMoreCommentsReply(ctx, rootId, maxId, limit)
}

// ListComments 列出评论的实现
func (c *commentService) ListComments(ctx context.Context, postId, minID, limit int64) ([]domain.Comment, error) {
	return c.repo.ListComments(ctx, postId, minID, limit)
}

func (c *commentService) GetTopCommentsReply(ctx context.Context, postId int64) (domain.Comment, error) {
	return c.repo.GetTopCommentsReply(ctx, postId)
}
