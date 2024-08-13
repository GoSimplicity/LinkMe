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
	historyRepo     repository.HistoryRepository // 历史记录
	producer        post.Producer
	searchRepo      repository.SearchRepository
	publishProducer publish.Producer
	l               *zap.Logger
}

func NewPostService(repo repository.PostRepository, l *zap.Logger, p post.Producer, historyRepo repository.HistoryRepository, searchRepo repository.SearchRepository, publishProducer publish.Producer) PostService {
	return &postService{
		repo:            repo,
		historyRepo:     historyRepo,
		searchRepo:      searchRepo,
		l:               l,
		producer:        p,
		publishProducer: publishProducer,
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

func (p *postService) Publish(ctx context.Context, post domain.Post) error {
	// 获取帖子详细信息
	dp, err := p.repo.GetPostById(ctx, post.ID, post.AuthorID)
	if err != nil {
		return fmt.Errorf("get post failed: %w", err)
	}

	// 使用context.WithCancel来管理goroutine的生命周期
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		// 确保goroutine中的panic不会导致程序崩溃
		defer func() {
			if r := recover(); r != nil {
				p.l.Error("panic occurred in produce publish event goroutine", zap.Any("error", r))
			}
		}()

		// 生产发布事件
		err := p.publishProducer.ProducePublishEvent(publish.PublishEvent{
			PostId:   dp.ID,
			Content:  dp.Content,
			Title:    dp.Title,
			AuthorID: dp.AuthorID,
		})

		if err != nil {
			p.l.Error("produce publish event failed", zap.Error(err))
			cancel() // 取消操作
		}
	}()

	return nil
}

// Withdraw 撤回帖子，移除线上数据库中的帖子
func (p *postService) Withdraw(ctx context.Context, post domain.Post) error {
	post.Status = domain.Withdrawn
	if err := p.searchRepo.DeletePostIndex(ctx, post.ID); err != nil {
		p.l.Error("delete post index failed", zap.Error(err))
	}
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

	// 异步存入历史记录
	go func() {
		err := (func() error {
			if er := p.historyRepo.SetHistory(ctx, dp); er != nil {
				p.l.Error("set history failed", zap.Error(er))
				return fmt.Errorf("set history failed: %w", er)
			}
			return nil
		})()
		if err != nil {
		}
	}()

	// 异步处理读取事件
	go func() {
		err := (func() error {
			if er := p.producer.ProduceReadEvent(post.ReadEvent{PostId: postId, Uid: uid}); er != nil {
				p.l.Error("produce read event failed", zap.Error(err))
				return fmt.Errorf("produce read event failed: %w", er)
			}
			return nil
		})()
		if err != nil {
		}
	}()

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

	go func() {
		err := p.searchRepo.DeletePostIndex(ctx, postId)
		if err != nil {
			p.l.Error("delete post index failed", zap.Error(err))
			return
		}
	}()

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
