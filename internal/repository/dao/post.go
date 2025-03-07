package dao

import (
	"context"
	"errors"
	"time"

	"github.com/GoSimplicity/LinkMe/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrPostNotFound  = errors.New("post not found")
	ErrInvalidParams = errors.New("invalid parameters")
	ErrPlateNotFound = errors.New("plate not found")
)

type PostDAO interface {
	Insert(ctx context.Context, post Post) (uint, error)
	Update(ctx context.Context, post Post) error
	UpdateStatus(ctx context.Context, postId uint, uid int64, status uint8) error
	GetById(ctx context.Context, postId uint, uid int64) (Post, error)
	GetPubById(ctx context.Context, postId uint) (PubPost, error)
	ListPub(ctx context.Context, pagination domain.Pagination) ([]PubPost, error)
	List(ctx context.Context, pagination domain.Pagination) ([]Post, error)
	Delete(ctx context.Context, postId uint, uid int64) error
	ListAll(ctx context.Context, pagination domain.Pagination) ([]Post, error)
	GetPost(ctx context.Context, postId uint) (Post, error)
	GetPostsCount(ctx context.Context) (int64, error)
}

type postDAO struct {
	l  *zap.Logger
	db *gorm.DB
}

type Post struct {
	gorm.Model
	Title        string `gorm:"size:255;not null"`            // 帖子标题
	Content      string `gorm:"type:text;not null"`           // 帖子内容
	Status       uint8  `gorm:"default:0"`                    // 帖子状态 
	Uid          int64  `gorm:"column:uid;index"`             // 作者ID
	Slug         string `gorm:"size:100;uniqueIndex"`         // 唯一标识
	CategoryID   int64  `gorm:"index"`                        // 分类ID
	PlateID      int64  `gorm:"index"`                        // 板块ID
	Plate        Plate  `gorm:"foreignKey:PlateID"`           // 关联板块
	Tags         string `gorm:"type:varchar(255);default:''"` // 标签
	CommentCount int64  `gorm:"default:0"`                    // 评论数
	IsSubmit     bool   `gorm:"default:false"`                // 是否提交审核
}

type PubPost struct {
	gorm.Model
	Title        string `gorm:"size:255;not null"`            // 帖子标题
	Content      string `gorm:"type:text;not null"`           // 帖子内容
	Status       uint8  `gorm:"default:0"`                    // 帖子状态
	Uid          int64  `gorm:"column:uid;index"`             // 作者ID
	Slug         string `gorm:"size:100;uniqueIndex"`         // 唯一标识
	CategoryID   int64  `gorm:"index"`                        // 分类ID
	PlateID      int64  `gorm:"index"`                        // 板块ID
	Plate        Plate  `gorm:"foreignKey:PlateID"`           // 关联板块
	Tags         string `gorm:"type:varchar(255);default:''"` // 标签
	CommentCount int64  `gorm:"default:0"`                    // 评论数
}

func NewPostDAO(db *gorm.DB, l *zap.Logger) PostDAO {
	return &postDAO{
		l:  l,
		db: db,
	}
}

// Insert 插入新帖子
func (p *postDAO) Insert(ctx context.Context, post Post) (uint, error) {
	if post.PlateID <= 0 {
		return 0, ErrInvalidParams
	}

	// 使用事务确保数据一致性
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查板块是否存在
		var count int64
		if err := tx.Model(&Plate{}).Where("id = ?", post.PlateID).Count(&count).Error; err != nil {
			p.l.Error("检查板块是否存在失败", zap.Error(err))
			return err
		}
		if count == 0 {
			return ErrPlateNotFound
		}
		// 创建帖子
		if err := tx.Create(&post).Error; err != nil {
			p.l.Error("创建帖子失败", zap.Error(err))
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return post.ID, nil
}

// Update 更新帖子信息
func (p *postDAO) Update(ctx context.Context, post Post) error {
	if post.ID == 0 || post.Uid == 0 {
		return ErrInvalidParams
	}

	// 更新帖子基本信息
	res := p.db.WithContext(ctx).Model(&Post{}).Where("id = ? AND uid = ?", post.ID, post.Uid).Updates(map[string]interface{}{
		"title":       post.Title,
		"content":     post.Content,
		"plate_id":    post.PlateID,
		"status":      post.Status,
		"is_submit":   post.IsSubmit,
		"category_id": post.CategoryID,
		"tags":        post.Tags,
		"updated_at":  time.Now(),
	})

	if res.Error != nil {
		p.l.Error("更新帖子失败", zap.Error(res.Error))
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}

// UpdateStatus 更新帖子状态
func (p *postDAO) UpdateStatus(ctx context.Context, postId uint, uid int64, status uint8) error {
	if postId == 0 || uid == 0 {
		return ErrInvalidParams
	}

	// 使用事务处理状态更新
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新原帖子状态
		updates := map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}

		// 如果是撤销状态或审核未通过(草稿状态)，同时更新 is_submit
		if status == domain.Withdrawn || status == domain.Draft {
			updates["is_submit"] = false
		}

		res := tx.Model(&Post{}).Where("id = ? AND uid = ?", postId, uid).Updates(updates)
		if res.Error != nil {
			p.l.Error("更新帖子状态失败", zap.Error(res.Error))
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrPostNotFound
		}

		// 如果是发布状态,创建已发布帖子副本
		if status == domain.Published {
			var post Post
			if err := tx.Where("id = ?", postId).First(&post).Error; err != nil {
				return err
			}

			// 使用 REPLACE INTO 语法，避免先删除再插入
			pubPost := PubPost{
				Model:        post.Model,
				Title:        post.Title,
				Content:      post.Content,
				Status:       post.Status,
				Uid:          post.Uid,
				Slug:         post.Slug,
				CategoryID:   post.CategoryID,
				PlateID:      post.PlateID,
				Plate:        post.Plate,
				Tags:         post.Tags,
				CommentCount: post.CommentCount,
			}

			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true, // 更新所有字段
			}).Create(&pubPost).Error; err != nil {
				return err
			}
		} else if status == domain.Withdrawn {
			// 如果是撤销状态,删除已发布的帖子
			if err := tx.Where("id = ?", postId).Delete(&PubPost{}).Error; err != nil {
				p.l.Error("删除已发布帖子失败", zap.Error(err))
				return err
			}
		}
		return nil
	})

	return err
}

