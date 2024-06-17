package models

import "time"

type User struct {
	ID           int64      `gorm:"primarykey"`
	CreateTime   int64      `gorm:"column:created_at;type:bigint"`
	UpdatedTime  int64      `gorm:"column:updated_at;type:bigint"`
	DeletedTime  int64      `gorm:"column:deleted_at;type:bigint;index"`
	Nickname     string     `gorm:"size:50"`
	PasswordHash string     `gorm:"not null"`
	Deleted      bool       `gorm:"column:deleted;default:false"`
	Birthday     *time.Time `gorm:"column:birthday;type:datetime"`
	Email        string     `gorm:"type:varchar(100);uniqueIndex"`
	Phone        *string    `gorm:"type:varchar(15);uniqueIndex"`
	About        string     `gorm:"type=varchar(4096)"`
	Profile      Profile
}

type Profile struct {
	ID       int64  `gorm:"primarykey"`
	UserId   int64  `gorm:"column:not null;type:bigint"`
	NickName string `gorm:"size:50"`
	Avatar   string `gorm:"type:text"`
	Bio      string `gorm:"type:text"`
}
