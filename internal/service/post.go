package service

import (
	"context"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/post"
	"github.com/GoSimplicity/LinkMe/internal/domain/events/publish"
	"github.com/GoSimplicity/LinkMe/internal/repository"
	"go.uber.org/zap"
)

type PostService interface {
	Create(ctx context.Context, post domain.Post) (uint, error)
	Update(ctx context.Context, post domain.Post) error
	Publish(ctx context.Context, post domain.Post) error
	Withdraw(ctx context.Context, post domain.Post) error
	GetDraftsByAuthor(ctx context.Context, postId uint, uid int64) (domain.Post, error)
	GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error)
	GetPublishedPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error)
	ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	Delete(ctx context.Context, postId uint, uid int64) error
	ListAllPost(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	GetPost(ctx context.Context, postId uint) (domain.Post, error)
	GetPostCount(ctx context.Context) (int64, error)
}

type postService struct {
	repo            repository.PostRepository
	producer        post.Producer
	publishProducer publish.Producer
	l               *zap.Logger
}

func NewPostService(repo repository.PostRepository, l *zap.Logger, p post.Producer, publishProducer publish.Producer) PostService {
	return &postService{
		repo:            repo,
		l:               l,
		producer:        p,
		publishProducer: publishProducer,
	}
}

// withAsyncCancel 装饰器函数，用来封装 goroutine 逻辑并处理错误和取消操作
func (p *postService) withAsyncCancel(_ context.Context, cancel context.CancelFunc, fn func() error) func() {
	return func() {
		go func() {
			// 确保 goroutine 中的 panic 不会导致程序崩溃
			defer func() {
				if r := recover(); r != nil {
					p.l.Error("panic occurred in async operation goroutine", zap.Any("error", r))
					cancel() // 取消操作
				}
			}()

			// 执行目标函数
			if err := fn(); err != nil {
				p.l.Error("async operation failed", zap.Error(err))
				cancel() // 取消操作
			}
		}()
	}
}

// Create 创建帖子，默认状态为草稿
func (p *postService) Create(ctx context.Context, post domain.Post) (uint, error) {
	post.Status = domain.Draft
	// 执行创建操作后默认将帖子状态设置为草稿状态
	return p.repo.Create(ctx, post)
}

// Update 更新帖子，默认状态为草稿
func (p *postService) Update(ctx context.Context, post domain.Post) error {
	post.Status = domain.Draft
	return p.repo.Update(ctx, post)
}

// Publish 发布帖子
func (p *postService) Publish(ctx context.Context, post domain.Post) error {
	// 获取帖子详细信息
	dp, err := p.repo.GetPostById(ctx, post.ID, post.AuthorID)
	if err != nil {
		return fmt.Errorf("get post failed: %w", err)
	}

	// 使用 context.WithCancel 来管理 goroutine 的生命周期
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 使用装饰器封装 goroutine 逻辑
	asyncPublish := p.withAsyncCancel(ctx, cancel, func() error {
		return p.publishProducer.ProducePublishEvent(publish.PublishEvent{
			PostId:   dp.ID,
			Content:  dp.Content,
			Title:    dp.Title,
			AuthorID: dp.AuthorID,
		})
	})
	asyncPublish()

	return nil
}

// Withdraw 撤回帖子，移除线上数据库中的帖子
func (p *postService) Withdraw(ctx context.Context, post domain.Post) error {
	post.Status = domain.Withdrawn

	return p.repo.UpdateStatus(ctx, post)
}

// GetDraftsByAuthor 获取作者的草稿帖子
func (p *postService) GetDraftsByAuthor(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	return p.repo.GetPostById(ctx, postId, uid)
}

// GetPostById 获取帖子详细信息
func (p *postService) GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	return p.repo.GetPostById(ctx, postId, uid)
}

// GetPublishedPostById 获取已发布的帖子详细信息
func (p *postService) GetPublishedPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	dp, err := p.repo.GetPublishedPostById(ctx, postId)
	if err != nil {
		return domain.Post{}, err
	}

	// 使用 context.WithCancel 来管理 goroutine 的生命周期
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 使用装饰器模式异步处理读取事件
	asyncReadEvent := p.withAsyncCancel(ctx, cancel, func() error {
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

// ListPublishedPosts 列出已发布的帖子
func (p *postService) ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return p.repo.ListPublishedPosts(ctx, pagination)
}

// Delete 删除帖子
func (p *postService) Delete(ctx context.Context, postId uint, uid int64) error {
	_, err := p.repo.GetPostById(ctx, postId, uid)
	if err != nil {
		return err
	}
	res := domain.Post{
		ID:       postId,
		Status:   domain.Deleted,
		AuthorID: uid,
	}

	return p.repo.Delete(ctx, res)
}

// ListAllPost 列出所有帖子
func (p *postService) ListAllPost(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return p.repo.ListAllPost(ctx, pagination)
}

// GetPost 获取帖子
func (p *postService) GetPost(ctx context.Context, postId uint) (domain.Post, error) {
	return p.repo.GetPost(ctx, postId)
}

// GetPostCount 获取帖子数量
func (p *postService) GetPostCount(ctx context.Context) (int64, error) {
	return p.repo.GetPostCount(ctx)
}
