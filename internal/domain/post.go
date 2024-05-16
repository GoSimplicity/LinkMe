package domain

type Post struct {
	ID           int64
	UserID       int64
	Title        string
	Content      string
	CreateTime   int64
	UpdatedTime  int64
	Status       string
	Visibility   string
	Slug         string
	CategoryID   int64
	Tags         string
	CommentCount int64
	ViewCount    int64
}
