package models

type HistoryRecord struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`
	PostID     int64  `gorm:"index;not null"`    // 与帖子ID关联
	Title      string `gorm:"size:255"`          // 文章标题
	Content    string `gorm:"type:text"`         // 文章内容
	ActionType string `gorm:"size:20;not null"`  // 操作类型，例如创建、更新、删除
	ActionTime int64  `gorm:"not null"`          // 操作时间
	AuthorID   int64  `gorm:"index;not null"`    // 操作者ID
	Status     string `gorm:"size:20"`           // 帖子状态
	Slug       string `gorm:"size:100"`          // 文章的唯一标识，用于生成友好URL
	CategoryID int64  `gorm:"index"`             // 关联分类表的外键
	Tags       string `gorm:"type:varchar(255)"` // 文章标签，以逗号分隔
}
