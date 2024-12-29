package dao

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	StatusLiked        = 1
	StatusUnliked      = 0
	StatusCollection   = 1
	StatusUnCollection = 0
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, postId uint) error
	BatchIncrReadCnt(ctx context.Context, postIds []uint) error
	InsertLikeInfo(ctx context.Context, lb UserLikeBiz) error
	DeleteLikeInfo(ctx context.Context, lb UserLikeBiz) error
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	DeleteCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetLikeInfo(ctx context.Context, postId uint, uid int64) (UserLikeBiz, error)
	GetCollectInfo(ctx context.Context, postId uint, uid int64) (UserCollectionBiz, error)
	Get(ctx context.Context, postId uint) (Interactive, error)
	GetByIds(ctx context.Context, postIds []uint) ([]Interactive, error)
}

type interactiveDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

// UserLikeBiz 用户点赞业务结构体
type UserLikeBiz struct {
	ID         int64 `gorm:"primaryKey;autoIncrement"`
	Uid        int64 `gorm:"index"`
	BizID      uint  `gorm:"index"`
	Status     int   `gorm:"type:int"`
	UpdateTime int64 `gorm:"column:updated_at;type:bigint;not null;index"`
	CreateTime int64 `gorm:"column:created_at;type:bigint"`
	Deleted    bool  `gorm:"column:deleted;default:false"`
}

// UserCollectionBiz 用户收藏业务结构体
type UserCollectionBiz struct {
	ID         int64 `gorm:"primaryKey;autoIncrement"`
	Uid        int64 `gorm:"index"`
	BizID      uint  `gorm:"index"`
	Status     int   `gorm:"column:status"`
	UpdateTime int64 `gorm:"column:updated_at;type:bigint;not null;index"`
	CreateTime int64 `gorm:"column:created_at;type:bigint"`
	Deleted    bool  `gorm:"column:deleted;default:false"`
}

// Interactive 互动信息结构体
type Interactive struct {
	ID           int64 `gorm:"primaryKey;autoIncrement"`
	BizID        uint  `gorm:"uniqueIndex"`
	ReadCount    int64 `gorm:"column:read_count"`
	LikeCount    int64 `gorm:"column:like_count"`
	CollectCount int64 `gorm:"column:collect_count"`
	UpdateTime   int64 `gorm:"column:updated_at;type:bigint;not null;index"`
	CreateTime   int64 `gorm:"column:created_at;type:bigint"`
}

func NewInteractiveDAO(db *gorm.DB, l *zap.Logger) InteractiveDAO {
	return &interactiveDAO{
		db: db,
		l:  l,
	}
}

func (i *interactiveDAO) getCurrentTime() int64 {
	return time.Now().UnixMilli()
}

// IncrReadCnt 增加阅读计数,使用UPSERT优化写入性能
func (i *interactiveDAO) IncrReadCnt(ctx context.Context, postId uint) error {
	now := i.getCurrentTime()
	return i.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_count": gorm.Expr("read_count + 1"),
			"updated_at": now,
		}),
	}).Create(&Interactive{
		BizID:      postId,
		ReadCount:  1,
		CreateTime: now,
		UpdateTime: now,
	}).Error
}

// BatchIncrReadCnt 批量增加阅读计数,使用事务提升性能
func (i *interactiveDAO) BatchIncrReadCnt(ctx context.Context, postIds []uint) error {
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txInc := NewInteractiveDAO(tx, i.l)
		for _, postId := range postIds {
			if err := txInc.IncrReadCnt(ctx, postId); err != nil {
				i.l.Error("增加阅读计数失败", zap.Error(err))
				return err
			}
		}
		return nil
	})
}

// InsertLikeInfo 插入点赞信息,使用事务保证数据一致性
func (i *interactiveDAO) InsertLikeInfo(ctx context.Context, lb UserLikeBiz) error {
	now := i.getCurrentTime()
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingLike UserLikeBiz
		err := tx.Where("uid = ? AND biz_id = ?", lb.Uid, lb.BizID).First(&existingLike).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			lb.CreateTime = now
			lb.UpdateTime = now
			lb.Status = StatusLiked
			if err = tx.Create(&lb).Error; err != nil {
				i.l.Error("创建点赞记录失败", zap.Error(err))
				return err
			}
		} else if err != nil {
			i.l.Error("查询点赞记录失败", zap.Error(err))
			return err
		} else {
			if err = tx.Model(&existingLike).Updates(map[string]interface{}{
				"status":     StatusLiked,
				"updated_at": now,
			}).Error; err != nil {
				i.l.Error("更新点赞记录失败", zap.Error(err))
				return err
			}
		}

		// 更新互动计数
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_count": gorm.Expr("like_count + 1"),
				"updated_at": now,
			}),
		}).Create(&Interactive{
			BizID:      lb.BizID,
			LikeCount:  1,
			UpdateTime: now,
			CreateTime: now,
		}).Error
	})
}

