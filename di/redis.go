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
