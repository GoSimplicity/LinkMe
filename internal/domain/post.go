package domain

import (
	"database/sql"
	"sync/atomic"
	"time"
)

const (
	Draft     uint8 = iota // 0: 草稿状态
	Published              // 1: 发布状态
	Withdrawn              // 2: 撤回状态
	Deleted                // 3: 删除状态
)

type Post struct {
	ID           uint
	Title        string
	Content      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
	AuthorID     int64
	Status       uint8
	Visibility   string
	PlateID      int64
	Slug         string
	CategoryID   int64
	Tags         string
	CommentCount int64
}

type Interactive struct {
	BizID        uint
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
