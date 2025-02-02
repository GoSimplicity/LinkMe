package ioc

import (
	"github.com/GoSimplicity/LinkMe/internal/job/interfaces"
	"github.com/GoSimplicity/LinkMe/internal/service"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

/*
 * Copyright 2024 Bamboo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * File: asynq.go
 * Description:
 */

func InitAsynqClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
	})
}

func InitAsynqServer() *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
		},
		asynq.Config{
			Concurrency: 10, // 设置并发数
		},
	)
}

func InitScheduler() *asynq.Scheduler {
	return asynq.NewScheduler(asynq.RedisClientOpt{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
	}, nil)
}

func InitRankingService(svc service.RankingService) interfaces.RankingService {
	return svc
}
