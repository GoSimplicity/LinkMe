package service

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/check"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"log"
	"time"
)

// 评论服务结构体
type commentService struct {
	repo          repository.CommentRepository
	checkProducer check.Producer
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
func NewCommentService(repo repository.CommentRepository, c check.Producer) CommentService {
	return &commentService{
		repo:          repo,
		checkProducer: c,
	}
}

// CreateComment 创建评论的实现
func (c *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	// 实现创建评论的逻辑

	// 异步发送审核事件
	// 设置超时上下文
	//ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//defer cancel()
	//comment.BizId = 1 // 表示审核业务类型为帖子
	//asyncPublish := general.WithAsyncCancel(ctx, cancel, func() error {
	//	return c.checkProducer.ProduceCheckEvent(check.CheckEvent{
	//		BizId:   comment.BizId, // 表示审核业务类型为帖子
	//		PostId:  uint(comment.PostId),
	//		Content: comment.Content,
	//		Uid:     comment.UserId,
	//	})
	//})
	//asyncPublish()

	// 创建评论
	commentId, err := c.repo.CreateComment(ctx, comment)
	if err != nil {
		return fmt.Errorf("发布评论失败: %w", err)
	}
	// 异步发送审核事件
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	//fmt.Println("uint(comment.Id):%d", uint(comment.Id))
	comment.BizId = 2 // 表示审核业务类型为评论类型

	go func() {
		// 异步执行检查事件的发送
		if err := c.checkProducer.ProduceCheckEvent(check.CheckEvent{
			BizId:   comment.BizId, // 表示审核业务类型为帖子
			PostId:  uint(commentId),
			Content: comment.Content,
			Uid:     comment.UserId,
		}); err != nil {
			// 处理错误，如果需要的话
			log.Printf("Error sending check event: %v", err)
		}
	}()

	return nil
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
