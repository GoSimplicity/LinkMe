package domain

// History represents a record of actions performed on a post.
type History struct {
	ID       int64
	PostID   int64
	Title    string
	Content  string
	Deleted  bool
	AuthorID int64
	Tags     string
}
