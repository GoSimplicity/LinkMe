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
	GetDraftsByAuthor(ctx context.Context, postId int64, uid int64) (domain.Post, error)         // 获取作者的草稿
	GetPostById(ctx context.Context, postId int64, uid int64) (domain.Post, error)               // 获取特定ID的帖子
	GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error)                 // 获取特定ID的已发布帖子
	ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) // 获取已发布的帖子列表，支持分页
	ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)          // 获取个人帖子列表，支持分页
	Delete(ctx context.Context, postId int64, uid int64) error                                   // 删除帖子
}

type postService struct {
	repo repository.PostRepository
	svc  InteractiveService
	l    *zap.Logger
}

func NewPostService(repo repository.PostRepository, l *zap.Logger, svc InteractiveService) PostService {
	return &postService{
		repo: repo,
		l:    l,
		svc:  svc,
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
		p.l.Error("db sync filed", zap.Error(err))
		return err
	}
	return p.repo.Update(ctx, post)
}

func (p *postService) Publish(ctx context.Context, post domain.Post) error {
	post.Status = domain.Published
	// 公开帖子时执行同步操作,添加帖子到线上库
	if _, err := p.repo.Sync(ctx, post); err != nil {
		p.l.Error("db sync filed", zap.Error(err))
		return err
	}
	return p.repo.UpdateStatus(ctx, post)
}

func (p *postService) Withdraw(ctx context.Context, post domain.Post) error {
	post.Status = domain.Withdrawn
	// 撤回帖子时执行同步操作,从线上库(mongodb)中移除帖子
	if _, err := p.repo.Sync(ctx, post); err != nil {
		p.l.Error("db sync filed", zap.Error(err))
		return err
	}
	return p.repo.UpdateStatus(ctx, post)
}

func (p *postService) GetDraftsByAuthor(ctx context.Context, postId int64, uid int64) (domain.Post, error) {
	dp, err := p.repo.GetDraftsByAuthor(ctx, postId, uid)
	if err != nil {
		p.l.Error("get post filed", zap.Error(err))
		return domain.Post{}, err
	}
	return dp, nil
}

func (p *postService) GetPostById(ctx context.Context, postId int64, uid int64) (domain.Post, error) {
	dp, err := p.repo.GetPostById(ctx, postId, uid)
	if err != nil {
		p.l.Error("get post filed", zap.Error(err))
		return domain.Post{}, err
	}
	return dp, err
}

func (p *postService) GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error) {
	dp, err := p.repo.GetPublishedPostById(ctx, postId)
	if err != nil {
		p.l.Error("get pub post filed", zap.Error(err))
		return domain.Post{}, err
	}
	// 增加阅读计数
	if er := p.svc.IncrReadCnt(ctx, "post", postId); er != nil {
		p.l.Error("incr read cnt filed", zap.Error(er))
	}

	return dp, nil
}

func (p *postService) ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return p.repo.ListPosts(ctx, pagination)
}

func (p *postService) ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	// 计算偏移量
	offset := int64(pagination.Page-1) * *pagination.Size
	pagination.Offset = &offset
	return p.repo.ListPublishedPosts(ctx, pagination)
}

func (p *postService) Delete(ctx context.Context, postId int64, uid int64) error {
	pd, err := p.repo.GetPostById(ctx, postId, uid)
	// 避免帖子被重复删除
	if err != nil || pd.Deleted != false {
		p.l.Error("delete post filed", zap.Error(err))
		return err
	}
	post := domain.Post{
		ID:     postId,
		Status: domain.Deleted,
		Author: domain.Author{
			Id: uid,
		},
	}
	if _, er := p.repo.Sync(ctx, post); er != nil {
		p.l.Error("db sync filed", zap.Error(er))
		return er
	}
	return p.repo.Delete(ctx, post)
}
