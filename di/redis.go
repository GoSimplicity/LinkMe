package di

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	// 初始化 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),            // 地址
		Password:     viper.GetString("redis.password"),        // 密码
		DB:           viper.GetInt("redis.db"),                 // 数据库
		PoolSize:     viper.GetInt("redis.pool_size"),          // 连接池大小
		MinIdleConns: viper.GetInt("redis.min_idle_conns"),     // 最小空闲连接数
		MaxRetries:   viper.GetInt("redis.max_retries"),        // 最大重试次数
		ReadTimeout:  viper.GetDuration("redis.read_timeout"),  // 读取超时时间
		WriteTimeout: viper.GetDuration("redis.write_timeout"), // 写入超时时间
	})

	return client
}
