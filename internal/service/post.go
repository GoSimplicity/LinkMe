package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/check"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/post"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"github.com/GoSimplicity/LinkMe/pkg/general"
	"go.uber.org/zap"
)

type PostService interface {
	Create(ctx context.Context, post domain.Post) (uint, error)
	Update(ctx context.Context, post domain.Post) error
	Publish(ctx context.Context, postId uint, uid int64) error
	Withdraw(ctx context.Context, postId uint, uid int64) error
	GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error)
	GetPublishPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error)
	ListPublishPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	Delete(ctx context.Context, postId uint, uid int64) error
}

type postService struct {
	repo          repository.PostRepository
	producer      post.Producer
	checkProducer check.Producer
	l             *zap.Logger
}

func NewPostService(repo repository.PostRepository, l *zap.Logger, p post.Producer, c check.Producer) PostService {
	return &postService{
		repo:          repo,
		l:             l,
		producer:      p,
		checkProducer: c,
	}
}

// Create 创建帖子，默认状态为草稿
func (p *postService) Create(ctx context.Context, post domain.Post) (uint, error) {
	return p.repo.Create(ctx, post)
}

// Update 更新帖子，默认状态为草稿
func (p *postService) Update(ctx context.Context, post domain.Post) error {
	return p.repo.Update(ctx, post)
}

// Publish 发布帖子
func (p *postService) Publish(ctx context.Context, postId uint, uid int64) error {
	// 获取帖子详细信息
	dp, err := p.repo.GetPostById(ctx, postId, uid)
	if err != nil {
		return fmt.Errorf("获取帖子失败: %w", err)
	}

	if dp.IsSubmit {
		return errors.New("帖子已提交审核，请勿重复提交")
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 异步发送审核事件
	asyncPublish := general.WithAsyncCancel(ctx, cancel, func() error {
		return p.checkProducer.ProduceCheckEvent(check.CheckEvent{
			PostId:  dp.ID,
			Content: dp.Content,
			Title:   dp.Title,
			Uid:     dp.Uid,
		})
	})

	asyncPublish()

	dp.IsSubmit = true
	// 更新帖子状态
	if err := p.repo.Update(ctx, dp); err != nil {
		p.l.Error("更新帖子失败", zap.Error(err))
		return fmt.Errorf("更新帖子失败: %w", err)
	}

	return nil
}

// Withdraw 撤回帖子，移除线上数据库中的帖子
func (p *postService) Withdraw(ctx context.Context, postId uint, uid int64) error {
	return p.repo.UpdateStatus(ctx, postId, uid, domain.Withdrawn)
}

// GetPostById 获取帖子详细信息
func (p *postService) GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	return p.repo.GetPostById(ctx, postId, uid)
}

// GetPublishPostById 获取已发布的帖子详细信息
func (p *postService) GetPublishPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	dp, err := p.repo.GetPublishPostById(ctx, postId)
	if err != nil {
		return domain.Post{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// 使用装饰器模式异步处理读取事件
	asyncReadEvent := general.WithAsyncCancel(ctx, cancel, func() error {
		if er := p.producer.ProduceReadEvent(post.ReadEvent{
			PostId:  postId,
			Uid:     uid,
			Title:   dp.Title,
			Content: dp.Content,
		}); er != nil {
			p.l.Error("produce read event failed", zap.Error(er))
			return fmt.Errorf("produce read event failed: %w", er)
		}
		return nil
	})
	asyncReadEvent()

	return dp, nil
}

// ListPosts 列出帖子
func (p *postService) ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return p.repo.ListPosts(ctx, pagination)
}

// ListPublishPosts 列出已发布的帖子
func (p *postService) ListPublishPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return p.repo.ListPublishPosts(ctx, pagination)
}

// Delete 删除帖子
func (p *postService) Delete(ctx context.Context, postId uint, uid int64) error {
	_, err := p.repo.GetPostById(ctx, postId, uid)
	if err != nil {
		return err
	}
	return p.repo.Delete(ctx, postId, uid)
}
