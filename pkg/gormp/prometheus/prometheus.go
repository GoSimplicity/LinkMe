package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type PrometheusCallbacks struct {
	operationMetrics *prometheus.SummaryVec // Prometheus 的 SummaryVec 用于跟踪操作时间的指标
}

func NewPrometheusCallbacks(opts prometheus.SummaryOpts) *PrometheusCallbacks {
	operationMetrics := prometheus.NewSummaryVec(opts, []string{"operation", "table"})
	prometheus.MustRegister(operationMetrics)
	return &PrometheusCallbacks{
		operationMetrics: operationMetrics,
	}
}

// Name 返回插件的名称
func (p *PrometheusCallbacks) Name() string {
	return "prometheus"
}

// Initialize 在 GORM 中注册 Prometheus 回调
func (p *PrometheusCallbacks) Initialize(db *gorm.DB) error {
	// 为各个操作类型注册 Prometheus 回调函数

	if err := db.Callback().Create().Before("gorm:create").Register("prometheus_create_before", p.beforeCallback()); err != nil {
		return err
	}
	if err := db.Callback().Create().After("gorm:create").Register("prometheus_create_after", p.afterCallback("CREATE")); err != nil {
		return err
	}

	if err := db.Callback().Query().Before("gorm:query").Register("prometheus_query_before", p.beforeCallback()); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:query").Register("prometheus_query_after", p.afterCallback("QUERY")); err != nil {
		return err
	}

	if err := db.Callback().Update().Before("gorm:update").Register("prometheus_update_before", p.beforeCallback()); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Register("prometheus_update_after", p.afterCallback("UPDATE")); err != nil {
		return err
	}

	if err := db.Callback().Delete().Before("gorm:delete").Register("prometheus_delete_before", p.beforeCallback()); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:delete").Register("prometheus_delete_after", p.afterCallback("DELETE")); err != nil {
		return err
	}

	if err := db.Callback().Raw().Before("gorm:raw").Register("prometheus_raw_before", p.beforeCallback()); err != nil {
		return err
	}
	if err := db.Callback().Raw().After("gorm:raw").Register("prometheus_raw_after", p.afterCallback("RAW")); err != nil {
		return err
	}

	if err := db.Callback().Row().Before("gorm:row").Register("prometheus_row_before", p.beforeCallback()); err != nil {
		return err
	}
	if err := db.Callback().Row().After("gorm:row").Register("prometheus_row_after", p.afterCallback("ROW")); err != nil {
		return err
	}

	return nil
}

// beforeCallback 是数据库操作前的回调，用于记录开始时间
func (p *PrometheusCallbacks) beforeCallback() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		start := time.Now()         // 记录当前时间为操作开始时间
		db.Set("start_time", start) // 将开始时间存储在 GORM 的数据库对象中
	}
}

// afterCallback 是数据库操作后的回调，用于记录操作时间并更新 Prometheus 指标
func (p *PrometheusCallbacks) afterCallback(operation string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		// 从数据库对象中获取开始时间
		startTime, ok := db.Get("start_time")
		if ok {
			start := startTime.(time.Time)               // 将获取的开始时间转换为 time.Time 类型
			duration := time.Since(start).Milliseconds() // 计算操作的持续时间（以毫秒为单位）
			// 将操作类型和表名作为标签，将持续时间记录到 Prometheus 指标中
			p.operationMetrics.WithLabelValues(operation, db.Statement.Table).Observe(float64(duration))
		}
	}
}
