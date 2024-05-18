package domain

type Author struct {
	Id   int64
	Name string
}
type Post struct {
	ID           int64
	UserID       int64
	Title        string
	Content      string
	CreateTime   int64
	UpdatedTime  int64
	Author       Author
	Status       PostStatus
	Visibility   string
	Slug         string
	CategoryID   int64
	Tags         string
	CommentCount int64
	ViewCount    int64
}

type Pagination struct {
	Page int    // 当前页码
	Size *int64 // 每页数据量
	// 以下字段通常在服务端内部使用，不需要客户端传递
	Offset *int64 // 数据偏移量
	Total  *int64 // 总数据量
}

type PostStatus uint8

const (
	Draft     PostStatus = iota // 草稿状态
	Published                   // 已发布状态
	Withdrawn                   // 撤回状态
	Deleted                     // 已删除状态
)

// String 方法用于将PostStatus转换为字符串
func (s PostStatus) String() string {
	switch s {
	case Draft:
		return "Draft"
	case Published:
		return "Published"
	case Withdrawn:
		return "Withdrawn"
	case Deleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}
