package models

// UserLikeBiz 用户点赞业务结构体
type UserLikeBiz struct {
	ID         int64 `gorm:"primaryKey;autoIncrement"`
	Uid        int64 `gorm:"index"`
	BizID      int64 `gorm:"index"`
	BizName    string
	Status     int
	UpdateTime int64
	CreateTime int64
}

// UserCollectionBiz 用户收藏业务结构体
type UserCollectionBiz struct {
	ID           int64 `gorm:"primaryKey;autoIncrement"`
	Uid          int64 `gorm:"index"`
	BizID        int64 `gorm:"index"`
	BizName      string
	CollectionId int64 `gorm:"index"`
	UpdateTime   int64
	CreateTime   int64
}

// Interactive 互动信息结构体
type Interactive struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	BizID        int64  `gorm:"uniqueIndex:biz_type_id"`
	BizName      string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCount    int64
	LikeCount    int64
	CollectCount int64
	UpdateTime   int64
	CreateTime   int64
}
