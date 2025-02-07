package domain

type PostSearch struct {
	Id       uint
	Title    string
	AuthorId int64
	Status   uint8
	Content  string
	Tags     []string
}

type UserSearch struct {
	Id       int64
	Username string
	RealName string
	Phone    *string
}

type CommentSearch struct {
	Id       uint   // 评论ID
	AuthorId int64  // 评论者ID
	Status   uint8  // 评论状态
	Content  string // 评论内容
}
