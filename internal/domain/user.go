package domain

import "time"

type User struct {
	ID         uint
	Phone      *string
	Email      string
	Nickname   string
	Password   string
	Birthday   *time.Time
	CreateTime int64
}
