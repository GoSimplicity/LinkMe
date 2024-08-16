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

func NewRelationDAO(db *gorm.DB, l *zap.Logger) RelationDAO {
	return &relationDAO{
		db: db,
		l:  l,
	}
}

func (r *relationDAO) getCurrentTime() int64 {
	return time.Now().UnixMilli()
}

func (r *relationDAO) ListFollowerRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]Relation, error) {
	var relations []Relation
	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)
	if err := r.db.WithContext(ctx).
		Where("follower_id = ? AND status = ?", followerID, FollowStatus).
		Offset(intOffset).
		Limit(intSize).
		Find(&relations).Error; err != nil {
		return nil, err
	}
	return relations, nil
}

func (r *relationDAO) ListFolloweeRelations(ctx context.Context, followeeID int64, pagination domain.Pagination) ([]Relation, error) {
	var relations []Relation
	intSize := int(*pagination.Size)
	intOffset := int(*pagination.Offset)
	if err := r.db.WithContext(ctx).
		Where("followee_id = ? AND status = ?", followeeID, FollowStatus).
		Offset(intOffset).
		Limit(intSize).
		Find(&relations).Error; err != nil {
		return nil, err
	}
	return relations, nil
}

func (r *relationDAO) FollowUser(ctx context.Context, followerID, followeeID int64) error {
	now := r.getCurrentTime()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建Relation记录
		if err := tx.Create(&Relation{
			FollowerID: followerID,
			FolloweeID: followeeID,
			Status:     FollowStatus,
			CreatedAt:  now,
			UpdatedAt:  now,
		}).Error; err != nil {
			return err
		}

		// 更新关注者的关注数
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"followee_count": gorm.Expr("followee_count + 1"),
				"updated_at":     now,
			}),
		}).Create(&RelationCount{
			UserID:        followerID,
			FolloweeCount: 1, // 初始值设为1，确保计数器存在
			CreatedAt:     now,
			UpdatedAt:     now,
		}).Error; err != nil {
			return err
		}

		// 更新被关注者的粉丝数
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"follower_count": gorm.Expr("follower_count + 1"),
				"updated_at":     now,
			}),
		}).Create(&RelationCount{
			UserID:        followeeID,
			FollowerCount: 1, // 初始值设为1，确保计数器存在
			CreatedAt:     now,
			UpdatedAt:     now,
		}).Error; err != nil {
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

func (r *relationDAO) FollowCount(ctx context.Context, userID int64) (RelationCount, error) {
	var relationCount RelationCount

	if err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, FollowStatus).First(&relationCount).Error; err != nil {
		r.l.Error("failed to get follower count", zap.Error(err))
		return RelationCount{}, err
	}

	return relationCount, nil
}

func (r *relationDAO) CancelFollowUser(ctx context.Context, followerID, followeeID int64) error {
	now := r.getCurrentTime()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除Relation记录
		if err := tx.Where("follower_id = ? AND followee_id = ? AND status = ?", followerID, followeeID, FollowStatus).
			Delete(&Relation{}).Error; err != nil {
			return err
		}

		// 更新关注者的关注数
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"followee_count": gorm.Expr("followee_count - 1"),
				"updated_at":     now,
			}),
		}).Create(&RelationCount{
			UserID:        followerID,
			FolloweeCount: 0, // 初始值设为0，确保计数器存在
			CreatedAt:     now,
			UpdatedAt:     now,
		}).Error; err != nil {
			return err
		}

		// 更新被关注者的粉丝数
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"follower_count": gorm.Expr("follower_count - 1"),
				"updated_at":     now,
			}),
		}).Create(&RelationCount{
			UserID:        followeeID,
			FollowerCount: 0, // 初始值设为0，确保计数器存在
			CreatedAt:     now,
			UpdatedAt:     now,
		}).Error; err != nil {
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
