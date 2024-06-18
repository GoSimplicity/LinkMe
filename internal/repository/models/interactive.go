package models

// UserLikeBiz 用户点赞业务结构体
type UserLikeBiz struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`                     // 点赞记录ID，主键，自增
	Uid        int64  `gorm:"index"`                                        // 用户ID，用于标识哪个用户点赞
	BizID      int64  `gorm:"index"`                                        // 业务ID，用于标识点赞的业务对象
	BizName    string `gorm:"type:varchar(255)"`                            // 业务名称
	Status     int    `gorm:"type:int"`                                     // 状态，用于表示点赞的状态（如有效、无效等）
	UpdateTime int64  `gorm:"column:updated_at;type:bigint;not null;index"` // 更新时间，Unix时间戳
	CreateTime int64  `gorm:"column:created_at;type:bigint"`                // 创建时间，Unix时间戳
	Deleted    bool   `gorm:"column:deleted;default:false"`                 // 删除标志，表示该记录是否被删除
}

// UserCollectionBiz 用户收藏业务结构体
type UserCollectionBiz struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`                     // 收藏记录ID，主键，自增
	Uid          int64  `gorm:"index"`                                        // 用户ID，用于标识哪个用户收藏
	BizID        int64  `gorm:"index"`                                        // 业务ID，用于标识收藏的业务对象
	BizName      string `gorm:"type:varchar(255)"`                            // 业务名称
	Status       int    `gorm:"column:status"`                                // 状态，用于表示收藏的状态（如有效、无效等）
	CollectionId int64  `gorm:"index"`                                        // 收藏ID，用于标识具体的收藏对象
	UpdateTime   int64  `gorm:"column:updated_at;type:bigint;not null;index"` // 更新时间，Unix时间戳
	CreateTime   int64  `gorm:"column:created_at;type:bigint"`                // 创建时间，Unix时间戳
	Deleted      bool   `gorm:"column:deleted;default:false"`                 // 删除标志，表示该记录是否被删除
}

// Interactive 互动信息结构体
type Interactive struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`                     // 互动记录ID，主键，自增
	BizID        int64  `gorm:"uniqueIndex:biz_type_id"`                      // 业务ID，用于标识互动的业务对象
	BizName      string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`    // 业务名称
	ReadCount    int64  `gorm:"column:read_count"`                            // 阅读数量
	LikeCount    int64  `gorm:"column:like_count"`                            // 点赞数量
	CollectCount int64  `gorm:"column:collect_count"`                         // 收藏数量
	UpdateTime   int64  `gorm:"column:updated_at;type:bigint;not null;index"` // 更新时间，Unix时间戳
	CreateTime   int64  `gorm:"column:created_at;type:bigint"`                // 创建时间，Unix时间戳
}
