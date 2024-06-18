package models

import "time"

// User 用户信息结构体
type User struct {
	ID           int64      `gorm:"primarykey"`                          // 用户ID，主键
	CreateTime   int64      `gorm:"column:created_at;type:bigint"`       // 创建时间，Unix时间戳
	UpdatedTime  int64      `gorm:"column:updated_at;type:bigint"`       // 更新时间，Unix时间戳
	DeletedTime  int64      `gorm:"column:deleted_at;type:bigint;index"` // 删除时间，Unix时间戳，用于软删除
	Nickname     string     `gorm:"size:50"`                             // 用户昵称，最大长度50
	PasswordHash string     `gorm:"not null"`                            // 密码哈希值，不能为空
	Deleted      bool       `gorm:"column:deleted;default:false"`        // 删除标志，表示该用户是否被删除
	Birthday     *time.Time `gorm:"column:birthday;type:datetime"`       // 生日，使用datetime类型
	Email        string     `gorm:"type:varchar(100);uniqueIndex"`       // 邮箱地址，唯一
	Phone        *string    `gorm:"type:varchar(15);uniqueIndex"`        // 手机号码，唯一
	About        string     `gorm:"type:varchar(4096)"`                  // 关于用户的介绍，最大长度4096
}
