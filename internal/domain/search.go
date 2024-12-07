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
	Email    string
	Nickname string
	Phone    *string
}
