package service

import (
	"LinkMe/internal/constants"
	"LinkMe/internal/domain"
	"LinkMe/internal/domain/events/post"
	"LinkMe/internal/repository"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

type PostService interface {
	Create(ctx context.Context, post domain.Post) (int64, error)                                 // 用于创建新帖子
	Update(ctx context.Context, post domain.Post) error                                          // 用于更新现有帖子
	Publish(ctx context.Context, post domain.Post) error                                         // 用于发布帖子
	Withdraw(ctx context.Context, post domain.Post) error                                        // 用于撤回帖子
	GetDraftsByAuthor(ctx context.Context, postId int64, uid int64) (domain.Post, error)         // 获取作者的草稿
	GetPostById(ctx context.Context, postId int64, uid int64) (domain.Post, error)               // 获取特定ID的帖子
	GetPublishedPostById(ctx context.Context, postId, uid int64) (domain.Post, error)            // 获取特定ID的已发布帖子
	ListPublishedPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error) // 获取已发布的帖子列表，支持分页
	ListPosts(ctx context.Context, pagination domain.Pagination) ([]domain.Post, error)          // 获取个人帖子列表，支持分页
	Delete(ctx context.Context, postId int64, uid int64) error                                   // 删除帖子
}

type postService struct {
	repo        repository.PostRepository
	historyRepo repository.HistoryRepository
	intSvc      InteractiveService
	checkSvc    CheckService
	checkRepo   repository.CheckRepository
	producer    post.Producer
	l           *zap.Logger
}

func NewPostService(repo repository.PostRepository, l *zap.Logger, intSvc InteractiveService, checkSvc CheckService, p post.Producer, historyRepo repository.HistoryRepository, checkRepo repository.CheckRepository) PostService {
	return &postService{
		repo:        repo,
		l:           l,
		intSvc:      intSvc,
		checkSvc:    checkSvc,
		producer:    p,
		historyRepo: historyRepo,
		checkRepo:   checkRepo,
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
		p.l.Error("db sync failed", zap.Error(err))
		return err
	}
	return p.repo.Update(ctx, post)
}

// Publish 发布帖子到审核
func (p *postService) Publish(ctx context.Context, post domain.Post) error {
	// 检查帖子是否存在
	po, err := p.checkRepo.FindByID(ctx, post.ID)
	if err != nil {
		return fmt.Errorf("无法找到帖子ID为 %d 的帖子: %w", post.ID, err)
	}
	// 检查帖子状态是否允许重新提交审核
	if po.Status == constants.PostUnApproved {
		po.Status = constants.PostUnderReview
		if err := p.checkRepo.UpdateStatus(ctx, domain.Check{Status: po.Status}); err != nil {
			p.l.Error("更新审核状态失败", zap.Error(err))
			return fmt.Errorf("更新审核状态失败: %w", err)
		}
	}
	// 获取帖子详细信息
	dp, err := p.repo.GetPostById(ctx, post.ID, post.Author.Id)
	if err != nil {
		p.l.Error("获取帖子失败", zap.Error(err))
		return fmt.Errorf("获取帖子失败: %w", err)
	}
	// 提交审核
	check := domain.Check{
		PostID:  dp.ID,
		Content: dp.Content,
		Title:   dp.Title,
		UserID:  dp.Author.Id,
	}
	checkId, err := p.checkSvc.SubmitCheck(ctx, check)
	if err != nil {
		p.l.Error("提交审核失败", zap.Error(err))
		return fmt.Errorf("提交审核失败: %w", err)
	}
	// 确保 checkId 有效
	if checkId == 0 {
		p.l.Error("提交审核失败，checkId 无效", zap.Int64("postID", post.ID))
		return errors.New("提交审核失败，checkId 无效")
	}
	return nil
}

func (p *postService) Withdraw(ctx context.Context, post domain.Post) error {
	post.Status = domain.Withdrawn
	// 撤回帖子时执行同步操作,从线上库(mongodb)中移除帖子
	if _, err := p.repo.Sync(ctx, post); err != nil {
		p.l.Error("db sync failed", zap.Error(err))
		return err
	}
	return p.repo.UpdateStatus(ctx, post)
}

func (p *postService) GetDraftsByAuthor(ctx context.Context, postId int64, uid int64) (domain.Post, error) {
	dp, err := p.repo.GetDraftsByAuthor(ctx, postId, uid)
	if err != nil {
		p.l.Error("get drafts by author failed", zap.Error(err))
		return domain.Post{}, err
	}
	return dp, nil
}

func (p *postService) GetPostById(ctx context.Context, postId int64, uid int64) (domain.Post, error) {
	dp, err := p.repo.GetPostById(ctx, postId, uid)
	if err != nil {
		p.l.Error("get post failed", zap.Error(err))
		return domain.Post{}, err
	}
	return dp, err
}

func (p *postService) GetPublishedPostById(ctx context.Context, postId, uid int64) (domain.Post, error) {
	dp, err := p.repo.GetPublishedPostById(ctx, postId)
	if err != nil {
		return domain.Post{}, err // 直接返回错误
	}
	// 存入历史记录
	if er := p.historyRepo.SetHistory(ctx, dp); er != nil {
		p.l.Error("set history failed", zap.Error(er))
	}
	// 异步处理读取事件
	go func() {
		// 生成读取事件
		if er := p.producer.ProduceReadEvent(post.ReadEvent{PostId: postId, Uid: uid}); er != nil {
			p.l.Error("produce read event failed", zap.Error(er))
		}

	}()
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
		p.l.Error("delete post failed", zap.Error(err))
		return err
	}
	res := domain.Post{
		ID:     postId,
		Status: domain.Deleted,
		Author: domain.Author{
			Id: uid,
		},
	}
	if _, er := p.repo.Sync(ctx, res); er != nil {
		p.l.Error("db sync failed", zap.Error(er))
		return er
	}
	return p.repo.Delete(ctx, res)
}
