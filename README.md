# LinkMe

LinkMe 是一个基于 Go 的论坛后端服务仓库。当前代码库仍然是单体部署形态，但运行时已经包含 HTTP API、异步任务、定时任务、Kafka 消费者和监控端点。

## 当前项目状态

- 单入口服务：`cmd/main.go`
- HTTP 服务：Gin
- 认证与权限：JWT + Casbin
- 数据存储：MySQL、Redis、Elasticsearch
- 异步链路：Kafka、Asynq、定时热榜任务
- 日志与监控：Zap JSON 日志、Prometheus、Grafana、ELK
- 提供商模式：短信、邮件、AI 审核默认可走 `mock`

相关仓库：

- 微服务版本：https://github.com/GoSimplicity/LinkMe-microservices
- 前端项目：https://github.com/GoSimplicity/LinkMe-web

## 运行链路

当前进程启动后会同时运行以下组件：

- Web API 服务，默认监听 `:9999`
- 健康检查接口：`/healthz`、`/readyz`
- Prometheus 指标服务，默认监听 `:9091`
- Asynq Worker
- Asynq Scheduler
- Kafka Consumers
- 热榜定时任务，当前实现为每小时刷新一次

## 核心能力

- 用户注册、登录、刷新令牌、短信登录、邮箱/短信验证码发送
- 用户资料查询与更新、用户列表、权限码查询
- 帖子草稿编辑、更新、发布、撤回、删除、详情、列表、按版块筛选
- 点赞、收藏、浏览历史、近期活动
- 评论创建、删除、分页加载、楼中楼回复
- 关注/取关、粉丝与关注列表、计数查询
- 版块管理、热榜查询、搜索
- 内容审核、角色、权限、菜单、API 资源管理
- 抽奖与秒杀活动接口

## 技术栈

- Go 1.22
- Gin
- Wire
- GORM
- Casbin
- Redis / Kafka / Elasticsearch / MySQL
- Asynq
- Prometheus / Grafana / ELK

说明：

- 仓库里仍保留了 MongoDB 初始化代码和 K8s YAML，但它不是当前默认启动链路的一部分。
- `Dockerfile` 当前是开发态 `modd` 运行镜像，不是最小化生产镜像。

## 快速开始

### 方式一：本地进程 + Docker Compose 依赖环境

1. 安装依赖并复制配置：

```bash
go mod tidy
cp config/config.example.yaml config/config.yaml
```

2. 启动依赖环境：

```bash
make init
docker-compose -f docker-compose-env.yaml up -d
```

3. 首次初始化数据库：

```bash
docker-compose -f docker-compose-env.yaml exec -T mysql sh -c 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" linkme' < deploy/init/linkme.sql
```

4. 按 `docker-compose-env.yaml` 的宿主机端口修改本地配置：

- MySQL：`localhost:43306`
- Redis：`localhost:46379`
- Kafka：`localhost:9094`
- Elasticsearch：`http://localhost:19200`

5. 启动服务：

```bash
go run ./cmd/main.go
```

### 方式二：Docker Compose 启动应用与网关

```bash
make init
docker-compose -f docker-compose-env.yaml up -d
make build
docker-compose up -d
```

默认入口：

- Nginx 网关：`http://localhost:8888`
- 应用端口：`http://localhost:9999`
- 指标端口：`http://localhost:9091/metrics`

## 常用命令

```bash
make fmt
make test-unit
make build-app
make env-up
make build
make up
make down
make deploy
make logs
```

## 配置说明

配置文件示例见 `config/config.example.yaml`。

- 默认配置文件路径：`config/config.yaml`
- 环境变量前缀：`LINKME_`
- 点号会被转换为下划线
- `kafka.addr` 同时兼容 YAML 数组和单个环境变量字符串

常用配置项：

- `server.addr`
- `metrics.addr`
- `db.dsn`
- `redis.addr`
- `redis.password`
- `kafka.addr`
- `es.addr`
- `sms.provider`
- `email.provider`
- `ark_api.provider`
- `cors.allow_all`
- `cors.allow_origins`

提供商建议：

- 本地联调默认用 `mock`
- 短信真实发送使用 `sms.provider=tencent`
- 邮件真实发送使用 `email.provider=qq`
- AI 审核真实调用使用 `ark_api.provider=ark`

## 目录结构

```text
.
├── cmd/main.go                 # 单进程启动入口
├── config/                     # Viper 配置与示例配置
├── deploy/                     # 初始化 SQL、Compose 依赖和 K8s YAML
├── doc/                        # 项目文档
├── internal/api/               # HTTP Handler 与请求参数
├── internal/constants/         # 常量定义
├── internal/domain/            # 领域模型与事件定义
├── internal/job/               # Asynq 任务与定时任务
├── internal/repository/        # 仓储实现（DB / Cache / ES）
├── internal/service/           # 业务服务
├── ioc/                        # 依赖注入与基础设施初始化
├── middleware/                 # 中间件
├── pkg/                        # 通用基础组件
└── utils/                      # 辅助工具与内容审核能力
```

## 文档索引

- [启动与调试文档](./doc/LinkMe项目启动文档.md)
- [功能模块总览](./doc/function_module.md)
- [Canal Kafka Connector 调试记录](./doc/canal-kafka-connector.md)
- [项目面试亮点（偏设计表达）](./doc/project_interview_highlights.md)

## 许可证

本项目使用 MIT 许可证，详情见 [LICENSE](./LICENSE)。
