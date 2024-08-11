package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	bloom "github.com/GoSimplicity/LinkMe/pkg/cache_plug/bloom"
	"github.com/GoSimplicity/LinkMe/pkg/cache_plug/local"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostRepository interface {
	Create(ctx context.Context, post domain.Post) (uint, error)
	Update(ctx context.Context, post domain.Post) error
	UpdateStatus(ctx context.Context, post domain.Post) error
	GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error)
	GetPublishedPostById(ctx context.Context, postId uint) (domain.Post, error)
	ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	Delete(ctx context.Context, post domain.Post) error
	ListAllPost(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)
	GetPost(ctx context.Context, postId uint) (domain.Post, error)
	GetPostCount(ctx context.Context) (int64, error)
}

type postRepository struct {
	dao dao.PostDAO
	l   *zap.Logger
	cb  *bloom.CacheBloom
	cl  *local.CacheManager
}

func NewPostRepository(dao dao.PostDAO, l *zap.Logger, cb *bloom.CacheBloom, cl *local.CacheManager) PostRepository {
	return &postRepository{
		dao: dao,
		l:   l,
		cb:  cb,
		cl:  cl,
	}
}

// Create 创建帖子
func (p *postRepository) Create(ctx context.Context, post domain.Post) (uint, error) {
	// 设置帖子的唯一标识符
	post.Slug = uuid.New().String() + strconv.Itoa(int(post.ID))

	id, err := p.dao.Insert(ctx, fromDomainPost(post))
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update 更新帖子
func (p *postRepository) Update(ctx context.Context, post domain.Post) error {
	return p.dao.UpdateById(ctx, fromDomainPost(post))
}

func (p *postRepository) UpdateStatus(ctx context.Context, post domain.Post) error {
	return p.dao.UpdateStatus(ctx, fromDomainPost(post))
}

// GetPostById 获取帖子详细信息
func (p *postRepository) GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	cacheKey := fmt.Sprintf("post:detail:%d:%d", postId, uid)

	// 使用布隆过滤器查询数据
	cachedPost, err := bloom.QueryData(p.cb, ctx, cacheKey, domain.Post{}, time.Minute*10)
	if err == nil && !isEmpty(cachedPost) {
		return cachedPost, nil
	}

	// 如果缓存未命中，则从数据库获取数据
	dp, err := p.dao.GetById(ctx, postId, uid)
	if err != nil {
		// 如果数据库查询失败，缓存空对象
		go func() {
			_ = p.cb.SetEmptyCache(context.Background(), cacheKey, time.Second*10)
		}()
		return domain.Post{}, err
	}

	// 将获取到的数据异步缓存
	go func() {
		if _, cacheErr := bloom.QueryData(p.cb, context.Background(), cacheKey, toDomainPost(dp), time.Minute*10); cacheErr != nil {
			p.l.Warn("更新布隆过滤器和缓存失败", zap.Error(cacheErr))
		}
	}()

	return toDomainPost(dp), nil
}

// GetPublishedPostById 获取已发布的帖子详细信息
func (p *postRepository) GetPublishedPostById(ctx context.Context, postId uint) (domain.Post, error) {
	cacheKey := fmt.Sprintf("post:pub:detail:%d", postId)

	// 使用布隆过滤器查询数据
	cachedPost, err := bloom.QueryData(p.cb, ctx, cacheKey, domain.Post{}, time.Minute*10)
	if err == nil && !isEmpty(cachedPost) {
		return cachedPost, nil
	}

	// 如果缓存未命中，直接返回数据库查询结果，同时异步更新缓存
	dp, err := p.dao.GetPubById(ctx, postId)
	if err != nil {
		// 如果数据库查询失败，缓存空对象
		go func() {
			_ = p.cb.SetEmptyCache(context.Background(), cacheKey, time.Second*10)
		}()
		return domain.Post{}, err
	}

	// 将获取到的数据异步缓存
	go func() {
		if _, cacheErr := bloom.QueryData(p.cb, context.Background(), cacheKey, toDomainListPubPost(dp), time.Minute*10); cacheErr != nil {
			p.l.Warn("更新布隆过滤器和缓存失败", zap.Error(cacheErr))
		}
	}()

	return toDomainListPubPost(dp), nil
}

