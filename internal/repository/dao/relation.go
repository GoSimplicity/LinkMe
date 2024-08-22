package dao

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

const (
	FollowStatus uint8 = iota // 关注
	ShieldStatus              // 屏蔽
	BlockStatus               // 拉黑
)

// RelationDAO 定义了与用户关系相关的接口
type RelationDAO interface {
	ListFollowerRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]Relation, error)
	ListFolloweeRelations(ctx context.Context, followeeID int64, pagination domain.Pagination) ([]Relation, error)
	FollowUser(ctx context.Context, followerID, followeeID int64) error
	CancelFollowUser(ctx context.Context, followerID, followeeID int64) error
	UpdateStatus(ctx context.Context, followerID, followeeID int64, status bool) error
	FollowCount(ctx context.Context, userID int64) (RelationCount, error)
}

type relationDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

// Relation 存储用户的关注数据
type Relation struct {
	ID         int64 `gorm:"column:id;primaryKey;autoIncrement"`                     // 主键ID
	FollowerID int64 `gorm:"column:follower_id;uniqueIndex:follower_id_followee_id"` // 关注者ID
	FolloweeID int64 `gorm:"column:followee_id;uniqueIndex:follower_id_followee_id"` // 被关注者ID
	Status     uint8 `gorm:"column:status"`                                          // 关系类型
	Deleted    bool  `gorm:"column:deleted"`                                         // 删除标志
	CreatedAt  int64 `gorm:"column:created_at"`                                      // 创建时间
	UpdatedAt  int64 `gorm:"column:updated_at"`                                      // 更新时间
}

// RelationCount 存储用户的粉丝和关注数量
type RelationCount struct {
	ID            int64 `gorm:"column:id;primaryKey;autoIncrement"` // 主键ID
	UserID        int64 `gorm:"column:user_id;unique"`              // 用户ID
	FollowerCount int64 `gorm:"column:follower_count"`              // 粉丝数量
	FolloweeCount int64 `gorm:"column:followee_count"`              // 关注数量
	CreatedAt     int64 `gorm:"column:created_at"`                  // 创建时间
	UpdatedAt     int64 `gorm:"column:updated_at"`                  // 更新时间
}

// NewRelationDAO 创建RelationDAO实例
func NewRelationDAO(db *gorm.DB, l *zap.Logger) RelationDAO {
	return &relationDAO{
		db: db,
		l:  l,
	}
}

// getCurrentTime 获取当前时间戳（毫秒）
func (r *relationDAO) getCurrentTime() int64 {
	return time.Now().UnixMilli()
}

// ListFollowerRelations 获取关注者的关系列表
func (r *relationDAO) ListFollowerRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]Relation, error) {
	var relations []Relation
	if err := r.db.WithContext(ctx).
		Where("follower_id = ? AND status = ?", followerID, FollowStatus).
		Offset(int(*pagination.Offset)).
		Limit(int(*pagination.Size)).
		Find(&relations).Error; err != nil {
		r.l.Error("failed to list follower relations", zap.Error(err))
		return nil, err
	}
	return relations, nil
}

// ListFolloweeRelations 获取被关注者的关系列表
func (r *relationDAO) ListFolloweeRelations(ctx context.Context, followeeID int64, pagination domain.Pagination) ([]Relation, error) {
	var relations []Relation
	if err := r.db.WithContext(ctx).
		Where("followee_id = ? AND status = ?", followeeID, FollowStatus).
		Offset(int(*pagination.Offset)).
		Limit(int(*pagination.Size)).
		Find(&relations).Error; err != nil {
		r.l.Error("failed to list followee relations", zap.Error(err))
		return nil, err
	}
	return relations, nil
}

// FollowUser 关注用户，并更新相关计数器
func (r *relationDAO) FollowUser(ctx context.Context, followerID, followeeID int64) error {
	now := r.getCurrentTime()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建Relation记录
		newRelation := &Relation{
			FollowerID: followerID,
			FolloweeID: followeeID,
			Status:     FollowStatus,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := tx.Create(newRelation).Error; err != nil {
			return err
		}

		// 更新关注者的关注数
		if err := r.updateRelationCount(tx, followerID, "followee_count", 1, now); err != nil {
			return err
		}

		// 更新被关注者的粉丝数
		if err := r.updateRelationCount(tx, followeeID, "follower_count", 1, now); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		r.l.Error("failed to follow user", zap.Error(err))
		return err
	}

	return nil
}

// CancelFollowUser 取消关注用户，并更新相关计数器
func (r *relationDAO) CancelFollowUser(ctx context.Context, followerID, followeeID int64) error {
	now := r.getCurrentTime()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除Relation记录
		if err := tx.Where("follower_id = ? AND followee_id = ? AND status = ?", followerID, followeeID, FollowStatus).
			Delete(&Relation{}).Error; err != nil {
			return err
		}

		// 更新关注者的关注数
		if err := r.updateRelationCount(tx, followerID, "followee_count", -1, now); err != nil {
			return err
		}

		// 更新被关注者的粉丝数
		if err := r.updateRelationCount(tx, followeeID, "follower_count", -1, now); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		r.l.Error("failed to cancel follow user", zap.Error(err))
		return err
	}
	return nil
}

// UpdateStatus 更新关系状态
func (r *relationDAO) UpdateStatus(ctx context.Context, followerID, followeeID int64, status bool) error {
	if err := r.db.WithContext(ctx).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Updates(map[string]any{
		"status":     status,
		"updated_at": r.getCurrentTime(),
	}).Error; err != nil {
		r.l.Error("failed to update status", zap.Error(err))
		return err
	}

	return nil
}

// FollowCount 获取用户的关注和粉丝数量
func (r *relationDAO) FollowCount(ctx context.Context, userID int64) (RelationCount, error) {
	var relationCount RelationCount

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&relationCount).Error; err != nil {
		r.l.Error("failed to get follower count", zap.Error(err))
		return RelationCount{}, err
	}

	return relationCount, nil
}

// updateRelationCount 更新关注或粉丝计数器
func (r *relationDAO) updateRelationCount(tx *gorm.DB, userID int64, field string, delta int64, timestamp int64) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			field:        gorm.Expr(field+" + ?", delta),
			"updated_at": timestamp,
		}),
	}).Create(&RelationCount{
		UserID:    userID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}).Error
}
