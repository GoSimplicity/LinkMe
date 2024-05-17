package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostRepository interface {
	Create(ctx context.Context, post domain.Post) (int64, error)                                                // 创建一个新的帖子
	Update(ctx context.Context, post domain.Post) error                                                         // 更新一个现有的帖子
	UpdateStatus(ctx context.Context, postId int64, status domain.PostStatus) error                             // 更新帖子的状态
	GetDraftsByAuthor(ctx context.Context, authorId int64, pagination domain.Pagination) ([]domain.Post, error) // 根据作者ID获取草稿帖子
	GetPostById(ctx context.Context, postId int64) (domain.Post, error)                                         // 根据ID获取一个帖子
	GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error)                                // 根据ID获取一个已发布的帖子
	ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)                // 获取已发布的帖子列表
	Delete(ctx context.Context, postId int64) error                                                             // 删除一个帖子
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

func (p *postRepository) Create(ctx context.Context, post domain.Post) (int64, error) {
	post.Slug = uuid.New().String()
	id, err := p.dao.Insert(ctx, fromDomainPost(post))
	if err != nil {
		p.l.Error("文章插入发生错误", zap.Error(err))
		return -1, err
	}
	return id, nil
}

func (p *postRepository) Update(ctx context.Context, post domain.Post) error {
	return p.dao.UpdateById(ctx, fromDomainPost(post))
}

func (p *postRepository) UpdateStatus(ctx context.Context, postId int64, status domain.PostStatus) error {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) GetDraftsByAuthor(ctx context.Context, authorId int64, pagination domain.Pagination) ([]domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) GetPostById(ctx context.Context, postId int64) (domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) Delete(ctx context.Context, postId int64) error {
	//TODO implement me
	panic("implement me")
}

// 将领域层对象转为dao层对象
func fromDomainPost(p domain.Post) models.Post {
	return models.Post{
		ID:           p.ID,
		Title:        p.Title,
		Content:      p.Content,
		CreateTime:   p.CreateTime,
		UpdatedTime:  p.UpdatedTime,
		Status:       p.Status.String(),
		Visibility:   p.Visibility,
		Slug:         p.Slug,
		CategoryID:   p.CategoryID,
		Tags:         p.Tags,
		CommentCount: p.CommentCount,
		ViewCount:    p.ViewCount,
	}
}

// 将dao层对象转为领域层对象
func toDomainPost(p models.Post) domain.Post {
	return domain.Post{
		ID:           p.ID,
		Title:        p.Title,
		Content:      p.Content,
		CreateTime:   p.CreateTime,
		UpdatedTime:  p.UpdatedTime,
		Visibility:   p.Visibility,
		Slug:         p.Slug,
		CategoryID:   p.CategoryID,
		Tags:         p.Tags,
		CommentCount: p.CommentCount,
		ViewCount:    p.ViewCount,
	}
}
