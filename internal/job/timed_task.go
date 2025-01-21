package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type TimedTask struct {
	l *zap.Logger
}

type TimedPayload struct {
	TaskName    string    `json:"task_name"`
	LastRunTime time.Time `json:"last_run_time"`
}

func NewTimedTask(l *zap.Logger) *TimedTask {
	return &TimedTask{
		l: l,
	}
}

func (t *TimedTask) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload TimedPayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("解析任务载荷失败: %v: %w", err, asynq.SkipRetry)
	}

	t.l.Info("开始处理定时任务",
		zap.String("task_name", payload.TaskName),
		zap.Time("last_run_time", payload.LastRunTime))

	taskCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 定义任务处理映射
	taskHandlers := map[string]func(context.Context) error{
		GetRankingTask: nil,
	}

	// 获取对应的处理函数
	handler, exists := taskHandlers[payload.TaskName]
	if !exists {
		return fmt.Errorf("未知的任务类型: %s", payload.TaskName)
	}

	// 执行任务处理
	if err := handler(taskCtx); err != nil {
		t.l.Error("任务执行失败",
			zap.String("task_name", payload.TaskName),
			zap.Error(err))
		return fmt.Errorf("%s: %w", payload.TaskName, err)
	}

	t.l.Info("成功完成任务", zap.String("task_name", payload.TaskName))
	return nil
}