// ListPosts 获取作者帖子的列表，使用热点数据永不过期的策略
func (p *postRepository) ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	cacheKey := fmt.Sprintf("post:pri:list:%d:%d", pagination.Uid, pagination.Page)

	var cachedPosts []domain.Post

	err := p.cl.Get(ctx, cacheKey, func() (interface{}, error) {
		// 如果缓存未命中，从数据库中加载数据
		pub, err := p.dao.List(ctx, pagination)
		if err != nil {
			// 如果数据库查询失败，缓存空对象以防止缓存穿透
			_ = p.cl.SetEmptyCache(ctx, cacheKey, time.Minute*10)
			return nil, err
		}

		posts := fromDomainSlicePost(pub)

		// 在后台异步刷新缓存
		go p.refreshCacheAsync(ctx, cacheKey, posts)

		// 返回从数据库加载的数据
		return posts, nil
	}, &cachedPosts)
	if err != nil {
		p.l.Warn("获取数据失败", zap.Error(err))
		return nil, err
	}

	return cachedPosts, nil
}

// ListPublishedPosts 获取已发布的帖子列表，使用热点数据永不过期的策略
func (p *postRepository) ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	cacheKey := fmt.Sprintf("post:pub:list:%d", pagination.Page)

	var cachedPosts []domain.Post

	err := p.cl.Get(ctx, cacheKey, func() (interface{}, error) {
		// 如果缓存未命中，从数据库中加载数据
		pub, err := p.dao.ListPub(ctx, pagination)
		if err != nil {
			// 如果数据库查询失败，缓存空对象以防止缓存穿透
			_ = p.cl.SetEmptyCache(ctx, cacheKey, time.Minute*5)
			return nil, err
		}

		posts := fromDomainSlicePubPostList(pub)

		// 在后台异步刷新缓存
		go p.refreshCacheAsync(ctx, cacheKey, posts)

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
	return p.dao.DeleteById(ctx, fromDomainPost(post))
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

// refreshCacheAsync 异步刷新缓存，确保热点数据的实时性
func (p *postRepository) refreshCacheAsync(ctx context.Context, cacheKey string, data interface{}) {
	// 设置较长的过期时间
	err := p.cl.Set(ctx, cacheKey, data, 24*time.Hour)
	if err != nil {
		p.l.Warn("异步刷新缓存失败", zap.Error(err))
	}
}

// 将领域层对象转为dao层对象
func fromDomainPost(p domain.Post) dao.Post {
	return dao.Post{
		Model:        gorm.Model{ID: p.ID},
		Title:        p.Title,
		Content:      p.Content,
		AuthorID:     p.AuthorID,
		Status:       p.Status,
		PlateID:      p.PlateID,
		Slug:         p.Slug,
		CategoryID:   p.CategoryID,
		Tags:         p.Tags,
		CommentCount: p.CommentCount,
	}
}

// 将dao层对象转为领域层对象
func fromDomainSlicePubPostList(post []dao.ListPubPost) []domain.Post {
	domainPosts := make([]domain.Post, len(post)) // 创建与输入切片等长的domain.Post切片
	for i, repoPost := range post {
		domainPosts[i] = domain.Post{
			ID:           repoPost.ID,
			Title:        repoPost.Title,
			Content:      repoPost.Content,
			CreatedAt:    repoPost.CreatedAt,
			UpdatedAt:    repoPost.UpdatedAt,
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
		AuthorID:     post.AuthorID,
	}
}

// 将dao层转化为领域层
func toDomainListPubPost(post dao.ListPubPost) domain.Post {
	return domain.Post{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		Status:       post.Status,
		PlateID:      post.PlateID,
		Slug:         post.Slug,
		CategoryID:   post.CategoryID,
		Tags:         post.Tags,
		CommentCount: post.CommentCount,
		AuthorID:     post.AuthorID,
	}
}
