package dao

import (
	. "LinkMe/internal/repository/models"
	"context"
	"errors"
	"log"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrRecordNotFound  = gorm.ErrRecordNotFound
	StatusLiked        = 1
	StatusUnliked      = 0
	StatusCollection   = 1
	StatusUnCollection = 0
)

// InteractiveDAO 互动数据访问对象接口
type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, biz []string, bizIds []int64) error
	InsertLikeInfo(ctx context.Context, lb UserLikeBiz) error
	DeleteLikeInfo(ctx context.Context, lb UserLikeBiz) error
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	DeleteCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error)
	GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error)
	Get(ctx context.Context, biz string, id int64) (Interactive, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error)
}

type interactiveDAO struct {
	db *gorm.DB
	l  *zap.Logger
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

// IncrReadCnt 增加阅读计数
func (i *interactiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := i.getCurrentTime()
	// 创建Interactive实例，用于存储阅读计数更新
	interactive := Interactive{
		BizName:    biz,
		BizID:      bizId,
		ReadCount:  1,
		CreateTime: now,
		UpdateTime: now,
	}
	// 使用Clauses处理数据库冲突，即在记录已存在时更新阅读计数
	return i.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_count": gorm.Expr("read_count + 1"),
			"updated_at": now,
		}),
	}).Create(&interactive).Error
}

func (i *interactiveDAO) BatchIncrReadCnt(ctx context.Context, biz []string, bizIds []int64) error {
	//TODO implement me
	panic("implement me")
}

// InsertLikeInfo 插入用户点赞信息和相关互动信息
func (i *interactiveDAO) InsertLikeInfo(ctx context.Context, lb UserLikeBiz) error {
	now := i.getCurrentTime()
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查点赞是否存在
		existingLike := UserLikeBiz{}
		err := tx.Where("uid = ? AND biz_id = ? AND biz_name = ?", lb.Uid, lb.BizID, lb.BizName).First(&existingLike).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				i.l.Error("query record filed", zap.Error(err))
				return err
			}
			// 创建新的点赞记录
			lb.CreateTime = now
			lb.UpdateTime = now
			lb.Status = StatusLiked
			if er := tx.Create(&lb).Error; er != nil {
				i.l.Error("create user like record filed ")
				return er
			}
		} else {
			// 更新现有的点赞记录
			existingLike.Status = StatusLiked
			existingLike.UpdateTime = now
			if er := tx.Save(&existingLike).Error; er != nil {
				i.l.Error("update user like record filed")
				return er
			}
		}
		// 更新互动信息中的点赞计数
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_count": gorm.Expr("`like_count` + 1"),
				"updated_at": now,
			}),
		}).Create(&Interactive{
			BizName:    lb.BizName,
			BizID:      lb.BizID,
			LikeCount:  1,
			UpdateTime: now,
			CreateTime: now,
		}).Error
	})
}

// DeleteLikeInfo 删除用户点赞信息和相关互动信息
func (i *interactiveDAO) DeleteLikeInfo(ctx context.Context, lb UserLikeBiz) error {
	now := i.getCurrentTime()
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 如果出现冲突，则更新updated_at和status字段
		err := tx.WithContext(ctx).Model(&UserLikeBiz{}).Where("uid = ? AND biz_id = ? AND biz_name = ?", lb.Uid, lb.BizID, lb.BizName).Updates(map[string]interface{}{
			"updated_at": now,
			"status":     StatusUnliked,
		}).Error
		if err != nil {
			i.l.Error("delete user like filed")
			return err
		}
		return tx.WithContext(ctx).Model(&Interactive{}).Where("biz_name = ? AND biz_id = ?", lb.BizName, lb.BizID).Updates(map[string]interface{}{
			"like_count": gorm.Expr("`like_count` - 1"),
			"updated_at": now,
		}).Error
	})
}

// InsertCollectionBiz 插入用户收藏信息和相关互动信息
func (i *interactiveDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := i.getCurrentTime()
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查用户是否已经收藏了该业务
		existingCollection := UserCollectionBiz{}
		err := tx.Where("uid = ? AND biz_id = ? AND biz_name = ?", cb.Uid, cb.BizID, cb.BizName).First(&existingCollection).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				i.l.Error("query record failed", zap.Error(err))
				return err
			}
			// 创建新的收藏记录
			cb.CreateTime = now
			cb.UpdateTime = now
			cb.Status = StatusCollection
			if er := tx.Create(&cb).Error; er != nil {
				i.l.Error("create user collection record failed", zap.Error(er))
				return er
			}
			log.Println(cb.Status)
		} else {
			// 更新现有的收藏记录
			existingCollection.Status = StatusCollection
			existingCollection.UpdateTime = now
			if er := tx.Save(&existingCollection).Error; er != nil {
				i.l.Error("update user collection record failed", zap.Error(er))
				return er
			}
			log.Println(cb.Status)
		}
		// 更新互动记录的收藏数
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_count": gorm.Expr("`collect_count` + 1"),
				"updated_at":    now,
			}),
		}).Create(&Interactive{
			BizName:      cb.BizName,
			BizID:        cb.BizID,
			CollectCount: 1,
			UpdateTime:   now,
			CreateTime:   now,
		}).Error
	})
}

// DeleteCollectionBiz 删除用户收藏信息和相关互动信息
func (i *interactiveDAO) DeleteCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := i.getCurrentTime()
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新用户收藏信息
		err := tx.WithContext(ctx).Model(&UserCollectionBiz{}).Where("uid = ? AND biz_id = ? AND biz_name = ? AND collection_id = ?", cb.Uid, cb.BizID, cb.BizName, cb.CollectionId).Updates(map[string]interface{}{
			"updated_at": now,
			"status":     StatusUnCollection,
		}).Error
		if err != nil {
			i.l.Error("delete user collection filed")
			return err
		}
		// 更新互动信息
		return tx.WithContext(ctx).Model(&Interactive{}).Where("biz_name = ? AND biz_id = ?", cb.BizName, cb.BizID).Updates(map[string]interface{}{
			"collect_count": gorm.Expr("`collect_count` - 1"),
			"updated_at":    now,
		}).Error
	})
}

func (i *interactiveDAO) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error) {
	var lb UserLikeBiz
	err := i.db.WithContext(ctx).Where("uid = ? AND biz_name = ? AND biz_id = ? AND status = ?", uid, biz, id, 1).First(&lb).Error
	return lb, err
}

func (i *interactiveDAO) GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error) {
	var cb UserCollectionBiz
	err := i.db.WithContext(ctx).Where("uid = ? AND biz_name = ? AND biz_id = ? AND status = ?", uid, biz, id, 1).First(&cb).Error
	return cb, err
}

func (i *interactiveDAO) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	var inc Interactive
	err := i.db.WithContext(ctx).Where("biz_name = ? AND biz_id = ?", biz, id).First(&inc).Error
	return inc, err
}

func (i *interactiveDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	var inc []Interactive
	err := i.db.WithContext(ctx).
		Where("biz_name = ? AND biz_id IN ?", biz, ids).
		Find(&inc).Error
	return inc, err
}
