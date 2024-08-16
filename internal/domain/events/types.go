package events

import "context"

type Consumer interface {
	Start(ctx context.Context) error
}
