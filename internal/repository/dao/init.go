package dao

import "gorm.io/gorm"

// InitTables 初始化数据库表
func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
	)
}
