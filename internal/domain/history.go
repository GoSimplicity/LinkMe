package domain

type History struct {
	ID         int64
	PostID     int64
	Title      string
	Content    string
	ActionType string
	ActionTime int64
	AuthorID   int64
	Status     string
	Slug       string
	CategoryID int64
	Tags       string
}
