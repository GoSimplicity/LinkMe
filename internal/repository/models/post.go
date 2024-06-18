package models

type Post struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Title        string `gorm:"size:255;not null"`                            // 文章标题
	Content      string `gorm:"type:text;not null"`                           // 文章内容
	CreateTime   int64  `gorm:"column:created_at;type:bigint;not null"`       // 创建时间
	UpdatedTime  int64  `gorm:"column:updated_at;type:bigint;not null;index"` // 更新时间
	DeletedTime  int64  `gorm:"column:deleted_at;type:bigint;index"`          // 删除时间
	Deleted      bool   `gorm:"column:deleted;default:false"`                 // 是否删除
	Status       string `gorm:"size:20;default:'draft'"`                      // 文章状态，如草稿、发布等
	Author       int64  `gorm:"column:author_id;index"`                       // 用户uid
	Slug         string `gorm:"size:100;uniqueIndex"`                         // 文章的唯一标识，用于生成友好URL
	CategoryID   int64  `gorm:"index"`                                        // 关联分类表的外键
	PlateID      int64  `gorm:"index"`                                        // 关联板块表的外键
	Plate        Plate  `gorm:"foreignKey:PlateID"`                           // 板块关系
	Tags         string `gorm:"type:varchar(255);default:''"`                 // 文章标签，以逗号分隔
	CommentCount int64  `gorm:"default:0"`                                    // 文章的评论数量
}
