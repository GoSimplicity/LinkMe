package job

import "github.com/hibiken/asynq"

type Routes struct {
	RefreshCache *RefreshCacheTask
}

func NewRoutes(refreshCache *RefreshCacheTask) *Routes {
	return &Routes{
		RefreshCache: refreshCache,
	}
}

func (r *Routes) RegisterHandlers() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(DeferRefreshPostCache, r.RefreshCache.ProcessTask)
	// 注册其他任务处理器

	return mux
}
