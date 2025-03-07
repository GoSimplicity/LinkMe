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
	ID           uint         `json:"id"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
	ReadCount    int64        `json:"read_count"`
	LikeCount    int64        `json:"like_count"`
	CollectCount int64        `json:"collect_count"`
	Uid          int64        `json:"uid"`
	Status       uint8        `json:"status"`
	PlateID      int64        `json:"plate_id"`
	Slug         string       `json:"slug"`
	CategoryID   int64        `json:"category_id"`
	Tags         string       `json:"tags"`
	CommentCount int64        `json:"comment_count"`
	IsSubmit     bool         `json:"is_submit"`
	Total        int64        `json:"total"`
}

type Interactive struct {
	BizID        uint  `json:"biz_id"`
	ReadCount    int64 `json:"read_count"`
	LikeCount    int64 `json:"like_count"`
	CollectCount int64 `json:"collect_count"`
	Liked        bool  `json:"liked"`
	Collected    bool  `json:"collected"`
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

func (p *Post) Abstract() string {
	// 将Content转换为一个rune切片
	str := []rune(p.Content)
	if len(str) > 128 {
		// 只保留前128个字符作为摘要
		str = str[:128]
	}
	return string(str)
}
