package domain

import "time"

type User struct {
	ID         int64
	Phone      *string
	Email      string
	Nickname   string
	Password   string
	Birthday   *time.Time
	CreateTime int64
	About      string
	Profile    Profile
}

type Profile struct {
	ID       int64
	UserID   int64
	Avatar   string
	NickName string
	Bio      string
}
