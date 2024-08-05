package domain

// History represents a record of actions performed on a post.
type History struct {
	ID       int64
	PostID   uint
	Title    string
	Content  string
	AuthorID int64
	Tags     string
}
