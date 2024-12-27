package general

import (
	"context"
)

func WithAsyncCancel(ctx context.Context, cancel context.CancelFunc, fn func() error) func() {
	return func() {
		go func() {
			// 监听 context 取消信号
			done := make(chan struct{})
			defer close(done)

			go func() {
				select {
				case <-ctx.Done():
					cancel()
				case <-done:
					return
				}
			}()

			// 确保 goroutine 中的 panic 不会导致程序崩溃
			defer func() {
				if r := recover(); r != nil {
					cancel() // 发生 panic 时取消操作
				}
			}()

			// 执行目标函数
			if err := fn(); err != nil {
				cancel() // 发生错误时取消操作
			}
		}()
	}
}
