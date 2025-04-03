/*
 * MIT License
 *
 * Copyright (c) 2024 Bamboo
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package di

import (
	"time"

	"github.com/GoSimplicity/LinkMe/internal/pkg/infra/database/migrations"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	DSN             string `yaml:"dsn"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"` // 单位：分钟
	LogLevel        string `yaml:"log_level"`
}

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	var config DBConfig

	if err := viper.UnmarshalKey("db", &config); err != nil {
		zap.L().Fatal("数据库配置解析失败", zap.Error(err))
	}

	if config.DSN == "" {
		zap.L().Fatal("数据库连接字符串(DSN)未配置")
	}

	// 设置默认值
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 10
	}
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 100
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = 60 // 默认60分钟
	}

	// 配置日志级别
	logLevel := getGormLogLevel(config.LogLevel)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		zap.L().Fatal("数据库连接失败", zap.Error(err), zap.String("dsn", config.DSN))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Fatal("获取数据库连接池失败", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Minute)

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		zap.L().Fatal("数据库连接测试失败", zap.Error(err))
	}

	zap.L().Info("数据库连接成功",
		zap.Int("最大空闲连接数", config.MaxIdleConns),
		zap.Int("最大打开连接数", config.MaxOpenConns),
		zap.Int("连接最大生命周期(分钟)", config.ConnMaxLifetime))

	// 初始化表
	if err = initTables(db); err != nil {
		zap.L().Fatal("数据库表初始化失败", zap.Error(err))
	}

	return db
}

// initTables 初始化数据库表
func initTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&migrations.User{},
	)
}

// getGormLogLevel 根据配置字符串返回对应的GORM日志级别
func getGormLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info // 默认为Info级别
	}
}
