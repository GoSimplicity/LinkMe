package models

type Check struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`                     // 审核ID
	PostID    int64  `gorm:"not null"`                                     // 帖子ID
	Content   string `gorm:"type:text;not null"`                           // 审核内容
	Title     string `gorm:"size:255;not null"`                            // 审核标签
	Author    int64  `gorm:"column:author_id;index"`                       // 提交审核的用户ID
	Status    string `gorm:"size:20;not null;default:'Pending'"`           // 审核状态
	Remark    string `gorm:"type:text"`                                    // 审核备注
	CreatedAt int64  `gorm:"column:created_at;type:bigint;not null"`       // 创建时间
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;not null;index"` // 更新时间
}
