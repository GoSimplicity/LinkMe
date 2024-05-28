package models

// UserLikeBiz 用户点赞业务结构体
type UserLikeBiz struct {
	ID         int64 `gorm:"primaryKey;autoIncrement"`
	Uid        int64 `gorm:"index"`
	BizID      int64 `gorm:"index"`
	BizName    string
	Status     int
	UpdateTime int64 `gorm:"column:updated_at;type:bigint;not null;index"`
	CreateTime int64 `gorm:"column:created_at;type:bigint"`
	Deleted    bool  `gorm:"column:deleted;default:false"`
}

// UserCollectionBiz 用户收藏业务结构体
type UserCollectionBiz struct {
	ID           int64 `gorm:"primaryKey;autoIncrement"`
	Uid          int64 `gorm:"index"`
	BizID        int64 `gorm:"index"`
	BizName      string
	Status       int   `gorm:"column:status"`
	CollectionId int64 `gorm:"index"`
	UpdateTime   int64 `gorm:"column:updated_at;type:bigint;not null;index"`
	CreateTime   int64 `gorm:"column:created_at;type:bigint"`
	Deleted      bool  `gorm:"column:deleted;default:false"`
}

// Interactive 互动信息结构体
type Interactive struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	BizID        int64  `gorm:"uniqueIndex:biz_type_id"`
	BizName      string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCount    int64  `gorm:"column:read_count"`
	LikeCount    int64  `gorm:"column:like_count"`
	CollectCount int64  `gorm:"column:collect_count"`
	UpdateTime   int64  `gorm:"column:updated_at;type:bigint;not null;index"`
	CreateTime   int64  `gorm:"column:created_at;type:bigint"`
}
