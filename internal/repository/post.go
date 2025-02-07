package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/job"
	"github.com/GoSimplicity/LinkMe/internal/repository/cache"
	"github.com/hibiken/asynq"

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
	dao         dao.PostDAO
	l           *zap.Logger
	cache       cache.PostCache
	asynqClient *asynq.Client
}

func NewPostRepository(dao dao.PostDAO, l *zap.Logger, cache cache.PostCache, asynqClient *asynq.Client) PostRepository {
	return &postRepository{
		dao:         dao,
		l:           l,
		cache:       cache,
		asynqClient: asynqClient,
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

	// 删除制作库列表缓存
	if err := p.cache.DelList(ctx, strconv.Itoa(int(post.ID))); err != nil {
		p.l.Error("删除帖子列表缓存失败", zap.Error(err))
	}

	return id, nil
}

// Update 更新帖子
func (p *postRepository) Update(ctx context.Context, post domain.Post) error {
	oldStatus := post.Status
	post.Status = domain.Draft

	if err := p.dao.Update(ctx, change.FromDomainPost(post)); err != nil {
		p.l.Error("更新帖子失败", zap.Error(err), zap.Uint("post_id", post.ID))
		return fmt.Errorf("更新帖子失败: %w", err)
	}

	if oldStatus == domain.Published {
		key := "*"

		if err := p.cache.DelPubList(ctx, key); err != nil {
			p.l.Error("删除帖子列表缓存失败", zap.Error(err))
		}

		if err := p.cache.DelPub(ctx, key); err != nil {
			p.l.Error("删除已发布帖子缓存失败", zap.Error(err))
		}

		// 延迟双删
		payload := job.Payload{Key: key, PostId: post.ID}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			p.l.Error("序列化任务负载失败", zap.Error(err))
		}

		task := asynq.NewTask(job.DeferRefreshPostCache, jsonPayload)
		if _, err := p.asynqClient.Enqueue(task, asynq.ProcessIn(time.Second*5)); err != nil {
			p.l.Error("延迟删除帖子缓存失败", zap.Error(err))
		}
	}

	// 删除制作库列表缓存
	if err := p.cache.DelList(ctx, strconv.Itoa(int(post.ID))); err != nil {
		p.l.Error("删除帖子列表缓存失败", zap.Error(err))
	}

	// 删除制作库缓存
	if err := p.cache.Del(ctx, int64(post.ID)); err != nil {
		p.l.Error("删除帖子缓存失败", zap.Error(err))
	}

	return nil
}

// UpdateStatus 更新帖子状态
func (p *postRepository) UpdateStatus(ctx context.Context, postId uint, uid int64, status uint8) error {
	// 更新帖子状态
	if err := p.dao.UpdateStatus(ctx, postId, uid, status); err != nil {
		p.l.Error("更新帖子状态失败",
			zap.Error(err),
			zap.Uint("post_id", postId),
			zap.Int64("uid", uid),
			zap.Uint8("status", status))
		return fmt.Errorf("更新帖子状态失败: %w", err)
	}

	// 如果不是草稿状态,需要清理缓存
	if status != domain.Draft {
		key := "*"

		if err := p.cache.DelPubList(ctx, key); err != nil {
			p.l.Error("删除帖子列表缓存失败", zap.Error(err))
		}

		// 延迟双删
		payload := job.Payload{Key: key}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			p.l.Error("序列化任务负载失败", zap.Error(err))
		}

		task := asynq.NewTask(job.DeferRefreshPostCache, jsonPayload)
		if _, err := p.asynqClient.Enqueue(task, asynq.ProcessIn(time.Second*5)); err != nil {
			p.l.Error("延迟删除帖子缓存失败", zap.Error(err))
		}
	}

	return nil
}

