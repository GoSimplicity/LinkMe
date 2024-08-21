## clone项目到本地
```bash
git clone git@github.com:GoSimplicity/LinkMe.git
```
如无法使用git clone可自行下载项目zip包解压
## 进入项目根目录初始化依赖
```bash
go mod tidy
```
## 修改config/config.yaml
```yaml
db:
  dsn: "root:root@tcp(ip_addr:端口)/linkme?charset=utf8mb4&parseTime=True&loc=Local"
redis:
  addr: "ip_addr:端口"
log:
  filepath: "logs/linkme.json"
mongodb:
  addr: "mongodb://ip_addr:端口"
kafka:
  addr: "ip_addr:端口"
es:
  addr: "http://ip_addr:端口"
sms:
  tencent:
    secretId: ""
    secretKey: ""
    endPoint: ""
    smsID: ""
    sign: ""
    templateID: ""

```
请自行将ip_addr和端口根据实际情况进行替换
## 部署中间件
### 使用docker-compose部署中间件
```bash
cd deploy
```
修改docker-compose.yaml
```bash
version: "3"
services:
  kafka:
    image: "bitnami/kafka:3.6.0"
    container_name: linkme-kafka
    restart: always
    ports:
      - "9092:9092"
      - "9094:9094"
    environment:
      - KAFKA_CFG_NODE_ID=0
      - KAFKA_CREATE_TOPICS=linkme_binlog:3:1
      - KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=true
      - KAFKA_CFG_PROCESS_ROLES=controller,broker
      - KAFKA_CFG_LISTENERS=PLAINTEXT://0.0.0.0:9092,CONTROLLER://:9093,EXTERNAL://0.0.0.0:9094
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092,EXTERNAL://localhost:9094 # 注意此处，如果不是在本机docker启动的，需要将这里的localhost改为宿主机ip
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_MESSAGE_MAX_BYTES=20971520
      - KAFKA_CFG_LOG_DIRS=/bitnami/kafka/data # 指定 Kafka 数据目录
    volumes:
      - /data/kafka:/bitnami/kafka/data # 注意此处/data/kafka修改为你实际的目录
  db:
    image: mysql:8.0
    container_name: linkme-mysql
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: linkme
    volumes:
      - /data/mysql:/var/lib/mysql # 注意此处/data/mysql修改为你实际的目录
      - ./init:/docker-entrypoint-initdb.d/ # 如果运行时此处报错，请改为绝对路径
  redis:
    image: redis:latest
    container_name: linkme-redis
    restart: always
    ports:
      - "6379:6379"
    command:
      - "redis-server"
      - "--bind"
      - "0.0.0.0"
      - "--protected-mode"
      - "no"
      - "--port"
      - "6379"
  mongo:
    image: mongo:latest
    container_name: linkme-mongo
    restart: always
    ports:
      - "27017:27017"
    environment:
      - MONGO_INITDB_DATABASE=linkme # 设置默认数据库名
    volumes:
      - /data/mongo:/data/db  # 注意此处/data/mongo修改为你实际的目录
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.12.2
    container_name: elasticsearch
    environment:
      - discovery.type=single-node
      - ES_JAVA_OPTS=-Xms1g -Xmx1g
    volumes:
      - /data/es/data:/usr/share/elasticsearch/data # 注意此处/data/es/data修改为你实际的目录
      - ./elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml # 如果运行时此处报错，请改为绝对路径
    ports:
      - "9200:9200"
      - "9300:9300"
  canal:
    image: canal/canal-server
    container_name: linkme-canal
    environment:
      - CANAL_IP=canal-server
      - CANAL_PORT=11111
      - CANAL_DESTINATIONS=linkme
    depends_on:
      - db
      - kafka
    ports:
      - "11111:11111"
    volumes:
      - ./sync:/home/admin/canal-server/conf/ # 如果运行时此处报错，请改为绝对路径
      - /data/canal/logs:/home/admin/canal-server/logs # 注意此处修改为你实际的目录
      - /data/canal/destinations:/home/admin/canal-server/destinations # 注意此处修改为你实际的目录
      - ./canal.properties:/home/admin/canal-server/conf/canal.properties # 如果运行时此处报错，请改为绝对路径
  prometheus:
    image: bitnami/prometheus:latest
    container_name: linkme-prometheus
    volumes:
      - ./prometheus.yml:/opt/bitnami/prometheus/conf/prometheus.yml # 如果运行时此处报错，请改为绝对路径
    ports:
      - "9090:9090"
```
### 使用k8s部署中间件
注意：所有操作均在deploy目录下进行
```bash
cd deploy/
```
```bash
kubectl create ns linkme # 创建命名空间
```
创建持久化目录
```yaml
mkdir /data
```
#### 部署mysql
```bash
mkdir /data/mysql && cp -rvf init /data/mysql/
kubectl apply -f yaml/mysql.yaml
```
#### 部署redis
```bash
kubectl apply -f yaml/redis.yaml
```
#### 部署mongodb
```bash
mkdir /data/mongo
kubectl apply -f yaml/mongo.yaml
```
#### 部署kafka
```bash
mkdir /data/kafka
# 修改yaml/kafka.yaml
- name: KAFKA_CFG_ADVERTISED_LISTENERS
  # 需将此处的ip_addr改为你的宿主机ip 端口不需要进行修改
  value: "PLAINTEXT://ip_addr:30880,EXTERNAL://ip_addr:9094"
kubectl apply -f yaml/kafka.yaml
```
![image.png](https://cdn.nlark.com/yuque/0/2024/png/40474956/1723806353368-90039049-24fb-414d-8a1c-7d5a83a900e8.png#averageHue=%23313030&clientId=u19722460-4426-4&from=paste&height=502&id=u664cfefd&originHeight=502&originWidth=836&originalType=binary&ratio=1&rotation=0&showTitle=false&size=101102&status=done&style=none&taskId=uda8cb3fc-5437-41d0-840b-e312682db90&title=&width=836)
#### 部署canal
```bash
mkdir -p /data/canal/conf
```
修改配置文件deploy/canal/canal.properties
![image.png](https://cdn.nlark.com/yuque/0/2024/png/40474956/1723806502298-63b20ed2-038b-4f59-9f52-bc616bafc44a.png#averageHue=%232f2e2d&clientId=u19722460-4426-4&from=paste&height=415&id=u7143b21b&originHeight=415&originWidth=877&originalType=binary&ratio=1&rotation=0&showTitle=false&size=60862&status=done&style=none&taskId=u5e4ca741-c95a-43c1-b72a-9ffd4512b4b&title=&width=877)
修改配置文件deploy/canal/sync/instance.properties
![image.png](https://cdn.nlark.com/yuque/0/2024/png/40474956/1723806551905-762560e9-a6a7-48b1-933e-74700ac5e166.png#averageHue=%232f2e2d&clientId=u19722460-4426-4&from=paste&height=540&id=u45837c7f&originHeight=540&originWidth=904&originalType=binary&ratio=1&rotation=0&showTitle=false&size=78590&status=done&style=none&taskId=u4dee3c89-99b2-4eb7-bc2e-04eebbf2451&title=&width=904)
然后需要将修改好的canal.properties和sync目录放到/data/canal/conf
```bash
cp canal.properties /data/canal/conf && cp -rvf sync /data/canal/conf
```
部署canal
```bash
kubectl apply -f yaml/canal.yaml
```
#### 部署es
```bash
kubectl apply -f yaml/es.yaml
```
#### 部署prometheus
修改配置文件
![image.png](https://cdn.nlark.com/yuque/0/2024/png/40474956/1723807279325-c6bafed9-5198-4e7b-b1e3-fca257670ecc.png#averageHue=%232e2c2b&clientId=u19722460-4426-4&from=paste&height=588&id=ua9682209&originHeight=588&originWidth=576&originalType=binary&ratio=1&rotation=0&showTitle=false&size=51227&status=done&style=none&taskId=u28bb0e6e-7130-4295-a01e-6ac31519ca4&title=&width=576)
```bash
mkdir /data/prometheus && cp prometheus.yml /data/prometheus
kubectl apply -f yaml/prometheus.yaml
```
#### 部署ELK
##### 部署logstash
修改配置文件
![image.png](https://cdn.nlark.com/yuque/0/2024/png/40474956/1723807808662-299b9b3e-9e89-4497-b833-e53fe66feb97.png#averageHue=%23516548&clientId=u19722460-4426-4&from=paste&height=625&id=u3885617c&originHeight=625&originWidth=896&originalType=binary&ratio=1&rotation=0&showTitle=false&size=110629&status=done&style=none&taskId=u63db72b3-8a98-4f03-918e-b45d63857df&title=&width=896)
将此处修改为实际IP地址
![image.png](https://cdn.nlark.com/yuque/0/2024/png/40474956/1723807854632-f537c415-04f5-488d-9dda-e12aba7d40a2.png#averageHue=%239e8539&clientId=u19722460-4426-4&from=paste&height=529&id=ua53f2aaf&originHeight=529&originWidth=931&originalType=binary&ratio=1&rotation=0&showTitle=false&size=82356&status=done&style=none&taskId=uba89c1eb-7cee-44cd-ba56-f3545ea887d&title=&width=931)
将此处修改为实际IP地址
```bash
mkdir -p /data/logstash/conf && cp logstash.conf /data/logstash/conf/ && cp logstash.yml /data/logstash/conf/
kubectl apply -f logstash.yaml
```
##### 部署kibana
修改配置文件
![image.png](https://cdn.nlark.com/yuque/0/2024/png/40474956/1723807992431-d9b5e099-4652-452d-bf6f-00b82e9383df.png#averageHue=%23384337&clientId=u19722460-4426-4&from=paste&height=648&id=u5abd58e7&originHeight=648&originWidth=911&originalType=binary&ratio=1&rotation=0&showTitle=false&size=131839&status=done&style=none&taskId=u5a22c49c-bfe8-4b19-bff0-1cc8c41e954&title=&width=911)
修改此处为实际IP
```bash
kubectl apply -f yaml/kibana.yaml
```
### 可能会遇到的问题
#### 配置文件确认都写对了而kakfa接收不到消息推送
```bash
# 首先检查log日志输出
# 检查kafka情况
# 最容易踩的坑是分区数量不一致，使用下面命令查看主题详情，如果分区数为1需要添加分区
```
##### 创建主题
```bash
kafka-topics --create --topic linkme_binlog --partitions 1 --replication-factor 1 --bootstrap-server 192.168.1.11:30880
```
##### 列出所有主题
```bash
kafka-topics --list --bootstrap-server 192.168.1.11:30880
```
##### 主题详情
```bash
kafka-topics --describe --topic test --bootstrap-server 192.168.1.11:30880
```
##### 添加分区
```bash
kafka-topics --alter --topic linkme_binlog --partitions 3 --bootstrap-server 192.168.1.11:30880
```
##### 删除主题
```bash
kafka-topics --delete --topic linkme_binlog --bootstrap-server 192.168.1.11:30880
```
##### 读取主题消息
```bash
kafka-console-consumer --brokers 192.168.1.11:30880 --topic linkme_binlog --offset oldest
```
## 启动项目
### 本地启动
```bash
go build -o . && ./LinkMe
```
### docker启动
```bash
# 在项目根目录下运行
make rebuild
```
