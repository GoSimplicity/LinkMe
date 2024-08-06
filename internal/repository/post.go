package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	bloom "github.com/GoSimplicity/LinkMe/pkg/cache_plug/bloom"
	"github.com/GoSimplicity/LinkMe/pkg/cache_plug/local"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"reflect"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostRepository interface {
	Create(ctx context.Context, post domain.Post) (uint, error)
	Update(ctx context.Context, post domain.Post) error
	UpdateStatus(ctx context.Context, post domain.Post) error
	GetDraftsByAuthor(ctx context.Context, postId uint, uid int64) (domain.Post, error)
	GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error)
	GetPublishedPostById(ctx context.Context, postId uint) (domain.Post, error)
	ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	Delete(ctx context.Context, post domain.Post) error
	Sync(ctx context.Context, post domain.Post) (uint, error)
	ListAllPost(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	GetPost(ctx context.Context, postId uint) (domain.Post, error)
	GetPostCount(ctx context.Context) (int64, error)
}

type postRepository struct {
	dao dao.PostDAO
	l   *zap.Logger
	c   cache.PostCache
	cb  *bloom.CacheBloom
	cl  *local.CacheManager
}

func NewPostRepository(dao dao.PostDAO, l *zap.Logger, c cache.PostCache, cb *bloom.CacheBloom, cl *local.CacheManager) PostRepository {
	return &postRepository{
		dao: dao,
		l:   l,
		c:   c,
		cb:  cb,
		cl:  cl,
	}
}

// Create 创建帖子
func (p *postRepository) Create(ctx context.Context, post domain.Post) (uint, error) {
	post.Slug = uuid.New().String()
	id, err := p.dao.Insert(ctx, fromDomainPost(post))
	if err != nil {
		return 0, err
	}
	// 删除缓存
	if er := p.c.DelFirstPage(ctx, id); er != nil {
		p.l.Warn("delete cache failed", zap.Error(er))
	}
	return id, nil
}

// Update 更新帖子
func (p *postRepository) Update(ctx context.Context, post domain.Post) error {
	// 更新数据库
	if err := p.dao.UpdateById(ctx, fromDomainPost(post)); err != nil {
		return err
	}
	// 删除缓存
	if er := p.c.DelFirstPage(ctx, post.ID); er != nil {
		p.l.Warn("delete cache failed", zap.Error(er))
	}
	return nil
}

func (p *postRepository) UpdateStatus(ctx context.Context, post domain.Post) error {
	// 更新数据库
	if err := p.dao.UpdateStatus(ctx, fromDomainPost(post)); err != nil {
		return err
	}
	// 删除缓存
	if er := p.c.DelFirstPage(ctx, post.ID); er != nil {
		p.l.Warn("delete cache failed", zap.Error(er))
	}
	return nil
}

// GetDraftsByAuthor 获取作者帖子详情
func (p *postRepository) GetDraftsByAuthor(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	// 定义缓存键
	cacheKey := fmt.Sprintf("post:draft:%d:%d", postId, uid)
	// 使用布隆过滤器查询数据
	cachedPost, err := bloom.QueryData(p.cb, ctx, cacheKey, domain.Post{}, time.Minute*10)
	if err == nil && !isEmpty(cachedPost) {
		// 如果缓存命中，直接返回缓存中的数据
		return cachedPost, nil
	}
	// 如果缓存未命中，从数据库中获取数据
	dp, err := p.dao.GetByAuthor(ctx, postId, uid)
	if err != nil {
		// 如果数据库查询失败，缓存空对象
		_ = p.cb.SetEmptyCache(ctx, cacheKey, time.Minute*10)
		return domain.Post{}, err
	}
	// 将数据存储到布隆过滤器与缓存中
	cachedPost, err = bloom.QueryData(p.cb, ctx, cacheKey, toDomainPost(dp), time.Minute*10)
	if err != nil {
		p.l.Warn("更新布隆过滤器和缓存失败", zap.Error(err))
	}
	return cachedPost, nil
}

