package domain

import "sync/atomic"

const (
	Draft     = "Draft"     // 草稿状态
	Published = "Published" // 发布状态
	Withdrawn = "Withdrawn" // 撤回状态
	Deleted   = "Deleted"   // 删除状态
)

type Author struct {
	Id   int64
	Name string
}

type Post struct {
	ID           int64
	Title        string
	Content      string
	CreateTime   int64
	UpdatedTime  int64
	DeletedTime  int64
	Deleted      bool
	Author       Author
	Status       string
	Visibility   string
	Slug         string
	CategoryID   int64
	Tags         string
	CommentCount int64
}

type Interactive struct {
	BizID        int64
	ReadCount    int64
	LikeCount    int64
	CollectCount int64
	Liked        bool
	Collected    bool
}

func (i *Interactive) IncrementReadCount() {
	atomic.AddInt64(&i.ReadCount, 1)
}

func (i *Interactive) IncrementLikeCount() {
	atomic.AddInt64(&i.LikeCount, 1)
}

func (i *Interactive) IncrementCollectCount() {
	atomic.AddInt64(&i.CollectCount, 1)
}

func (p Post) Abstract() string {
	// 将Content转换为一个rune切片
	str := []rune(p.Content)
	if len(str) > 128 {
		// 只保留前128个字符作为摘要
		str = str[:128]
	}
	return string(str)
}
