package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type PostRepository interface {
	Create(ctx context.Context, post domain.Post) (int64, error)                                                // 创建一个新的帖子
	Update(ctx context.Context, post domain.Post) error                                                         // 更新一个现有的帖子
	UpdateStatus(ctx context.Context, post domain.Post) error                                                   // 更新帖子的状态
	GetDraftsByAuthor(ctx context.Context, authorId int64, pagination domain.Pagination) ([]domain.Post, error) // 根据作者ID获取草稿帖子
	GetPostById(ctx context.Context, postId int64) (domain.Post, error)                                         // 根据ID获取一个帖子
	GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error)                                // 根据ID获取一个已发布的帖子
	ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)                // 获取已发布的帖子列表
	Delete(ctx context.Context, postId int64) error                                                             // 删除一个帖子
	Sync(ctx context.Context, post domain.Post) (int64, error)                                                  // 用于同步帖子记录

}

type postRepository struct {
	dao dao.PostDAO
	l   *zap.Logger
	c   cache.PostCache
}

func NewPostRepository(dao dao.PostDAO, l *zap.Logger, c cache.PostCache) PostRepository {
	return &postRepository{
		dao: dao,
		l:   l,
		c:   c,
	}
}

func (p *postRepository) Create(ctx context.Context, post domain.Post) (int64, error) {
	post.Slug = uuid.New().String()
	id, err := p.dao.Insert(ctx, fromDomainPost(post))
	if err != nil {
		p.l.Error("文章插入发生错误", zap.Error(err))
		return -1, err
	}
	// 删除缓存
	if er := p.c.DelFirstPage(ctx, id); er != nil {
		p.l.Warn("删除缓存失败", zap.Error(er))
		return -1, er
	}
	return id, nil
}

func (p *postRepository) Update(ctx context.Context, post domain.Post) error {
	// 删除缓存
	if err := p.c.DelFirstPage(ctx, post.ID); err != nil {
		p.l.Warn("删除缓存失败", zap.Error(err))
		return err
	}
	// 更新数据库
	if err := p.dao.UpdateById(ctx, fromDomainPost(post)); err != nil {
		p.l.Error("更新帖子失败", zap.Error(err))
		return err
	}
	return nil
}

func (p *postRepository) UpdateStatus(ctx context.Context, post domain.Post) error {
	now := time.Now().UnixMilli()
	post.UpdatedTime = now
	// 删除缓存
	if err := p.c.DelFirstPage(ctx, post.ID); err != nil {
		p.l.Warn("删除缓存失败", zap.Error(err))
		return err
	}
	// 更新数据库
	if err := p.dao.UpdateStatus(ctx, fromDomainPost(post)); err != nil {
		p.l.Error("更新帖子状态失败", zap.Error(err))
		return err
	}
	return nil
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
	// 尝试从缓存中获取第一页的帖子摘要
	posts, err := p.c.GetFirstPage(ctx, pagination.Uid)
	if err == nil && posts != nil {
		// 如果缓存命中，直接返回缓存中的数据
		return posts, nil
	}
	// 如果缓存未命中，从数据库中获取数据
	pub, err := p.dao.ListPub(ctx, pagination)
	if err != nil {
		p.l.Error("公开文章获取失败", zap.Error(err))
		return nil, err
	}
	posts = toDomainPost(pub)
	// 由于缓存未命中，这里选择更新缓存
	if er := p.c.SetFirstPage(ctx, pagination.Uid, posts); er != nil {
		p.l.Warn("缓存设置失败", zap.Error(er))
	}
	return posts, nil
}

func (p *postRepository) Delete(ctx context.Context, postId int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) Sync(ctx context.Context, post domain.Post) (int64, error) {
	// 执行同步操作前确保为最新状态
	err := p.dao.UpdateStatus(ctx, fromDomainPost(post))
	if err != nil {
		return -1, err
	}
	// 获取帖子详情，以检查状态是否发生变化
	mp, err := p.dao.GetById(ctx, post.ID)
	if err != nil {
		p.l.Error("获取post失败", zap.Error(err))
		return -1, err
	}
	// 检查状态是否发生变化，如果发生变化，删除缓存
	if mp.Status != post.Status {
		// 删除缓存
		if er := p.c.DelFirstPage(ctx, post.ID); er != nil {
			p.l.Warn("删除缓存失败", zap.Error(er))
		}
	}
	// 同步操作
	id, err := p.dao.Sync(ctx, fromDomainPost(post))
	if err != nil {
		p.l.Error("同步post失败", zap.Error(err))
		return -1, err
	}
	return id, nil
}

// 将领域层对象转为dao层对象
func fromDomainPost(p domain.Post) models.Post {
	return models.Post{
		ID:           p.ID,
		Title:        p.Title,
		Content:      p.Content,
		CreateTime:   p.CreateTime,
		UpdatedTime:  p.UpdatedTime,
		Author:       p.Author.Id,
		Status:       p.Status,
		Visibility:   p.Visibility,
		Slug:         p.Slug,
		CategoryID:   p.CategoryID,
		Tags:         p.Tags,
		CommentCount: p.CommentCount,
		ViewCount:    p.ViewCount,
	}
}

// 将dao层对象转为领域层对象
func toDomainPost(p []models.Post) []domain.Post {
	domainPosts := make([]domain.Post, len(p)) // 创建与输入切片等长的domain.Post切片
	for i, repoPost := range p {
		domainPosts[i] = domain.Post{
			ID:           repoPost.ID,
			Title:        repoPost.Title,
			Content:      repoPost.Content,
			CreateTime:   repoPost.CreateTime,
			UpdatedTime:  repoPost.UpdatedTime,
			Status:       repoPost.Status,
			Visibility:   repoPost.Visibility,
			Slug:         repoPost.Slug,
			CategoryID:   repoPost.CategoryID,
			Tags:         repoPost.Tags,
			CommentCount: repoPost.CommentCount,
			ViewCount:    repoPost.ViewCount,
		}
	}
	return domainPosts
}
