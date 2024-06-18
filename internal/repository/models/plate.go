package models

type Plate struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`      // 板块ID
	Name        string `gorm:"size:255;not null;uniqueIndex"` // 板块名称
	Description string `gorm:"type:text"`                     // 板块描述
	CreateTime  int64  `gorm:"column:created_at;type:bigint"` // 创建时间
	UpdatedTime int64  `gorm:"column:updated_at;type:bigint"` // 更新时间
	DeletedTime int64  `gorm:"column:deleted_at;type:bigint"` // 删除时间
	Deleted     bool   `gorm:"column:deleted;default:false"`  // 是否删除
	Uid         int64  `gorm:"index"`                         // 板主id
	Posts       []Post `gorm:"foreignKey:PlateID"`            // 帖子关系
}
