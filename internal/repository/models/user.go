package models

// User 用户信息结构体
type User struct {
	ID           int64   `gorm:"primarykey"`                          // 用户ID，主键
	CreateTime   int64   `gorm:"column:created_at;type:bigint"`       // 创建时间，Unix时间戳
	UpdatedTime  int64   `gorm:"column:updated_at;type:bigint"`       // 更新时间，Unix时间戳
	DeletedTime  int64   `gorm:"column:deleted_at;type:bigint;index"` // 删除时间，Unix时间戳，用于软删除
	PasswordHash string  `gorm:"not null"`                            // 密码哈希值，不能为空
	Deleted      bool    `gorm:"column:deleted;default:false"`        // 删除标志，表示该用户是否被删除
	Email        string  `gorm:"type:varchar(100);uniqueIndex"`       // 邮箱地址，唯一
	Phone        *string `gorm:"type:varchar(15);uniqueIndex"`        // 手机号码，唯一
	Profile      Profile `gorm:"foreignKey:UserID;references:ID"`     // 关联的用户资料
}

// Profile 用户资料信息结构体
type Profile struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`         // 用户资料ID，主键
	UserID   int64  `gorm:"not null;index"`                   // 用户ID，外键
	NickName string `gorm:"size:50"`                          // 昵称，最大长度50
	Avatar   string `gorm:"type:text"`                        // 头像URL
	About    string `gorm:"type:text"`                        // 个人简介
	Birthday string `gorm:"column:birthday;type:varchar(10)"` // 生日
}
