package job

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	GetRankingTask = "get_ranking"
)

type TimedScheduler struct {
	scheduler *asynq.Scheduler
}

func NewTimedScheduler(scheduler *asynq.Scheduler) *TimedScheduler {
	return &TimedScheduler{
		scheduler: scheduler,
	}
}

func (s *TimedScheduler) RegisterTimedTasks() error {
	// K8s 客户端刷新任务 - 每5分钟
	if err := s.registerTask(
		GetRankingTask,
		"@every 5m",
	); err != nil {
		return err
	}

	return nil
}

func (s *TimedScheduler) registerTask(taskName, cronSpec string) error {
	payload := TimedPayload{
		TaskName:    taskName,
		LastRunTime: time.Now(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(DeferTimedTask, payloadBytes)
	_, err = s.scheduler.Register(cronSpec, task)
	return err
}

func (s *TimedScheduler) Run() error {
	return s.scheduler.Run()
}

func (s *TimedScheduler) Stop() {
	s.scheduler.Shutdown()
}
