package domain

// History represents a record of actions performed on a post.
type History struct {
	PostID  uint
	Title   string
	Content string
	Uid     int64
	Tags    string
}
