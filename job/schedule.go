package job

import (
	"LinkMe/internal/service"
	"context"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"time"
)

type Scheduler struct {
	dbTimeout time.Duration
	svc       service.CronJobService
	executors map[string]Executor // 执行器映射
	l         *zap.Logger
	limiter   *semaphore.Weighted // 信号量，限制并发任务数量
}

func NewScheduler(svc service.CronJobService, l *zap.Logger) *Scheduler {
	return &Scheduler{
		svc:       svc,
		dbTimeout: time.Second,
		limiter:   semaphore.NewWeighted(100),
		l:         l,
		executors: map[string]Executor{},
	}
}

// RegisterExecutor 注册任务执行器
func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.executors[exec.Name()] = exec
}

// Schedule 调度任务执行
func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		// 检查上下文是否取消
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// 获取信号量，限制并发任务数量
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		// 设置数据库操作超时上下文
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			// 如果出错，等待10秒后重试
			time.Sleep(10 * time.Second)
			continue
		}
		// 查找任务的执行器
		exec, ok := s.executors[j.Executor]
		if !ok {
			s.l.Error("Executor not found")
			continue
		}
		// 启动协程执行任务
		go func() {
			defer func() {
				s.limiter.Release(1) // 释放信号量
				j.CancelFunc()       // 调用任务的取消函数
			}()
			// 执行任务
			er := exec.Exec(ctx, j)
			if er != nil {
				s.l.Error("Failed to execute job", zap.Error(er))
				return
			}
			// 重置任务的下次执行时间
			er = s.svc.ResetNextTime(ctx, j)
			if er != nil {
				s.l.Error("Failed to reset next execution time", zap.Error(er))
			}
		}()
	}
}
