package config

import (
	"errors"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// InitViper 初始化viper配置
func InitViper() {
	setDefaults()

	configFile := pflag.String("config", "config/config.yaml", "配置文件路径")
	pflag.Parse()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(*configFile)
	viper.SetEnvPrefix("LINKME")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	var notFound viper.ConfigFileNotFoundError
	if err != nil && !errors.As(err, &notFound) {
		panic(err)
	}
}

func setDefaults() {
	viper.SetDefault("server.addr", ":9999")
	viper.SetDefault("metrics.addr", ":9091")
	viper.SetDefault("log.dir", "logs")
	viper.SetDefault("jwt.auth_expire", 30)
	viper.SetDefault("jwt.refresh_expire", 150)
	viper.SetDefault("sms.provider", "mock")
	viper.SetDefault("email.provider", "mock")
	viper.SetDefault("ark_api.provider", "mock")
	viper.SetDefault("es.bootstrap_indexes", true)
	viper.SetDefault("cors.allow_all", false)
	viper.SetDefault("cors.allow_origins", []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	})
}
