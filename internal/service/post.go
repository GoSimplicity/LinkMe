package service

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository"
	"context"
	"go.uber.org/zap"
)

type PostService interface {
	Create(ctx context.Context, post domain.Post) (int64, error)                                 // 用于创建新帖子
	Update(ctx context.Context, post domain.Post) error                                          // 用于更新现有帖子
	Publish(ctx context.Context, post domain.Post) error                                         // 用于发布帖子
	Withdraw(ctx context.Context, post domain.Post) error                                        // 用于撤回帖子
	GetDraftsByAuthor(ctx context.Context, authorId int64) ([]domain.Post, error)                // 获取作者的草稿
	GetPostById(ctx context.Context, postId int64) (domain.Post, error)                          // 获取特定ID的帖子
	GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error)                 // 获取特定ID的已发布帖子
	ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) // 获取已发布的帖子列表，支持分页
	Delete(ctx context.Context, postId int64) error                                              // 删除帖子
}

type postService struct {
	repo repository.PostRepository
	l    *zap.Logger
}

func NewPostService(repo repository.PostRepository, l *zap.Logger) PostService {
	return &postService{
		repo: repo,
		l:    l,
	}
}

func (p *postService) Create(ctx context.Context, post domain.Post) (int64, error) {
	post.Status = domain.Draft
	// 执行创建操作后默认将帖子状态设置为草稿状态
	return p.repo.Create(ctx, post)
}

func (p *postService) Update(ctx context.Context, post domain.Post) error {
	post.Status = domain.Draft
	// 执行更新操作后默认将帖子状态设置为草稿状态,需手动执行发布操作
	if _, err := p.repo.Sync(ctx, post); err != nil {
		p.l.Error("数据库同步失败", zap.Error(err))
		return err
	}
	return p.repo.Update(ctx, post)
}

func (p *postService) Publish(ctx context.Context, post domain.Post) error {
	post.Status = domain.Published
	// 公开帖子时执行同步操作,添加帖子到线上库
	if _, err := p.repo.Sync(ctx, post); err != nil {
		p.l.Error("数据库同步失败", zap.Error(err))
		return err
	}
	return p.repo.UpdateStatus(ctx, post)
}

func (p *postService) Withdraw(ctx context.Context, post domain.Post) error {
	post.Status = domain.Withdrawn
	// 撤回帖子时执行同步操作,从线上库(mongodb)中移除帖子
	if _, err := p.repo.Sync(ctx, post); err != nil {
		p.l.Error("数据库同步失败", zap.Error(err))
		return err
	}
	return p.repo.UpdateStatus(ctx, post)
}

func (p *postService) GetDraftsByAuthor(ctx context.Context, authorId int64) ([]domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postService) GetPostById(ctx context.Context, postId int64) (domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postService) GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postService) ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return p.repo.ListPublishedPosts(ctx, pagination)
}

func (p *postService) Delete(ctx context.Context, postId int64) error {
	//TODO implement me
	panic("implement me")
}