// GetPostById 获取帖子详细信息
func (p *postRepository) GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	// 定义缓存键
	cacheKey := fmt.Sprintf("post:detail:%d:%d", postId, uid)
	// 使用布隆过滤器查询数据
	cachedPost, err := bloom.QueryData(p.cb, ctx, cacheKey, domain.Post{}, time.Minute*10)
	if err == nil && !isEmpty(cachedPost) {
		// 如果缓存命中，直接返回缓存中的数据
		return cachedPost, nil
	}
	// 如果缓存未命中，从数据库中获取数据
	post, err := p.dao.GetById(ctx, postId, uid)
	if err != nil {
		// 如果数据库查询失败，缓存空对象
		_ = p.cb.SetEmptyCache(ctx, cacheKey, time.Minute*10)
		return domain.Post{}, err
	}
	// 将数据存储到布隆过滤器与缓存中
	cachedPost, err = bloom.QueryData(p.cb, ctx, cacheKey, toDomainPost(post), time.Minute*10)
	if err != nil {
		p.l.Warn("更新布隆过滤器和缓存失败", zap.Error(err))
	}
	return cachedPost, nil
}

// GetPublishedPostById 获取已发布的帖子详细信息
func (p *postRepository) GetPublishedPostById(ctx context.Context, postId uint) (domain.Post, error) {
	// 定义缓存键
	cacheKey := fmt.Sprintf("post:pub:detail:%d", postId)
	// 使用布隆过滤器查询数据
	cachedPost, err := bloom.QueryData(p.cb, ctx, cacheKey, domain.Post{}, time.Minute*10)
	if err == nil && !isEmpty(cachedPost) {
		// 如果缓存命中，直接返回缓存中的数据
		return cachedPost, nil
	}
	// 如果缓存未命中，从数据库中获取数据
	dp, err := p.dao.GetPubById(ctx, postId)
	if err != nil {
		// 如果数据库查询失败，缓存空对象
		_ = p.cb.SetEmptyCache(ctx, cacheKey, time.Minute*10)
		return domain.Post{}, err
	}
	// 将数据存储到布隆过滤器与缓存中
	cachedPost, err = bloom.QueryData(p.cb, ctx, cacheKey, toDomainPost(dp), time.Minute*10)
	if err != nil {
		p.l.Warn("更新布隆过滤器和缓存失败", zap.Error(err))
	}
	return cachedPost, nil
}

func (p *postRepository) ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	// 定义缓存键
	cacheKey := fmt.Sprintf("post:pri:%d:%d", pagination.Uid, pagination.Page)
	// 尝试从缓存中获取数据
	var cachedPosts []domain.Post
	err := p.cl.Get(ctx, cacheKey, func() (interface{}, error) {
		// 如果缓存未命中，从数据库中加载数据
		pub, err := p.dao.List(ctx, pagination)
		if err != nil {
			// 如果数据库查询失败，缓存空对象以防止缓存穿透
			_ = p.cl.SetEmptyCache(ctx, cacheKey, time.Minute*10)
			return nil, err
		}
		// 将从数据库加载的数据转换为 domain.Post 类型
		posts := fromDomainSlicePost(pub)
		return posts, nil
	}, &cachedPosts)
	if err != nil {
		p.l.Warn("获取数据失败", zap.Error(err))
		return nil, err
	}
	return cachedPosts, nil
}

func (p *postRepository) ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	// 定义缓存键
	cacheKey := fmt.Sprintf("post:pub:%d:%d", pagination.Uid, pagination.Page)
	// 尝试从缓存中获取数据
	var cachedPosts []domain.Post
	err := p.cl.Get(ctx, cacheKey, func() (interface{}, error) {
		// 如果缓存未命中，从数据库中加载数据
		pub, err := p.dao.ListPub(ctx, pagination)
		if err != nil {
			// 如果数据库查询失败，缓存空对象以防止缓存穿透
			_ = p.cl.SetEmptyCache(ctx, cacheKey, time.Minute*5)
			return nil, err
		}
		// 将从数据库加载的数据转换为 domain.Post 类型
		posts := fromDomainSlicePost(pub)
		return posts, nil
	}, &cachedPosts)
	if err != nil {
		p.l.Warn("获取数据失败", zap.Error(err))
		return nil, err
	}
	return cachedPosts, nil
}

