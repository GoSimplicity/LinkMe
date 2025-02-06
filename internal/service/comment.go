package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/check"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/pkg/general"
	"golang.org/x/sync/errgroup"
)

type commentService struct {
	repo          repository.CommentRepository
	checkProducer check.Producer
}

type CommentService interface {
	CreateComment(ctx context.Context, comment domain.Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
	ListComments(ctx context.Context, postId, minID, limit int64) ([]domain.Comment, error)
	GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]domain.Comment, error)
	GetTopCommentsReply(ctx context.Context, postId int64) (domain.Comment, error)
}

func NewCommentService(repo repository.CommentRepository, c check.Producer) CommentService {
	return &commentService{
		repo:          repo,
		checkProducer: c,
	}
}

// CreateComment 创建评论的实现
func (c *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	// 参数校验
	if comment.Content == "" {
		return fmt.Errorf("评论内容不能为空")
	}

	// 创建评论
	commentId, err := c.repo.CreateComment(ctx, comment)
	if err != nil || commentId == 0 {
		return fmt.Errorf("发布评论失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	general.WithAsyncCancel(ctx, cancel, func() error {
		// 异步发送审核事件
		go func() {
			event := check.CheckEvent{
				BizId:   2, // 表示审核业务类型为评论
				PostId:  uint(commentId),
				Content: comment.Content,
				Uid:     comment.UserId,
			}

			if err := c.checkProducer.ProduceCheckEvent(event); err != nil {
				log.Printf("发送评论审核事件失败: %v", err)
				return
			}
		}()
		return nil
	})()

	return nil
}

// DeleteComment 删除评论的实现
func (c *commentService) DeleteComment(ctx context.Context, commentId int64) error {
	return c.repo.DeleteComment(ctx, commentId)
}

// GetMoreCommentsReply 获取更多评论回复的实现
func (c *commentService) GetMoreCommentsReply(ctx context.Context, rootId, maxId, limit int64) ([]domain.Comment, error) {
	return c.repo.GetMoreCommentsReply(ctx, rootId, maxId, limit)
}

// ListComments 列出评论的实现
func (c *commentService) ListComments(ctx context.Context, postId, minID, limit int64) ([]domain.Comment, error) {
	// 获取评论列表
	comments, err := c.repo.ListComments(ctx, postId, minID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取评论列表失败: %w", err)
	}

	// 初始化返回的评论列表
	domainComments := make([]domain.Comment, 0, len(comments))
	var eg errgroup.Group

	// 并发处理每个评论
	for i, comment := range comments {
		domainComments = append(domainComments, comment)

		i := i // 创建副本避免闭包问题
		commentId := comment.Id
		eg.Go(func() error {
			// 获取最新的三条子评论
			subComments, err := c.repo.GetMoreCommentsReply(ctx, commentId, 0, 3)
			if err != nil {
				return fmt.Errorf("获取子评论失败: %w", err)
			}

			domainComments[i].Children = subComments
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return domainComments, nil
}

func (c *commentService) GetTopCommentsReply(ctx context.Context, postId int64) (domain.Comment, error) {
	return c.repo.GetTopCommentsReply(ctx, postId)
}
