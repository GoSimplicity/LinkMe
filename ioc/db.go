package ioc

import (
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	prometheus3 "github.com/GoSimplicity/LinkMe/pkg/gormp/prometheus"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
	"log"
)

type config struct {
	DSN string `yaml:"dsn"`
}

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	var c config

	if err := viper.UnmarshalKey("db", &c); err != nil {
		panic(fmt.Errorf("init failed：%v", err))
	}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}
	// 初始化表

	if err = dao.InitTables(db); err != nil {
		panic(err)
	}

	// 注册 Prometheus 插件
	if err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "linkme", // Prometheus中标识数据库的名称
		RefreshInterval: 5,        // 监控数据刷新间隔，单位为秒
	})); err != nil {
		log.Println("register prometheus plugin failed")
	}

	// 创建并注册自定义的 PrometheusCallbacks 插件，用于监控gorm操作执行时间
	prometheusPlugin := prometheus3.NewPrometheusCallbacks(prometheus2.SummaryOpts{
		Namespace: "linkme",                                          // 命名空间
		Subsystem: "gorm",                                            // 子系统
		Name:      "operation_duration_seconds",                      // 指标名称
		Help:      "Duration of GORM database operations in seconds", // 指标帮助信息
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})

	if err = db.Use(prometheusPlugin); err != nil {
		log.Println("register custom prometheus callbacks plugin failed:", err)
	}

	return db
}
