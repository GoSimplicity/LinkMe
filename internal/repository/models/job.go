package models

type Job struct {
	Id          int64  `gorm:"primaryKey,autoIncrement"`
	Name        string `gorm:"type:varchar(128);unique"`
	Executor    string
	Expression  string
	Cfg         string
	Status      int
	Version     int
	NextTime    int64 `gorm:"index"`
	CreateTime  int64 `gorm:"column:created_at;type:bigint;not null"`
	UpdatedTime int64 `gorm:"column:updated_at;type:bigint;not null"`
}
