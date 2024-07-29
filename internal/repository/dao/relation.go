package dao

import (
	"context"
	"github.com/GoSimplicity/LinkMe/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	FollowStatus uint8 = iota
	ShieldStatus
	BlockStatus
)

type RelationDAO interface {
	ListRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]Relation, error)
	GetRelationInfo(ctx context.Context, followerID, followeeID int64) (Relation, error)
	FollowUser(ctx context.Context, followerID, followeeID int64) error
	UpdateStatus(ctx context.Context, followerID, followeeID int64, status bool) error
	FollowerCount(ctx context.Context, userID int64) (int64, error)
	FolloweeCount(ctx context.Context, userID int64) (int64, error)
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

// RelationStatus 存储关注关系的状态
type RelationStatus struct {
	ID           int64 `gorm:"column:id;primaryKey;autoIncrement"` // 主键ID
	FollowerID   int64 `gorm:"column:follower_id"`                 // 关注者ID
	FolloweeID   int64 `gorm:"column:followee_id"`                 // 被关注者ID
	RelationType uint8 `gorm:"column:relation_type"`               // 关系类型
}

func NewRelationDAO(db *gorm.DB, l *zap.Logger) RelationDAO {
	return &relationDAO{
		db: db,
		l:  l,
	}
}

func (r relationDAO) ListRelations(ctx context.Context, followerID int64, pagination domain.Pagination) ([]Relation, error) {
	// TODO 实现方法
	panic("implement me")
}

func (r relationDAO) GetRelationInfo(ctx context.Context, followerID, followeeID int64) (Relation, error) {
	// TODO 实现方法
	panic("implement me")
}

func (r relationDAO) FollowUser(ctx context.Context, followerID, followeeID int64) error {
	// TODO 实现方法
	panic("implement me")
}

func (r relationDAO) UpdateStatus(ctx context.Context, followerID, followeeID int64, status bool) error {
	// TODO 实现方法
	panic("implement me")
}

func (r relationDAO) FollowerCount(ctx context.Context, userID int64) (int64, error) {
	// TODO 实现方法
	panic("implement me")
}

func (r relationDAO) FolloweeCount(ctx context.Context, userID int64) (int64, error) {
	// TODO 实现方法
	panic("implement me")
}
