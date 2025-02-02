package interfaces

import (
	"context"
	"time"
)

type Post struct {
	ID        uint
	Title     string
	Content   string
	UpdatedAt time.Time
}

type RankingService interface {
	TopN(ctx context.Context) error
}
