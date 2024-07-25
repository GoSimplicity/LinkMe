package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
)

type commentService struct {
	repo repository.CommentRepository
}

type CommentService interface {
	CreateComment(ctx context.Context, comment domain.Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
	ListComment(ctx context.Context, Pagination domain.Pagination) ([]domain.Comment, error)
	GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error)
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}

// CreateComment implements CommentService.
func (c *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	panic("unimplemented")
}

// DeleteComment implements CommentService.
func (c *commentService) DeleteComment(ctx context.Context, commentId int64) error {
	panic("unimplemented")
}

// GetMoreCommentReply implements CommentService.
func (c *commentService) GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error) {
	panic("unimplemented")
}

// ListComment implements CommentService.
func (c *commentService) ListComment(ctx context.Context, Pagination domain.Pagination) ([]domain.Comment, error) {
	panic("unimplemented")
}
