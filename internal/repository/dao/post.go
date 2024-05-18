package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var (
	PostUpdateERROR = errors.New("ID 不对或者创作者不对")
)

type PostDAO interface {
	Insert(ctx context.Context, pst Post) (int64, error)                                          // 创建一个新的帖子记录
	UpdateById(ctx context.Context, pst Post) error                                               // 根据ID更新一个帖子记录
	Sync(ctx context.Context, post Post) (int64, error)                                           // 用于同步帖子记录
	SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error                      // 同步帖子的状态
	UpdateStatus(ctx context.Context, postId int64, post Post) error                              // 更新帖子的状态
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Post, error)            // 根据作者ID获取帖子记录
	GetById(ctx context.Context, id int64) (Post, error)                                          // 根据ID获取一个帖子记录
	GetPubById(ctx context.Context, id int64) (PublishedPost, error)                              // 根据ID获取一个已发布的帖子记录
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]PublishedPost, error) // 获取已发布的帖子记录列表
}

type postDAO struct {
	//client *mongo.Client
	l  *zap.Logger
	db *gorm.DB
}

func NewPostDAO(db *gorm.DB, l *zap.Logger) PostDAO {
	return &postDAO{
		//client: client,
		l:  l,
		db: db,
	}
}

func (p *postDAO) Insert(ctx context.Context, pst Post) (int64, error) {
	now := time.Now().UnixMilli()
	pst.CreateTime = now
	pst.UpdatedTime = now
	if err := p.db.WithContext(ctx).Create(&pst).Error; err != nil {
		p.l.Error("帖子插入数据库发生错误", zap.Error(err))
		return -1, err
	}
	return pst.ID, nil
}

func (p *postDAO) UpdateById(ctx context.Context, pst Post) error {
	now := time.Now().UnixMilli()
	res := p.db.WithContext(ctx).Model(&pst).Where("id = ? AND author_id = ?", pst.ID, pst.Author).Updates(map[string]any{
		"title":      pst.Title,
		"content":    pst.Content,
		"status":     pst.Status,
		"updated_at": now,
	})
	if res.Error != nil {
		p.l.Error("帖子更新失败", zap.Error(res.Error))
		return res.Error
	}
	if res.RowsAffected == 0 {
		p.l.Error("帖子更新失败", zap.Error(PostUpdateERROR))
		return PostUpdateERROR
	}
	return nil
}

func (p *postDAO) UpdateStatus(ctx context.Context, postId int64, post Post) error {
	now := time.Now().UnixMilli()
	if err := p.db.WithContext(ctx).Model(&Post{}).Where("id = ?", postId).
		Updates(map[string]any{
			"status":     post.Status,
			"updated_at": now,
		}).Error; err != nil {
		p.l.Error("帖子状态更新失败", zap.Error(err))
		return err
	}
	return nil
}

func (p *postDAO) Sync(ctx context.Context, Post Post) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postDAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func (p *postDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postDAO) GetById(ctx context.Context, id int64) (Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postDAO) GetPubById(ctx context.Context, id int64) (PublishedPost, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]PublishedPost, error) {
	//TODO implement me
	panic("implement me")
}
