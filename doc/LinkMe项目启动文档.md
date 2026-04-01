# LinkMe 启动与调试文档

本文档只描述当前仓库代码对应的启动方式，不再保留历史截图和过时的部署步骤。

## 1. 当前启动形态

当前服务由 `cmd/main.go` 单入口启动，进程内会同时拉起：

- Gin Web API
- Prometheus 指标服务
- Asynq Worker
- Asynq Scheduler
- Kafka Consumers

因此，除应用本身外，最少还需要准备：

- MySQL
- Redis
- Kafka
- Elasticsearch

短信、邮件、AI 审核在默认配置下可以走 `mock`，本地联调时不强制要求真实提供商。

## 2. 前置要求

- Go 1.22+
- Docker / Docker Compose

初始化依赖：

```bash
go mod tidy
cp config/config.example.yaml config/config.yaml
```

## 3. 推荐方式：本地进程 + Docker Compose 依赖环境

这是当前最适合开发调试的启动方式。

### 3.1 启动依赖环境

```bash
make init
docker-compose -f docker-compose-env.yaml up -d
```

依赖环境暴露到宿主机的关键端口如下：

| 组件 | 端口 |
| --- | --- |
| MySQL | `43306` |
| Redis | `46379` |
| Kafka | `9094` |
| Elasticsearch | `19200` |
| Kibana | `5601` |
| Prometheus | `9090` |
| Grafana | `3001` |
| Asynqmon | `8980` |

### 3.2 初始化数据库

`docker-compose-env.yaml` 当前不会自动导入 SQL，因此首次启动需要手动导入：

```bash
docker-compose -f docker-compose-env.yaml exec -T mysql sh -c 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" linkme' < deploy/init/linkme.sql
```

### 3.3 修改本地配置

如果你直接使用 `docker-compose-env.yaml`，需要把 `config/config.yaml` 中的地址改成宿主机映射端口，而不是 `config.example.yaml` 里的默认端口。

建议至少确认下面这些配置：

```yaml
server:
  addr: ":9999"

metrics:
  addr: ":9091"

db:
  dsn: "root:v6SxhWHyZC7S@tcp(localhost:43306)/linkme?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "localhost:46379"
  password: "v6SxhWHyZC7S"

kafka:
  addr:
    - "localhost:9094"

es:
  addr: "http://localhost:19200"

sms:
  provider: "mock"

email:
  provider: "mock"

ark_api:
  provider: "mock"
```

说明：

- `kafka.addr` 在 YAML 中建议写数组
- 如果使用环境变量覆盖，也可以写单值字符串，例如 `LINKME_KAFKA_ADDR=localhost:9094`
- `config/config.yaml` 仅用于本地，不应提交到仓库

### 3.4 启动应用

```bash
go run ./cmd/main.go
```

启动后可以检查：

- 应用健康检查：`http://localhost:9999/healthz`
- 就绪检查：`http://localhost:9999/readyz`
- 指标端点：`http://localhost:9091/metrics`

## 4. Docker Compose 启动应用

如果你希望把应用和 Nginx 网关也一起拉起来，可以继续执行：

```bash
make build
docker-compose up -d
```

说明：

- 根目录 `docker-compose.yaml` 里的 `linkme` 服务通过环境变量注入配置，不依赖本地 `config/config.yaml`
- 该启动方式当前更偏向开发联调，因为 `Dockerfile` 使用 `modd` 作为容器启动命令

访问入口：

- 网关：`http://localhost:8888`
- 应用：`http://localhost:9999`

## 5. 常用命令

```bash
make fmt
make test-unit
make build-app
make env-up
make build
make up
make down
make logs
```

## 6. K8s 部署说明

仓库里仍保留了 `deploy/yaml/` 下的 Kubernetes 清单：

- `mysql.yaml`
- `redis.yaml`
- `kafka.yaml`
- `es.yaml`
- `mongo.yaml`
- `prometheus.yaml`
- `grafana.yaml`
- `kibana.yaml`
- `logstash.yaml`
- `canal.yaml`

这部分清单可以作为部署起点，但当前没有再维护成“一键照抄即可运行”的文档流程。实际使用时需要根据你的集群存储、NodePort、域名、镜像源和安全策略自行调整。

## 7. 常见问题

### 7.1 依赖都启动了，但应用仍连不上

优先检查是不是直接拿了 `config/config.example.yaml` 默认值。

`docker-compose-env.yaml` 的宿主机端口不是：

- MySQL `3306`
- Redis `6379`
- Kafka `9092`
- Elasticsearch `9200`

而是：

- MySQL `43306`
- Redis `46379`
- Kafka `9094`
- Elasticsearch `19200`

### 7.2 登录、发短信、发邮件、AI 审核失败

先确认 provider 是否仍然是 `mock`：

- `sms.provider`
- `email.provider`
- `ark_api.provider`

如果切换到了真实提供商，必须同时配置对应密钥。

### 7.3 搜索或事件链路不工作

优先检查：

- Kafka 是否可达
- Elasticsearch 是否可达
- 应用日志中 Kafka Consumer 是否启动成功
- SQL 是否已经导入

Canal Connector 的本地调试命令已单独整理到：

- [doc/canal-kafka-connector.md](./canal-kafka-connector.md)
