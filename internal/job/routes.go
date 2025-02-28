package job

import "github.com/hibiken/asynq"

type Routes struct {
	RefreshCache *RefreshCacheTask
	TimedTask    *TimedTask
}

func NewRoutes(refreshCache *RefreshCacheTask, timedTask *TimedTask) *Routes {
	return &Routes{
		RefreshCache: refreshCache,
		TimedTask:    timedTask,
	}
}

func (r *Routes) RegisterHandlers() *asynq.ServeMux {
	mux := asynq.NewServeMux()

	mux.HandleFunc(RefreshPostCache, r.RefreshCache.ProcessTask)
	mux.HandleFunc(DeferTimedTask, r.TimedTask.ProcessTask)

	return mux
}
