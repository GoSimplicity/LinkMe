package models

// Role 表示一个角色，包含角色的基本信息
type Role struct {
	ID          int64  `gorm:"primarykey"`                          // 角色ID，主键
	Name        string `gorm:"type:varchar(100);uniqueIndex"`       // 角色名称，唯一索引
	CreateTime  int64  `gorm:"column:created_at;type:bigint"`       // 创建时间戳
	UpdatedTime int64  `gorm:"column:updated_at;type:bigint"`       // 更新时间戳
	DeletedTime int64  `gorm:"column:deleted_at;type:bigint;index"` // 删除时间戳，索引
}

// Permission 表示一个权限，包含权限的基本信息
type Permission struct {
	ID          int64  `gorm:"primarykey"`                          // 权限ID，主键
	Name        string `gorm:"type:varchar(100);uniqueIndex"`       // 权限名称，唯一索引
	CreateTime  int64  `gorm:"column:created_at;type:bigint"`       // 创建时间戳
	UpdatedTime int64  `gorm:"column:updated_at;type:bigint"`       // 更新时间戳
	DeletedTime int64  `gorm:"column:deleted_at;type:bigint;index"` // 删除时间戳，索引
}

// RolePermission 表示角色和权限的关联关系
type RolePermission struct {
	RoleID       int64 `gorm:"primarykey"`                          // 角色ID，主键
	PermissionID int64 `gorm:"primarykey"`                          // 权限ID，主键
	CreateTime   int64 `gorm:"column:created_at;type:bigint"`       // 创建时间戳
	UpdatedTime  int64 `gorm:"column:updated_at;type:bigint"`       // 更新时间戳
	DeletedTime  int64 `gorm:"column:deleted_at;type:bigint;index"` // 删除时间戳，索引
}

// UserRole 表示用户和角色的关联关系
type UserRole struct {
	UserID      int64 `gorm:"primarykey"`                          // 用户ID，主键
	RoleID      int64 `gorm:"primarykey"`                          // 角色ID，主键
	CreateTime  int64 `gorm:"column:created_at;type:bigint"`       // 创建时间戳
	UpdatedTime int64 `gorm:"column:updated_at;type:bigint"`       // 更新时间戳
	DeletedTime int64 `gorm:"column:deleted_at;type:bigint;index"` // 删除时间戳，索引
}
