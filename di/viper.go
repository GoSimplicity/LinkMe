package di

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitViper 初始化配置管理
func InitViper() error {
	// 获取项目根目录
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取项目根目录失败: %v", err)
	}

	// 设置配置文件路径
	configFile := pflag.String("config", filepath.Join(dir, "config", "config.yaml"), "配置文件路径")
	pflag.Parse()

	// 设置配置文件
	viper.SetConfigFile(*configFile)

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("配置文件未找到: %s", *configFile)
		}
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 设置默认值
	viper.SetDefault("server.addr", ":8080")

	zap.L().Info("配置加载成功", zap.String("config_file", *configFile))
	return nil
}
