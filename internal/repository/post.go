package repository

import (
	"LinkMe/internal/domain"
	"LinkMe/internal/repository/cache"
	"LinkMe/internal/repository/dao"
	"LinkMe/internal/repository/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
	ListAllPost(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	GetPost(ctx context.Context, id int64) (domain.Post, error)
	GetPostCount(ctx context.Context) (int64, error)
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
		return err
	}
	return nil
}

func (p *postRepository) GetDraftsByAuthor(ctx context.Context, postId int64, uid int64) (domain.Post, error) {
	dp, err := p.dao.GetByAuthor(ctx, postId, uid)
	if err != nil {
		return domain.Post{}, err
	}
	return toDomainPost(dp), nil
}

func (p *postRepository) GetPostById(ctx context.Context, postId int64, uid int64) (domain.Post, error) {
	post, err := p.dao.GetById(ctx, postId, uid)
	if err != nil {
		return domain.Post{}, err
	}
	return toDomainPost(post), nil
}

func (p *postRepository) GetPublishedPostById(ctx context.Context, postId int64) (domain.Post, error) {
	dp, err := p.dao.GetPubById(ctx, postId)
	if err != nil {
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
		return nil, err
	}
	posts = fromDomainSlicePost(pub)
	// 如果缓存未命中，这里选择异步更新缓存
	go func() {
		ctx1, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if er := p.c.SetFirstPage(ctx1, pagination.Uid, posts); er != nil {
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
		} else if err != nil {
			p.l.Warn("获取缓存失败", zap.Error(err))
		}
	}
	// 如果缓存未命中，从数据库中获取数据
	pub, err := p.dao.ListPub(ctx, pagination)
	if err != nil {
		return nil, err
	}
	posts = fromDomainSlicePost(pub)
	if pagination.Page == 1 {
		// 由于缓存未命中，这里选择异步更新缓存
		go func() {
			ctx1, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			if er := p.c.SetPubFirstPage(ctx1, pagination.Uid, posts); er != nil {
				p.l.Warn("设置缓存失败", zap.Error(er))
			}
		}()
	}
	return posts, nil
}

func (p *postRepository) Delete(ctx context.Context, post domain.Post) error {
	// 删除缓存
	if err := p.c.DelFirstPage(ctx, post.ID); err != nil {
		return err
	}
	return p.dao.DeleteById(ctx, post)
}

func (p *postRepository) Sync(ctx context.Context, post domain.Post) (int64, error) {
	// 执行同步操作前确保为最新状态
	err := p.dao.UpdateStatus(ctx, fromDomainPost(post))
	if err != nil {
		return -1, fmt.Errorf("更新帖子状态失败: %w", err)
	}
	// 获取帖子详情，以检查状态是否发生变化
	mp, err := p.dao.GetById(ctx, post.ID, post.Author.Id)
	if err != nil {
		return -1, fmt.Errorf("获取帖子失败: %w", err)
	}
	// 检查状态是否发生变化，如果发生变化，删除缓存
	if mp.Status != post.Status {
		// 删除缓存
		if er := p.c.DelFirstPage(ctx, post.ID); er != nil {
			p.l.Warn("删除缓存失败", zap.Error(er))
		}
	}
	// 如果帖子是已发布状态，删除缓存
	if mp.Status == domain.Published {
		if er := p.c.DelPubFirstPage(ctx, post.ID); er != nil {
			p.l.Warn("删除发布缓存失败", zap.Error(er))
		}
	}
	// 执行同步操作
	id, er := p.dao.Sync(ctx, fromDomainPost(post))
	if er != nil {
		return -1, fmt.Errorf("帖子同步失败: %w", er)
	}
	// 同步完成后删除发布缓存
	if e := p.c.DelPubFirstPage(ctx, post.ID); e != nil {
		return -1, fmt.Errorf("删除发布缓存失败: %w", e)
	}
	return id, nil
}

func (p *postRepository) ListAllPost(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	posts, err := p.dao.ListAllPost(ctx, pagination)
	if err != nil {
		return nil, err
	}
	return fromDomainSlicePost(posts), nil
}

func (p *postRepository) GetPost(ctx context.Context, id int64) (domain.Post, error) {
	post, err := p.dao.GetPost(ctx, id)
	if err != nil {
		return domain.Post{}, err
	}
	return toDomainPost(post), nil
}

func (p *postRepository) GetPostCount(ctx context.Context) (int64, error) {
	count, err := p.dao.GetPostCount(ctx)
	if err != nil {
		return -1, err
	}
	return count, nil
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
		PlateID:      p.PlateID,
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
			PlateID:      repoPost.PlateID,
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
		PlateID:      post.PlateID,
		Slug:         post.Slug,
		CategoryID:   post.CategoryID,
		Tags:         post.Tags,
		CommentCount: post.CommentCount,
		Author:       domain.Author{Id: post.Author},
	}
}
