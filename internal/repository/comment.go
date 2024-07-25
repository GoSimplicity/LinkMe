package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"context"
	"database/sql"
	"time"
)

type commentRepository struct {
	dao dao.CommentDAO
}

type CommentRepository interface {
	CreateComment(ctx context.Context, comment domain.Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
	ListComment(ctx context.Context, Pagination domain.Pagination) ([]domain.Comment, error)
	GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error)
}

func NewCommentService(dao dao.CommentDAO) CommentRepository {
	return &commentRepository{
		dao: dao,
	}
}

// CreateComment implements CommentRepository.
func (c *commentRepository) CreateComment(ctx context.Context, comment domain.Comment) error {
	return c.dao.CreateComment(ctx, c.toDAOComment(comment))
}

// DeleteComment implements CommentRepository.
func (c *commentRepository) DeleteComment(ctx context.Context, commentId int64) error {
	panic("unimplemented")
}

// GetMoreCommentReply implements CommentRepository.
func (c *commentRepository) GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error) {
	panic("unimplemented")
}

// ListComment implements CommentRepository.
func (c *commentRepository) ListComment(ctx context.Context, Pagination domain.Pagination) ([]domain.Comment, error) {
	panic("unimplemented")
}

func (c *commentRepository) toDAOComment(comment domain.Comment) dao.Comment {
	daoComment := dao.Comment{
		Id:        comment.Id,
		UserId:    comment.UserId,
		Biz:       comment.Biz,
		BizId:     comment.BizId,
		PostId:    comment.PostId,
		Content:   comment.Content,
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	if comment.RootComment != nil {
		daoComment.RootID = sql.NullInt64{
			Valid: true,
			Int64: comment.RootComment.Id,
		}
	}
	if comment.ParentComment != nil {
		daoComment.PID = sql.NullInt64{
			Valid: true,
			Int64: comment.ParentComment.Id,
		}
	}
	return daoComment
}