// GetById 根据ID获取帖子
func (p *postDAO) GetById(ctx context.Context, postId uint, uid int64) (Post, error) {
	if postId == 0 || uid == 0 {
		return Post{}, ErrInvalidParams
	}

	var post Post
	err := p.db.WithContext(ctx).Where("uid = ? AND id = ?", uid, postId).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.l.Debug("未找到帖子", zap.Uint("id", postId), zap.Int64("uid", uid))
			return Post{}, ErrPostNotFound
		}
		p.l.Error("获取帖子失败", zap.Error(err))
		return Post{}, err
	}
	return post, nil
}

// List 获取帖子列表
func (p *postDAO) List(ctx context.Context, pagination domain.Pagination) ([]Post, error) {
	if pagination.Size == nil || pagination.Offset == nil {
		return nil, ErrInvalidParams
	}

	var posts []Post
	query := p.db.WithContext(ctx).Model(&Post{})
	// 如果指定了用户ID，则只查询该用户的帖子
	if pagination.Uid > 0 {
		query = query.Where("uid = ?", pagination.Uid)
	}

	err := query.Limit(int(*pagination.Size)).Offset(int(*pagination.Offset)).Find(&posts).Error
	if err != nil {
		p.l.Error("获取帖子列表失败", zap.Error(err))
		return nil, err
	}
	return posts, nil
}

// GetPubById 获取已发布的帖子
func (p *postDAO) GetPubById(ctx context.Context, postId uint) (PubPost, error) {
	if postId == 0 {
		return PubPost{}, ErrInvalidParams
	}

	var post PubPost
	err := p.db.WithContext(ctx).Where("id = ?", postId).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			p.l.Debug("未找到已发布的帖子", zap.Error(err))
			return PubPost{}, ErrPostNotFound
		}
		p.l.Error("获取已发布帖子失败", zap.Error(err))
		return PubPost{}, err
	}
	return post, nil
}

// ListPub 获取已发布帖子列表
func (p *postDAO) ListPub(ctx context.Context, pagination domain.Pagination) ([]PubPost, error) {
	if pagination.Size == nil || pagination.Offset == nil {
		return nil, ErrInvalidParams
	}

	var posts []PubPost
	err := p.db.WithContext(ctx).Model(&PubPost{}).Limit(int(*pagination.Size)).Offset(int(*pagination.Offset)).Find(&posts).Error
	if err != nil {
		p.l.Error("获取已发布帖子列表失败", zap.Error(err))
		return nil, err
	}
	return posts, nil
}

// Delete 删除帖子
func (p *postDAO) Delete(ctx context.Context, postId uint, uid int64) error {
	if postId == 0 || uid == 0 {
		return ErrInvalidParams
	}

	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除已发布的帖子
		if err := tx.Where("id = ? AND uid = ?", postId, uid).Delete(&PubPost{}).Error; err != nil {
			p.l.Error("删除已发布帖子失败", zap.Error(err))
			return err
		}

		// 软删除原帖子
		if err := tx.Where("id = ? AND uid = ?", postId, uid).Delete(&Post{}).Error; err != nil {
			p.l.Error("删除帖子失败", zap.Error(err))
			return err
		}

		return nil
	})
}

// GetPost 获取帖子
func (p *postDAO) GetPost(ctx context.Context, postId uint) (Post, error) {
	if postId == 0 {
		return Post{}, ErrInvalidParams
	}

	var post Post
	if err := p.db.WithContext(ctx).First(&post, postId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Post{}, ErrPostNotFound
		}
		p.l.Error("获取帖子失败", zap.Error(err))
		return Post{}, err
	}

	return post, nil
}

// ListAll 获取所有帖子
func (p *postDAO) ListAll(ctx context.Context, pagination domain.Pagination) ([]Post, error) {
	if pagination.Size == nil || pagination.Offset == nil {
		return nil, ErrInvalidParams
	}

	var posts []Post
	err := p.db.WithContext(ctx).Model(&Post{}).
		Limit(int(*pagination.Size)).
		Offset(int(*pagination.Offset)).
		Find(&posts).Error
	if err != nil {
		p.l.Error("获取所有帖子列表失败", zap.Error(err))
		return nil, err
	}

	return posts, nil
}

// GetPostsCount 获取帖子总数
func (p *postDAO) GetPostsCount(ctx context.Context) (int64, error) {
	var count int64
	if err := p.db.WithContext(ctx).Model(&Post{}).Count(&count).Error; err != nil {
		p.l.Error("获取帖子总数失败", zap.Error(err))
		return 0, err
	}
	return count, nil
}
