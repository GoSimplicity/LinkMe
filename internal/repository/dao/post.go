package dao

import (
	"LinkMe/internal/domain"
	. "LinkMe/internal/repository/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var (
	PostUpdateERROR = errors.New("ID 不对或者创作者不对")
)

type PostDAO interface {
	Insert(ctx context.Context, post Post) (int64, error)                              // 创建一个新的帖子记录
	UpdateById(ctx context.Context, post Post) error                                   // 根据ID更新一个帖子记录
	Sync(ctx context.Context, post Post) (int64, error)                                // 用于同步帖子记录
	UpdateStatus(ctx context.Context, post Post) error                                 // 更新帖子的状态
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Post, error) // 根据作者ID获取帖子记录
	GetById(ctx context.Context, id int64) (Post, error)                               // 根据ID获取一个帖子记录
	GetPubById(ctx context.Context, id int64) (PublishedPost, error)                   // 根据ID获取一个已发布的帖子记录
	ListPub(ctx context.Context, pagination domain.Pagination) ([]Post, error)         // 获取已发布的帖子记录列表
	DeleteById(ctx context.Context, post domain.Post) error
}

type postDAO struct {
	client *mongo.Client
	l      *zap.Logger
	db     *gorm.DB
}

func NewPostDAO(db *gorm.DB, l *zap.Logger, client *mongo.Client) PostDAO {
	return &postDAO{
		client: client,
		l:      l,
		db:     db,
	}
}

// Insert 创建一个新的帖子记录(mysql)
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

// UpdateById 通过Id更新帖子
func (p *postDAO) UpdateById(ctx context.Context, post Post) error {
	now := time.Now().UnixMilli()
	res := p.db.WithContext(ctx).Model(&post).Where("id = ? AND author_id = ?", post.ID, post.Author).Updates(map[string]any{
		"title":      post.Title,
		"content":    post.Content,
		"status":     post.Status,
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

// UpdateStatus 更新帖子状态
func (p *postDAO) UpdateStatus(ctx context.Context, post Post) error {
	now := time.Now().UnixMilli()
	if err := p.db.WithContext(ctx).Model(&Post{}).Where("id = ?", post.ID).
		Updates(map[string]any{
			"status":     post.Status,
			"updated_at": now,
		}).Error; err != nil {
		p.l.Error("帖子状态更新失败", zap.Error(err))
		return err
	}
	return nil
}

// Sync 同步线上库(mongodb)与制作库(mysql)
func (p *postDAO) Sync(ctx context.Context, post Post) (int64, error) {
	var mysqlPost Post
	var mongoPost Post
	now := time.Now().UnixMilli()
	post.UpdatedTime = now
	// 根据id查询帖子信息
	if err := p.db.WithContext(ctx).Where("id = ?", post.ID).First(&mysqlPost).Error; err != nil {
		return -1, err
	}
	// 只有在帖子为公开状态才会进行同步
	if post.Status == domain.Published {
		// 判断当前id的帖子是否已经被同步
		if err := p.client.Database("linkme").Collection("posts").FindOne(ctx, bson.M{"id": post.ID}).Decode(&mongoPost); err == nil {
			// 如果MongoDB中已存在相同ID的文章，则不执行同步
			return -1, errors.New("文章已存在于MongoDB")
		}
		// 如果没同步则执行同步操作
		if _, err := p.client.Database("linkme").Collection("posts").InsertOne(ctx, mysqlPost); err != nil {
			return -1, err
		}
	} else {
		// 进入到这里说明帖子状态非公开状态,或从公开状态变为非公开状态
		// 我们的mongodb数据库只储存状态为公开状态的帖子
		if _, err := p.client.Database("linkme").Collection("posts").DeleteOne(ctx, bson.M{"id": post.ID}); err != nil {
			return -1, err
		}
	}
	return mysqlPost.ID, nil
}

// GetById 根据ID获取一个帖子记录
func (p *postDAO) GetById(ctx context.Context, id int64) (Post, error) {
	var post Post
	err := p.db.WithContext(ctx).Where("id = ? AND deleted = ?", id, false).First(&post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p.l.Error("没有查到该记录", zap.Error(err))
		return Post{}, err
	}
	return post, err
}

// GetPubById 根据ID获取一个已发布的帖子记录
func (p *postDAO) GetPubById(ctx context.Context, id int64) (PublishedPost, error) {
	var pubPost PublishedPost
	status := domain.Published
	err := p.db.WithContext(ctx).Where("id = ? AND status = ?", id, status).First(&pubPost).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p.l.Error("没有查到该记录", zap.Error(err))
		return PublishedPost{}, err
	}
	return pubPost, err
}

// GetByAuthor 根据作者ID获取帖子记录
func (p *postDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Post, error) {
	var posts []Post
	err := p.db.WithContext(ctx).Where("author_id = ? AND deleted = ?", uid, false).Offset(offset).Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// ListPub 查询公开帖子
func (p *postDAO) ListPub(ctx context.Context, pagination domain.Pagination) ([]Post, error) {
	status := domain.Published
	// 设置查询超时时间
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// 指定数据库与集合
	collection := p.client.Database("linkme").Collection("posts")
	filter := bson.M{
		"status": status,
	}
	// 设置分页查询参数
	opts := options.FindOptions{
		Skip:  pagination.Offset,
		Limit: pagination.Size,
	}
	var posts []Post
	cursor, err := collection.Find(ctx, filter, &opts)
	if err != nil {
		p.l.Error("数据库查询失败", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)
	// 将获取到的查询结果解码到posts结构体中
	if err = cursor.All(ctx, &posts); err != nil {
		p.l.Error("解码查询结果失败", zap.Error(err))
		return nil, err
	}
	if len(posts) == 0 {
		p.l.Debug("查询没有返回任何结果")
	}
	return posts, nil
}

// DeleteById 通过id删除帖子
func (p *postDAO) DeleteById(ctx context.Context, post domain.Post) error {
	now := time.Now().UnixMilli()
	// 使用事务来确保操作的原子性
	tx := p.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新帖子的删除时间，这里假设您有一个表示删除状态的字段，例如 `IsDeleted`
	if err := tx.Model(Post{}).Where("id = ?", post.ID).Update("deleted_at", now).Update("status", domain.Deleted).Update("deleted", true).Error; err != nil {
		tx.Rollback()
		p.l.Error("更新帖子的删除时间出错", zap.Error(err))
		return err
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		p.l.Error("提交事务出错", zap.Error(err))
		return err
	}
	return nil
}
