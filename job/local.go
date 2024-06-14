package job

import (
	"LinkMe/internal/domain"
	"context"
	"fmt"
)

// Executor 定义了任务执行器接口
type Executor interface {
	Name() string
	Exec(ctx context.Context, dj domain.Job) error
	RegisterFunc(name string, fn func(ctx context.Context, dj domain.Job) error)
}

// LocalFuncExecutor 本地方法执行器
type LocalFuncExecutor struct {
	localFunc map[string]func(ctx context.Context, dj domain.Job) error
}

// NewLocalFuncExecutor 创建并初始化 LocalFuncExecutor 实例
func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{localFunc: map[string]func(
		ctx context.Context, dj domain.Job) error{},
	}
}

// Name 返回执行器名称
func (l *LocalFuncExecutor) Name() string {
	return "local"
}

// RegisterFunc 注册本地执行函数
func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, dj domain.Job) error) {
	l.localFunc[name] = fn
}

// Exec 执行注册的本地方法
func (l *LocalFuncExecutor) Exec(ctx context.Context, dj domain.Job) error {
	fn, ok := l.localFunc[dj.Name]
	if !ok {
		return fmt.Errorf("local function not registered: %s", dj.Name)
	}
	return fn(ctx, dj)
}
