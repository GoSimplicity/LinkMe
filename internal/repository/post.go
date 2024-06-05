package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PostRepository 帖子仓库接口
type PostRepository interface {
	Create(ctx context.Context, post domain.Post) (int64, error)
	Update(ctx context.Context, post domain.Post) error
	UpdateStatus(ctx context.Context, post domain.Post) error
	GetDraftsByAuthor(ctx context.Context, postId int64, uid int64) (domain.Post, error)
	GetPostById(ctx context.Context, postId int64, uid int64) (domain.Post, error)
	GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error)
	ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	Delete(ctx context.Context, post domain.Post) error
	Sync(ctx context.Context, post domain.Post) (int64, error)
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
		p.l.Error("post insert failed", zap.Error(err))
		return -1, err
	}
	// 删除缓存
	if er := p.c.DelFirstPage(ctx, id); er != nil {
		p.l.Warn("delete cache failed", zap.Error(er))
	}
	return id, nil
}

func (p *postRepository) Update(ctx context.Context, post domain.Post) error {
	// 删除缓存
	if er := p.c.DelFirstPage(ctx, post.ID); er != nil {
		p.l.Warn("delete cache failed", zap.Error(er))
	}
	// 更新数据库
	if err := p.dao.UpdateById(ctx, fromDomainPost(post)); err != nil {
		p.l.Error("update post failed", zap.Error(err))
		return err
	}
	return nil
}

func (p *postRepository) UpdateStatus(ctx context.Context, post domain.Post) error {
	now := time.Now().UnixMilli()
	post.UpdatedTime = now
	// 删除缓存
	if er := p.c.DelFirstPage(ctx, post.ID); er != nil {
		p.l.Warn("delete cache failed", zap.Error(er))
	}
	// 更新数据库
	if err := p.dao.UpdateStatus(ctx, fromDomainPost(post)); err != nil {
		p.l.Error("update post status failed", zap.Error(err))
		return err
	}
	return nil
}

func (p *postRepository) GetDraftsByAuthor(ctx context.Context, postId int64, uid int64) (domain.Post, error) {
	dp, err := p.dao.GetByAuthor(ctx, postId, uid)
	if err != nil {
		p.l.Error("get post failed by uid", zap.Error(err))
		return domain.Post{}, err
	}
	return toDomainPost(dp), nil
}

func (p *postRepository) GetPostById(ctx context.Context, postId int64, uid int64) (domain.Post, error) {
	post, err := p.dao.GetById(ctx, postId, uid)
	if err != nil {
		p.l.Error("get post failed by id", zap.Error(err))
		return domain.Post{}, err
	}
	return toDomainPost(post), nil
}

func (p *postRepository) GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error) {
	dp, err := p.dao.GetPubById(ctx, postId)
	if err != nil {
		p.l.Error("get pub post failed by id", zap.Error(err))
		return domain.Post{}, err
	}
	return toDomainPost(dp), nil
}

func (p *postRepository) ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	var posts []domain.Post
	var err error
	if pagination.Page == 1 {
		posts, err = p.c.GetFirstPage(ctx, pagination.Uid)
		if err == nil && len(posts) != 0 {
			// 如果缓存命中，直接返回缓存中的数据
			return posts, nil
		}
	}
	// 如果缓存未命中，从数据库中获取数据
	pub, err := p.dao.List(ctx, pagination)
	if err != nil {
		p.l.Error("get pub post failed", zap.Error(err))
		return nil, err
	}
	posts = fromDomainSlicePost(pub)
	// 如果缓存未命中，这里选择异步更新缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if er := p.c.SetFirstPage(ctx, pagination.Uid, posts); er != nil {
			p.l.Warn("set cache failed", zap.Error(er))
		}
	}()
	return posts, nil
}

func (p *postRepository) ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	var posts []domain.Post
	var err error
	if pagination.Page == 1 {
		// 尝试从缓存中获取第一页的帖子摘要
		posts, err = p.c.GetPubFirstPage(ctx, pagination.Uid)
		if err == nil && len(posts) != 0 {
			// 如果缓存命中，直接返回缓存中的数据
			return posts, nil
		}
	}
	// 如果缓存未命中，从数据库中获取数据
	pub, err := p.dao.ListPub(ctx, pagination)
	if err != nil {
		p.l.Error("get pub post failed", zap.Error(err))
		return nil, err
	}
	posts = fromDomainSlicePost(pub)
	// 由于缓存未命中，这里选择异步更新缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if er := p.c.SetPubFirstPage(ctx, pagination.Uid, posts); er != nil {
			p.l.Warn("set cache failed", zap.Error(er))
		}
	}()
	return posts, nil
}

func (p *postRepository) Delete(ctx context.Context, post domain.Post) error {
	// 删除缓存
	if err := p.c.DelFirstPage(ctx, post.ID); err != nil {
		p.l.Warn("delete cache failed", zap.Error(err))
		return err
	}
	return p.dao.DeleteById(ctx, post)
}

func (p *postRepository) Sync(ctx context.Context, post domain.Post) (int64, error) {
	// 执行同步操作前确保为最新状态
	err := p.dao.UpdateStatus(ctx, fromDomainPost(post))
	if err != nil {
		return -1, err
	}
	// 获取帖子详情，以检查状态是否发生变化
	mp, err := p.dao.GetById(ctx, post.ID, post.Author.Id)
	if err != nil {
		p.l.Error("get post failed", zap.Error(err))
		return -1, err
	}
	// 检查状态是否发生变化，如果发生变化，删除缓存
	if mp.Status != post.Status {
		// 删除缓存
		if er := p.c.DelFirstPage(ctx, post.ID); er != nil {
			p.l.Warn("delete cache failed", zap.Error(er))
		}
	}
	if mp.Status == domain.Published {
		if er := p.c.DelPubFirstPage(ctx, post.ID); er != nil {
			p.l.Warn("delete cache failed", zap.Error(er))
		}
	}
	// 同步操作
	id, err := p.dao.Sync(ctx, fromDomainPost(post))
	if err != nil {
		p.l.Error("post sync failed", zap.Error(err))
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
		Slug:         p.Slug,
		CategoryID:   p.CategoryID,
		Tags:         p.Tags,
		CommentCount: p.CommentCount,
	}
}

// 将dao层对象转为领域层对象
func fromDomainSlicePost(post []models.Post) []domain.Post {
	domainPosts := make([]domain.Post, len(post)) // 创建与输入切片等长的domain.Post切片
	for i, repoPost := range post {
		domainPosts[i] = domain.Post{
			ID:           repoPost.ID,
			Title:        repoPost.Title,
			Content:      repoPost.Content,
			CreateTime:   repoPost.CreateTime,
			UpdatedTime:  repoPost.UpdatedTime,
			Status:       repoPost.Status,
			Slug:         repoPost.Slug,
			CategoryID:   repoPost.CategoryID,
			Tags:         repoPost.Tags,
			CommentCount: repoPost.CommentCount,
		}
	}
	return domainPosts
}

// 将dao层转化为领域层
func toDomainPost(post models.Post) domain.Post {
	return domain.Post{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		CreateTime:   post.CreateTime,
		UpdatedTime:  post.UpdatedTime,
		Status:       post.Status,
		Slug:         post.Slug,
		CategoryID:   post.CategoryID,
		Tags:         post.Tags,
		CommentCount: post.CommentCount,
	}
}
