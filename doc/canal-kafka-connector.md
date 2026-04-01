# Canal Kafka Connector 调试记录

以下命令用于本地调试 Debezium MySQL Connector，不参与应用构建或运行：

```bash
# 拉取镜像
docker pull quay.io/debezium/connect

# 创建容器
docker run -it --rm --name linkme-connect -p 8083:8083 \
  -e GROUP_ID=1 \
  -e CONFIG_STORAGE_TOPIC=my_connect_configs \
  -e OFFSET_STORAGE_TOPIC=my_connect_offsets \
  -e STATUS_STORAGE_TOPIC=my_connect_statuses \
  -e BOOTSTRAP_SERVERS=<kafka-host>:9092 \
  --link linkme-kafka:linkme-kafka --link linkme-mysql:linkme-mysql \
  --network linkme_default \
  quay.io/debezium/connect

# 创建 connector
curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d '
{
  "name": "linkme-connector",
  "config": {
    "connector.class": "io.debezium.connector.mysql.MySqlConnector",
    "tasks.max": "1",
    "database.hostname": "linkme-mysql",
    "database.port": "3306",
    "database.user": "root",
    "database.password": "<mysql-password>",
    "database.server.id": "184054",
    "database.server.name": "linkme",
    "database.include.list": "linkme",
    "schema.history.internal.kafka.bootstrap.servers": "<kafka-host>:9092",
    "schema.history.internal.kafka.topic": "schema-changes.linkme",
    "topic.prefix": "oracle"
  }
}'
```