// DeleteLikeInfo 删除点赞信息
func (i *interactiveDAO) DeleteLikeInfo(ctx context.Context, lb UserLikeBiz) error {
	now := i.getCurrentTime()
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询用户点赞记录并检查状态
		var likeBiz UserLikeBiz
		if err := tx.Where("uid = ? AND biz_id = ?", lb.Uid, lb.BizID).First(&likeBiz).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				i.l.Error("用户未点赞,无法取消", zap.Error(err))
				return errors.New("用户未点赞,无法取消")
			}
			i.l.Error("查询点赞记录失败", zap.Error(err))
			return err
		}

		if likeBiz.Status == StatusUnliked {
			i.l.Error("点赞已取消,请勿重复操作")
			return errors.New("点赞已取消,请勿重复操作")
		}

		// 分别更新点赞状态和互动计数
		if err := tx.Model(&UserLikeBiz{}).
			Where("uid = ? AND biz_id = ?", lb.Uid, lb.BizID).
			Updates(map[string]interface{}{
				"status":     StatusUnliked,
				"updated_at": now,
			}).Error; err != nil {
			i.l.Error("更新点赞状态失败", zap.Error(err))
			return err
		}

		if err := tx.Model(&Interactive{}).
			Where("biz_id = ?", lb.BizID).
			Updates(map[string]interface{}{
				"like_count": gorm.Expr("CASE WHEN like_count > 0 THEN like_count - 1 ELSE 0 END"),
				"updated_at": now,
			}).Error; err != nil {
			i.l.Error("更新互动计数失败", zap.Error(err))
			return err
		}

		return nil
	})
}

// InsertCollectionBiz 插入收藏信息
func (i *interactiveDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := i.getCurrentTime()
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingCollection UserCollectionBiz
		err := tx.Where("uid = ? AND biz_id = ?", cb.Uid, cb.BizID).First(&existingCollection).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			cb.CreateTime = now
			cb.UpdateTime = now
			cb.Status = StatusCollection
			if err = tx.Create(&cb).Error; err != nil {
				i.l.Error("创建收藏记录失败", zap.Error(err))
				return err
			}
		} else if err != nil {
			i.l.Error("查询收藏记录失败", zap.Error(err))
			return err
		} else {
			if err = tx.Model(&existingCollection).Updates(map[string]interface{}{
				"status":     StatusCollection,
				"updated_at": now,
			}).Error; err != nil {
				i.l.Error("更新收藏记录失败", zap.Error(err))
				return err
			}
		}

		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_count": gorm.Expr("collect_count + 1"),
				"updated_at":    now,
			}),
		}).Create(&Interactive{
			BizID:        cb.BizID,
			CollectCount: 1,
			UpdateTime:   now,
			CreateTime:   now,
		}).Error
	})
}

// DeleteCollectionBiz 删除收藏信息
func (i *interactiveDAO) DeleteCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := i.getCurrentTime()
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询用户收藏记录并检查状态
		var collectionBiz UserCollectionBiz
		if err := tx.Where("uid = ? AND biz_id = ?", cb.Uid, cb.BizID).First(&collectionBiz).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				i.l.Error("用户未收藏,无法取消", zap.Error(err))
				return errors.New("用户未收藏,无法取消")
			}
			i.l.Error("查询收藏记录失败", zap.Error(err))
			return err
		}

		if collectionBiz.Status == StatusUnCollection {
			i.l.Error("收藏已取消,请勿重复操作")
			return errors.New("收藏已取消,请勿重复操作")
		}

		// 分别更新收藏状态和互动计数
		if err := tx.Model(&UserCollectionBiz{}).
			Where("uid = ? AND biz_id = ?", cb.Uid, cb.BizID).
			Updates(map[string]interface{}{
				"status":     StatusUnCollection,
				"updated_at": now,
			}).Error; err != nil {
			i.l.Error("更新收藏状态失败", zap.Error(err))
			return err
		}

		if err := tx.Model(&Interactive{}).
			Where("biz_id = ?", cb.BizID).
			Updates(map[string]interface{}{
				"collect_count": gorm.Expr("CASE WHEN collect_count > 0 THEN collect_count - 1 ELSE 0 END"),
				"updated_at":    now,
			}).Error; err != nil {
			i.l.Error("更新互动计数失败", zap.Error(err))
			return err
		}

		return nil
	})
}

// GetLikeInfo 获取点赞信息
func (i *interactiveDAO) GetLikeInfo(ctx context.Context, postId uint, uid int64) (UserLikeBiz, error) {
	var lb UserLikeBiz
	err := i.db.WithContext(ctx).
		Where("uid = ? AND biz_id = ? AND status = ?", uid, postId, StatusLiked).
		First(&lb).Error
	return lb, err
}

// GetCollectInfo 获取收藏信息
func (i *interactiveDAO) GetCollectInfo(ctx context.Context, postId uint, uid int64) (UserCollectionBiz, error) {
	var cb UserCollectionBiz
	err := i.db.WithContext(ctx).
		Where("uid = ? AND biz_id = ? AND status = ?", uid, postId, StatusCollection).
		First(&cb).Error
	return cb, err
}

// Get 获取单个互动信息
func (i *interactiveDAO) Get(ctx context.Context, postId uint) (Interactive, error) {
	var inc Interactive
	err := i.db.WithContext(ctx).Where("biz_id = ?", postId).First(&inc).Error
	return inc, err
}

// GetByIds 批量获取互动信息
func (i *interactiveDAO) GetByIds(ctx context.Context, postIds []uint) ([]Interactive, error) {
	var inc []Interactive
	err := i.db.WithContext(ctx).Where("biz_id IN ?", postIds).Find(&inc).Error
	return inc, err
}
