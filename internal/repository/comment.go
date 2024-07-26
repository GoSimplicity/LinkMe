package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"context"
	"database/sql"
	"time"

	"golang.org/x/sync/errgroup"
)

// 评论仓库结构体
type commentRepository struct {
	dao dao.CommentDAO
}

// 评论仓库接口
type CommentRepository interface {
	CreateComment(ctx context.Context, comment domain.Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
	ListComments(ctx context.Context, biz string, bizID, minID, limit int64) ([]domain.Comment, error)
	GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error)
}

// 创建新的评论服务
func NewCommentService(dao dao.CommentDAO) CommentRepository {
	return &commentRepository{
		dao: dao,
	}
}

// 创建评论
func (c *commentRepository) CreateComment(ctx context.Context, comment domain.Comment) error {
	return c.dao.CreateComment(ctx, c.toDAOComment(comment))
}

// 删除评论
func (c *commentRepository) DeleteComment(ctx context.Context, commentId int64) error {
	return c.dao.DeleteCommentById(ctx, commentId)
}

// 获取更多评论回复（未实现）
func (c *commentRepository) GetMoreCommentReply(ctx context.Context, commentId int64, pagination domain.Pagination, Id int64) ([]domain.Comment, error) {
	panic("unimplemented")
}

// 列出评论
func (c *commentRepository) ListComments(ctx context.Context, biz string, bizID, minID, limit int64) ([]domain.Comment, error) {
	// 从DAO层获取评论列表
	daoComments, err := c.dao.FindCommentsByBiz(ctx, biz, bizID, minID, limit)
	if err != nil {
		return nil, err
	}

	// 初始化返回的领域模型评论列表
	domainComments := make([]domain.Comment, 0, len(daoComments))
	var errGroup errgroup.Group

	// 遍历每个评论
	for _, daoComment := range daoComments {
		// 将当前评论转换为领域模型评论
		domainComment := c.toDomainComment(daoComment)
		domainComments = append(domainComments, domainComment)

		// 并发获取子评论
		errGroup.Go(func(dc dao.Comment, dcm *domain.Comment) func() error {
			return func() error {
				subComments, err := c.dao.FindRepliesByRid(ctx, dc.Id, 0, 3)
				if err != nil {
					return err
				}
				// 将子评论转换为领域模型评论并添加到当前评论的子评论列表中
				childrenComments := make([]domain.Comment, 0, len(subComments))
				for _, subComment := range subComments {
					childrenComments = append(childrenComments, c.toDomainComment(subComment))
				}
				dcm.Children = childrenComments
				return nil
			}
		}(daoComment, &domainComments[len(domainComments)-1]))
	}

	// 等待所有并发任务完成并返回结果
	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	return domainComments, nil
}

// 将领域模型评论转换为DAO评论
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

// 将DAO评论转换为领域模型评论
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
	}
	if daoComment.PID.Valid {
		domainComment.ParentComment = &domain.Comment{
			Id: daoComment.PID.Int64,
		}
	}
	if daoComment.RootId.Valid {
		domainComment.RootComment = &domain.Comment{
			Id: daoComment.RootId.Int64,
		}
	}
	return domainComment
}
