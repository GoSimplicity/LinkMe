package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	bloom "github.com/GoSimplicity/LinkMe/pkg/cachep/bloom"
	"github.com/GoSimplicity/LinkMe/pkg/cachep/local"
	"gorm.io/gorm"

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
		p.l.Error("创建帖子失败", zap.Error(err))
		return 0, fmt.Errorf("创建帖子失败: %w", err)
	}

	return id, nil
}

// Update 更新帖子
func (p *postRepository) Update(ctx context.Context, post domain.Post) error {
	if err := p.dao.UpdateById(ctx, fromDomainPost(post)); err != nil {
		p.l.Error("更新帖子失败", zap.Error(err), zap.Uint("post_id", post.ID))
		return fmt.Errorf("更新帖子失败: %w", err)
	}
	return nil
}

// UpdateStatus 更新帖子状态
func (p *postRepository) UpdateStatus(ctx context.Context, post domain.Post) error {
	if err := p.dao.UpdateStatus(ctx, fromDomainPost(post)); err != nil {
		p.l.Error("更新帖子状态失败", zap.Error(err), zap.Uint("post_id", post.ID))
		return fmt.Errorf("更新帖子状态失败: %w", err)
	}
	return nil
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
			if err := p.cb.SetEmptyCache(ctx, cacheKey, time.Second*10); err != nil {
				p.l.Error("设置空缓存失败", zap.Error(err))
			}
		}()

		p.l.Error("获取帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		return domain.Post{}, fmt.Errorf("获取帖子失败: %w", err)
	}

	// 将获取到的数据异步缓存
	go func(ctx context.Context) {
		newCtx := context.Background()
		if err := p.cb.SetEmptyCache(newCtx, cacheKey, time.Second*10); err != nil {
			p.l.Error("设置空缓存失败", zap.Error(err))
		}
	}(ctx)

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
			if err := p.cb.SetEmptyCache(context.Background(), cacheKey, time.Second*10); err != nil {
				p.l.Error("设置空缓存失败", zap.Error(err))
			}
		}()
		p.l.Error("获取已发布帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		return domain.Post{}, fmt.Errorf("获取已发布帖子失败: %w", err)
	}

	// 将获取到的数据异步缓存
	go func() {
		if _, cacheErr := bloom.QueryData(p.cb, context.Background(), cacheKey, toDomainListPubPost(dp), time.Minute*10); cacheErr != nil {
			p.l.Error("更新布隆过滤器和缓存失败", zap.Error(cacheErr))
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
			if err := p.cl.SetEmptyCache(ctx, cacheKey, time.Minute*10); err != nil {
				p.l.Error("设置空缓存失败", zap.Error(err))
			}
			p.l.Error("获取帖子列表失败", zap.Error(err))
			return nil, fmt.Errorf("获取帖子列表失败: %w", err)
		}

		posts := fromDomainSlicePost(pub)
		return posts, nil
	}, &cachedPosts)
	if err != nil {
		p.l.Error("获取数据失败", zap.Error(err))
		return nil, fmt.Errorf("获取数据失败: %w", err)
	}

	return cachedPosts, nil
}

// ListPublishedPosts 获取已发布的帖子列表
func (p *postRepository) ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	cacheKey := fmt.Sprintf("post:pub:list:%d", pagination.Page)
	var cachedPosts []domain.Post

	err := p.cl.Get(ctx, cacheKey, func() (interface{}, error) {
		// 如果缓存未命中，从数据库中加载数据
		pub, err := p.dao.ListPub(ctx, pagination)
		if err != nil {
			p.l.Error("从数据库获取已发布帖子列表失败", zap.Error(err))
			// 如果数据库查询失败，缓存空对象以防止缓存穿透,使用较短的过期时间
			if err := p.cl.SetEmptyCache(ctx, cacheKey, time.Minute*5); err != nil {
				p.l.Error("设置空缓存失败", zap.Error(err))
			}
			return nil, fmt.Errorf("从数据库获取已发布帖子列表失败: %w", err)
		}

		if len(pub) == 0 {
			p.l.Info("没有找到已发布的帖子")
			return []domain.Post{}, nil
		}

		posts := fromDomainSlicePubPostList(pub)
		return posts, nil
	}, &cachedPosts)

	if err != nil {
		p.l.Error("获取已发布帖子列表失败",
			zap.Error(err),
			zap.Int("page", pagination.Page))
		return nil, fmt.Errorf("获取已发布帖子列表失败: %w", err)
	}

	return cachedPosts, nil
}

// Delete 删除帖子
func (p *postRepository) Delete(ctx context.Context, post domain.Post) error {
	if err := p.dao.DeleteById(ctx, fromDomainPost(post)); err != nil {
		p.l.Error("删除帖子失败", zap.Error(err), zap.Uint("post_id", post.ID))
		return fmt.Errorf("删除帖子失败: %w", err)
	}
	return nil
}

// ListAllPost 列出所有帖子
func (p *postRepository) ListAllPost(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	posts, err := p.dao.ListAllPost(ctx, pagination)
	if err != nil {
		p.l.Error("获取所有帖子失败", zap.Error(err))
		return nil, fmt.Errorf("获取所有帖子失败: %w", err)
	}
	return fromDomainSlicePost(posts), nil
}

// GetPost 获取帖子
func (p *postRepository) GetPost(ctx context.Context, postId uint) (domain.Post, error) {
	post, err := p.dao.GetPost(ctx, postId)
	if err != nil {
		p.l.Error("获取帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		return domain.Post{}, fmt.Errorf("获取帖子失败: %w", err)
	}

	return toDomainPost(post), nil
}

// GetPostCount 获取帖子数量
func (p *postRepository) GetPostCount(ctx context.Context) (int64, error) {
	count, err := p.dao.GetPostCount(ctx)
	if err != nil {
		p.l.Error("获取帖子数量失败", zap.Error(err))
		return 0, fmt.Errorf("获取帖子数量失败: %w", err)
	}

	return count, nil
}

// isEmpty 判断帖子是否为空
func isEmpty(post domain.Post) bool {
	return reflect.DeepEqual(post, domain.Post{})
}

// fromDomainPost 将领域层对象转为dao层对象
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

// fromDomainSlicePubPostList 将dao层对象转为领域层对象
func fromDomainSlicePubPostList(post []dao.ListPubPost) []domain.Post {
	domainPosts := make([]domain.Post, len(post))
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

// fromDomainSlicePost 将dao层对象转为领域层对象
func fromDomainSlicePost(post []dao.Post) []domain.Post {
	domainPosts := make([]domain.Post, len(post))
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

// toDomainPost 将dao层转化为领域层
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

// toDomainListPubPost 将dao层转化为领域层
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
