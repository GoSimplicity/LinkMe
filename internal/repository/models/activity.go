package models

type RecentActivity struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	UserID      int64  `gorm:"column:user_id;not null" json:"user_id"`
	Description string `gorm:"type:varchar(255);not null"`
	Time        string `gorm:"type:varchar(255);not null"`
}
