package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/GoSimplicity/LinkMe/pkg/change"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostRepository interface {
	Create(ctx context.Context, post domain.Post) (uint, error)
	Update(ctx context.Context, post domain.Post) error
	UpdateStatus(ctx context.Context, postId uint, uid int64, status uint8) error
	GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error)
	GetPublishPostById(ctx context.Context, postId uint) (domain.Post, error)
	ListPublishPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	Delete(ctx context.Context, postId uint, uid int64) error
	GetPost(ctx context.Context, postId uint) (domain.Post, error)
	ListAllPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
}

type postRepository struct {
	dao dao.PostDAO
	l   *zap.Logger
}

func NewPostRepository(dao dao.PostDAO, l *zap.Logger) PostRepository {
	return &postRepository{
		dao: dao,
		l:   l,
	}
}

// Create 创建帖子
func (p *postRepository) Create(ctx context.Context, post domain.Post) (uint, error) {
	// 设置帖子的唯一标识符
	post.Slug = uuid.New().String() + strconv.Itoa(int(post.ID))

	id, err := p.dao.Insert(ctx, change.FromDomainPost(post))
	if err != nil {
		p.l.Error("创建帖子失败", zap.Error(err))
		return 0, fmt.Errorf("创建帖子失败: %w", err)
	}

	return id, nil
}

// Update 更新帖子
func (p *postRepository) Update(ctx context.Context, post domain.Post) error {
	if err := p.dao.Update(ctx, change.FromDomainPost(post)); err != nil {
		p.l.Error("更新帖子失败", zap.Error(err), zap.Uint("post_id", post.ID))
		return fmt.Errorf("更新帖子失败: %w", err)
	}
	return nil
}

// UpdateStatus 更新帖子状态
func (p *postRepository) UpdateStatus(ctx context.Context, postId uint, uid int64, status uint8) error {
	if err := p.dao.UpdateStatus(ctx, postId, uid, status); err != nil {
		p.l.Error("更新帖子状态失败", zap.Error(err), zap.Uint("post_id", postId))
		return fmt.Errorf("更新帖子状态失败: %w", err)
	}
	return nil
}

// GetPostById 获取帖子详细信息
func (p *postRepository) GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	dp, err := p.dao.GetById(ctx, postId, uid)
	if err != nil {
		p.l.Error("获取帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		return domain.Post{}, fmt.Errorf("获取帖子失败: %w", err)
	}

	return change.ToDomainPost(dp), nil
}

// GetPublishedPostById 获取已发布的帖子详细信息
func (p *postRepository) GetPublishPostById(ctx context.Context, postId uint) (domain.Post, error) {
	dp, err := p.dao.GetPubById(ctx, postId)
	if err != nil {
		p.l.Error("获取已发布帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		return domain.Post{}, err
	}

	return change.ToDomainListPubPost(dp), nil
}

// ListPosts 获取作者帖子的列表
func (p *postRepository) ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	pub, err := p.dao.List(ctx, pagination)
	if err != nil {
		p.l.Error("获取帖子列表失败", zap.Error(err))
		return nil, fmt.Errorf("获取帖子列表失败: %w", err)
	}

	return change.FromDomainSlicePost(pub), nil
}

// ListPublishedPosts 获取已发布的帖子列表
func (p *postRepository) ListPublishPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	pub, err := p.dao.ListPub(ctx, pagination)
	if err != nil {
		p.l.Error("从数据库获取已发布帖子列表失败", zap.Error(err))
		return nil, fmt.Errorf("从数据库获取已发布帖子列表失败: %w", err)
	}

	if len(pub) == 0 {
		p.l.Info("没有找到已发布的帖子")
		return []domain.Post{}, nil
	}

	return change.FromDomainSlicePubPostList(pub), nil
}

// Delete 删除帖子
func (p *postRepository) Delete(ctx context.Context, postId uint, uid int64) error {
	if err := p.dao.Delete(ctx, postId, uid); err != nil {
		p.l.Error("删除帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		return fmt.Errorf("删除帖子失败: %w", err)
	}
	return nil
}

// GetPost 获取帖子
func (p *postRepository) GetPost(ctx context.Context, postId uint) (domain.Post, error) {
	dp, err := p.dao.GetPost(ctx, postId)
	if err != nil {
		p.l.Error("获取帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		return domain.Post{}, fmt.Errorf("获取帖子失败: %w", err)
	}

	return change.ToDomainPost(dp), nil
}

// ListAllPosts 获取所有帖子
func (p *postRepository) ListAllPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	posts, err := p.dao.ListAll(ctx, pagination)
	if err != nil {
		p.l.Error("获取所有帖子列表失败", zap.Error(err))
		return nil, fmt.Errorf("获取所有帖子列表失败: %w", err)
	}

	if len(posts) == 0 {
		p.l.Info("没有找到任何帖子")
		return []domain.Post{}, nil
	}

	return change.FromDomainSlicePost(posts), nil
}
