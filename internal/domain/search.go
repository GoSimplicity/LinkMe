package domain

import "time"

type PostSearch struct {
	Id       uint
	Title    string
	AuthorId int64
	Status   uint8
	Content  string
	Tags     string
}

type UserSearch struct {
	Id       int64
	Nickname string
	Birthday time.Time
	Email    string
	Phone    string
	About    string
}

type CommentSearch struct {
	Id       uint   // 评论ID
	AuthorId int64  // 评论者ID
	Status   uint8  // 评论状态
	Content  string // 评论内容
}