// Delete 删除帖子
func (p *postRepository) Delete(ctx context.Context, post domain.Post) error {
	// 删除缓存
	if err := p.c.DelFirstPage(ctx, post.ID); err != nil {
		return err
	}
	return p.dao.DeleteById(ctx, fromDomainPost(post))
}

// Sync 同步帖子状态
func (p *postRepository) Sync(ctx context.Context, post domain.Post) (uint, error) {
	err := p.dao.UpdateStatus(ctx, fromDomainPost(post))
	if err != nil {
		return 0, fmt.Errorf("update post status failed: %w", err)
	}
	mp, err := p.dao.GetById(ctx, post.ID, post.Author.Id)
	if err != nil {
		return 0, fmt.Errorf("get post failed: %w", err)
	}
	eg, ctx := errgroup.WithContext(ctx)
	if mp.Status != post.Status {
		eg.Go(func() error {
			return p.c.DelFirstPage(ctx, post.ID)
		})
	}
	if mp.Status == domain.Published {
		eg.Go(func() error {
			return p.c.DelPubFirstPage(ctx, post.ID)
		})
	}
	id, err := p.dao.Sync(ctx, fromDomainPost(post))
	if err != nil {
		return 0, fmt.Errorf("sync post failed: %w", err)
	}
	eg.Go(func() error {
		return p.c.DelPubFirstPage(ctx, post.ID)
	})
	if er := eg.Wait(); er != nil {
		return 0, er
	}
	return id, nil
}

// ListAllPost 列出所有帖子
func (p *postRepository) ListAllPost(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	posts, err := p.dao.ListAllPost(ctx, pagination)
	if err != nil {
		return nil, err
	}
	return fromDomainSlicePost(posts), nil
}

// GetPost 获取帖子
func (p *postRepository) GetPost(ctx context.Context, postId uint) (domain.Post, error) {
	post, err := p.dao.GetPost(ctx, postId)
	if err != nil {
		return domain.Post{}, err
	}
	return toDomainPost(post), nil
}

// GetPostCount 获取帖子数量
func (p *postRepository) GetPostCount(ctx context.Context) (int64, error) {
	count, err := p.dao.GetPostCount(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func isEmpty(post domain.Post) bool {
	return reflect.DeepEqual(post, domain.Post{})
}

// 将领域层对象转为dao层对象
func fromDomainPost(p domain.Post) dao.Post {
	return dao.Post{
		Model:        gorm.Model{ID: p.ID},
		Title:        p.Title,
		Content:      p.Content,
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
func fromDomainSlicePost(post []dao.Post) []domain.Post {
	domainPosts := make([]domain.Post, len(post)) // 创建与输入切片等长的domain.Post切片
	for i, repoPost := range post {
		domainPosts[i] = domain.Post{
			ID:           repoPost.ID,
			Title:        repoPost.Title,
			Content:      repoPost.Content,
			CreatedAt:    repoPost.CreatedAt,
			UpdatedAt:    repoPost.UpdatedAt,
			DeletedAt:    sql.NullTime(repoPost.DeletedAt),
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
func toDomainPost(post dao.Post) domain.Post {
	return domain.Post{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		DeletedAt:    sql.NullTime(post.DeletedAt),
		Status:       post.Status,
		PlateID:      post.PlateID,
		Slug:         post.Slug,
		CategoryID:   post.CategoryID,
		Tags:         post.Tags,
		CommentCount: post.CommentCount,
		Author:       domain.Author{Id: post.Author},
	}
}