// GetPostById 获取帖子详细信息
func (p *postRepository) GetPostById(ctx context.Context, postId uint, uid int64) (domain.Post, error) {
	// 先从缓存中获取
	post, err := p.cache.Get(ctx, int64(postId))
	if err == nil {
		return post, nil
	}

	dp, err := p.dao.GetById(ctx, postId, uid)
	if err != nil {
		p.l.Error("获取帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		return domain.Post{}, fmt.Errorf("获取帖子失败: %w", err)
	}

	result := change.ToDomainPost(dp)

	go func() {
		if err := p.cache.Set(ctx, result); err != nil {
			p.l.Error("设置帖子缓存失败", zap.Error(err), zap.Uint("post_id", postId))
		}
	}()

	return result, nil
}

// GetPublishPostById 获取已发布的帖子详细信息
func (p *postRepository) GetPublishPostById(ctx context.Context, postId uint) (domain.Post, error) {
	// 构建缓存key
	cacheKey := fmt.Sprintf("pub:%d", postId)

	// 先从缓存中获取
	post, err := p.cache.GetPub(ctx, cacheKey)
	if err == nil {
		return post, nil
	}

	dp, err := p.dao.GetPubById(ctx, postId)
	if err != nil {
		p.l.Error("获取已发布帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		// 缓存空对象
		if err := p.cache.SetEmpty(ctx, int64(postId)); err != nil {
			p.l.Error("设置空对象缓存失败", zap.Error(err))
		}
		return domain.Post{}, err
	}

	result := change.ToDomainPubPost(dp)

	go func() {
		if err := p.cache.SetPub(ctx, cacheKey, result); err != nil {
			p.l.Error("设置已发布帖子缓存失败", zap.Error(err), zap.Uint("post_id", postId))
		}
	}()

	return result, nil
}

// ListPosts 获取作者帖子的列表
func (p *postRepository) ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	// 先从缓存中获取
	posts, err := p.cache.GetList(ctx, pagination.Page, int(*pagination.Size))
	if err == nil && len(posts) > 0 {
		return posts, nil
	}

	// 缓存未命中,从数据库获取
	pub, err := p.dao.List(ctx, pagination)
	if err != nil {
		p.l.Error("获取帖子列表失败",
			zap.Error(err),
			zap.Int("page", pagination.Page))
		return nil, fmt.Errorf("获取帖子列表失败: %w", err)
	}

	// 转换数据
	result := change.FromDomainSlicePost(pub)

	// 异步写入缓存
	go func() {
		if err := p.cache.SetList(ctx, pagination.Page, int(*pagination.Size), result); err != nil {
			p.l.Error("缓存帖子列表失败",
				zap.Error(err),
				zap.Int("page", pagination.Page))
		}
	}()

	return result, nil
}

// ListPublishPosts 获取已发布的帖子列表
func (p *postRepository) ListPublishPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) {
	// 先从缓存中获取
	posts, err := p.cache.GetPubList(ctx, pagination.Page, int(*pagination.Size))
	if err == nil && len(posts) > 0 {
		return posts, nil
	}

	// 缓存未命中,从数据库获取
	pub, err := p.dao.ListPub(ctx, pagination)
	if err != nil {
		p.l.Error("从数据库获取已发布帖子列表失败", zap.Error(err))
		return nil, fmt.Errorf("从数据库获取已发布帖子列表失败: %w", err)
	}

	if len(pub) == 0 {
		p.l.Info("没有找到已发布的帖子")
		// 缓存空对象
		if err := p.cache.SetEmpty(ctx, int64(pagination.Page)); err != nil {
			p.l.Error("设置空对象缓存失败", zap.Error(err))
		}
		return []domain.Post{}, nil
	}

	result := change.FromDomainSlicePubPostList(pub)

	go func() {
		// 如果是前三页
		if pagination.Page <= 3 {
			if err := p.cache.PreHeat(ctx, result); err != nil {
				p.l.Error("预热已发布帖子列表缓存失败", zap.Error(err))
			}
		} else {
			if err := p.cache.SetPubList(ctx, pagination.Page, int(*pagination.Size), result); err != nil {
				p.l.Error("缓存已发布帖子列表失败",
					zap.Error(err),
					zap.Int("page", pagination.Page))
			}
		}
	}()

	return result, nil
}

// Delete 删除帖子
func (p *postRepository) Delete(ctx context.Context, postId uint, uid int64) error {
	// 获取帖子
	post, err := p.GetPostById(ctx, postId, uid)
	if err != nil {
		return fmt.Errorf("获取帖子失败: %w", err)
	}

	if post.ID == 0 {
		return fmt.Errorf("帖子不存在")
	}

	// 删除帖子
	if err := p.dao.Delete(ctx, postId, uid); err != nil {
		p.l.Error("删除帖子失败", zap.Error(err), zap.Uint("post_id", postId))
		return fmt.Errorf("删除帖子失败: %w", err)
	}

	if post.Status == domain.Published {
		// 删除已发布的帖子
		if err := p.cache.DelPub(ctx, strconv.Itoa(int(postId))); err != nil {
			p.l.Error("删除已发布帖子缓存失败", zap.Error(err))
		}

		if err := p.cache.DelPubList(ctx, "*"); err != nil {
			p.l.Error("删除已发布帖子列表缓存失败", zap.Error(err))
		}

		// 延迟双删
		payload := job.Payload{Key: "*", PostId: postId}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			p.l.Error("序列化任务负载失败", zap.Error(err))
		}

		task := asynq.NewTask(job.DeferRefreshPostCache, jsonPayload)
		if _, err := p.asynqClient.Enqueue(task, asynq.ProcessIn(time.Second*5)); err != nil {
			p.l.Error("延迟删除帖子缓存失败", zap.Error(err))
		}
	}

	go func() {
		// 删除制作库列表缓存
		if err := p.cache.DelList(ctx, strconv.Itoa(int(postId))); err != nil {
			p.l.Error("删除帖子列表缓存失败", zap.Error(err))
		}

		// 删除制作库缓存
		if err := p.cache.Del(ctx, int64(postId)); err != nil {
			p.l.Error("删除帖子缓存失败", zap.Error(err))
		}
	}()

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
