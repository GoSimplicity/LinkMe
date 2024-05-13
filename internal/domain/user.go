package domain

import "time"

type User struct {
	ID         uint
	Phone      *string
	Email      *string
	Nickname   string
	Username   string
	Password   string
	Birthday   time.Time
	CreateTime time.Time
}
