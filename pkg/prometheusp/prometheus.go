package prometheusp

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MetricsPlugin struct {
	Namespace           string
	Subsystem           string
	InstanceID          string
	responseTimeVec     *prometheus.SummaryVec // 用于记录HTTP请求响应时间的SummaryVec指标
	activeRequestsGauge prometheus.Gauge       // 用于记录当前活跃HTTP请求数的Gauge指标
}

func NewMetricsPlugin(namespace, subsystem, instanceID string) *MetricsPlugin {
	return &MetricsPlugin{
		Namespace:  namespace,
		Subsystem:  subsystem,
		InstanceID: instanceID,
	}
}

// RegisterMetrics 注册所有 Prometheus 指标
// 包括响应时间SummaryVec和活跃请求数Gauge
func (m *MetricsPlugin) RegisterMetrics() {
	// 注册响应时间指标
	m.responseTimeVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      "gin_response_time",          // 指标名称
		Help:      "HTTP request response time", // 指标的帮助信息
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID, // 常量标签，标识具体的实例
		},
		Objectives: map[float64]float64{ // 定义不同百分位的目标值
			0.5:   0.01,   // 50%分位
			0.75:  0.01,   // 75%分位
			0.9:   0.01,   // 90%分位
			0.99:  0.001,  // 99%分位
			0.999: 0.0001, // 99.9%分位
		},
	}, []string{"method", "pattern", "status"}) // 标签：HTTP方法、路径模式、状态码
	prometheus.MustRegister(m.responseTimeVec) // 注册指标到Prometheus

	// 注册活跃请求数指标
	m.activeRequestsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      "active_requests",                // 指标名称
		Help:      "Number of active HTTP requests", // 指标的帮助信息
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID, // 常量标签，标识具体的实例
		},
	})
	prometheus.MustRegister(m.activeRequestsGauge) // 注册指标到Prometheus
}

// TrackResponseTimeMiddleware 返回一个 Gin 中间件，用于跟踪 HTTP 请求的响应时间
func (m *MetricsPlugin) TrackResponseTimeMiddleware() gin.HandlerFunc {
	// 在请求处理完成后记录响应时间
	return func(ctx *gin.Context) {
		start := time.Now() // 记录请求开始时间
		defer func() {
			duration := time.Since(start).Milliseconds() // 计算请求处理的持续时间
			method := ctx.Request.Method                 // 获取请求方法
			pattern := ctx.FullPath()                    // 获取请求的路径模式
			status := ctx.Writer.Status()                // 获取响应的状态码
			// 将响应时间记录到SummaryVec中
			m.responseTimeVec.WithLabelValues(method, pattern, strconv.Itoa(status)).
				Observe(float64(duration))
		}()
		ctx.Next() // 继续处理请求
	}
}

// TrackActiveRequestsMiddleware 返回一个 Gin 中间件，用于跟踪活跃 HTTP 请求数
func (m *MetricsPlugin) TrackActiveRequestsMiddleware() gin.HandlerFunc {
	// 在请求进入时增加计数，在请求完成时减少计数
	return func(ctx *gin.Context) {
		m.activeRequestsGauge.Inc()       // 活跃请求数加1
		defer m.activeRequestsGauge.Dec() // 请求完成后活跃请求数减1
		ctx.Next()                        // 继续处理请求
	}
}

// Apply 将所有的 Prometheus 中间件应用到 Gin 路由
func (m *MetricsPlugin) Apply(router *gin.Engine) {
	router.Use(m.TrackResponseTimeMiddleware())   // 添加响应时间跟踪中间件
	router.Use(m.TrackActiveRequestsMiddleware()) // 添加活跃请求数跟踪中间件
}

// AddCustomGauge 添加一个自定义的 Gauge 类型的指标
func (m *MetricsPlugin) AddCustomGauge(name, help string, constLabels map[string]string) prometheus.Gauge {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   m.Namespace,
		Subsystem:   m.Subsystem,
		Name:        name,        // 指标名称
		Help:        help,        // 指标的帮助信息
		ConstLabels: constLabels, // 常量标签
	})
	prometheus.MustRegister(gauge)
	return gauge
}

// AddCustomSummaryVec 添加一个自定义的 SummaryVec 类型的指标
func (m *MetricsPlugin) AddCustomSummaryVec(name, help string, labels []string, objectives map[float64]float64, constLabels map[string]string) *prometheus.SummaryVec {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   m.Namespace,
		Subsystem:   m.Subsystem,
		Name:        name,        // 指标名称
		Help:        help,        // 指标的帮助信息
		Objectives:  objectives,  // 百分位目标值
		ConstLabels: constLabels, // 常量标签
	}, labels) // 标签数组
	prometheus.MustRegister(summaryVec)
	return summaryVec
}
