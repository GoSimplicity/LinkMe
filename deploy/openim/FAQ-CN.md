# Docker Compose 常见问题及解决方法

[toc]

## 1. 配置文件管理

在使用 OpenIM 的新版本（version >= 3.2.0）时，对配置文件的管理变得尤为重要。配置文件不仅为应用程序提供了必要的运行参数，还确保了系统运行的稳定性和可靠性。

### 1.1 生成配置文件

为了生成配置文件，OpenIM 提供了两种方式。一种是通过 `Makefile`，另一种是直接执行初始化脚本。

#### 使用 Makefile

对于熟悉 Makefile 的开发者，这是一个快捷且友好的方式。只需要在项目根目录执行以下命令：

```bash
make init
```

这会触发 `Makefile` 中的相关命令，最终生成所需的配置文件。

#### 使用初始化脚本

对于不想使用 `Makefile` 的用户，或者那些对其不太熟悉的人，我们提供了一个更直接的方式来生成配置文件。只需执行以下命令：

```bash
./scripts/init-config.sh
```

无论选择哪种方式，都会生成相同的配置文件。因此，可以根据个人喜好和环境来选择最适合自己的方法。

### 1.2 验证配置文件

生成配置文件后，最好对其进行验证，确保它能够满足应用程序的要求。验证的标志如下：

[日志输出内容...]

这些日志输出确保了配置文件已正确生成并可以被 OpenIM 服务正确解析。

### 1.3 配置文件的修改与管理

配置文件通常不需要频繁修改。但在某些情况下，例如当更改数据库连接参数或修改其他关键参数时，可能需要调整它。

建议通过环境变量的方式去配置和管理 ~

建议在修改配置文件之前，先备份原始文件。这样，如果出现问题，可以很容易地回滚到原始状态。

此外，对于在团队中使用 OpenIM 的情况，建议使用版本控制系统（如 Git）来管理配置文件。这可以确保团队成员都使用相同的配置，并能够跟踪任何修改。



## 2. Docker Compose 不支持 gateway

Docker Compose 是一个工具，用于定义和运行多容器的 Docker 应用程序。有时，你可能会遇到不支持特定功能，如 `gateway` 的问题。下面是详细的指南，包括问题、原因、解决方法和调试技巧。

### 2.1 问题描述

当使用 Docker Compose 文件定义网络时，尝试设置 gateway 参数可能会出现以下错误：

```
ERROR: The Compose file './docker-compose.yaml' is invalid because:
networks.openim-server.ipam.config value Additional properties are not allowed ('gateway' was unexpected)
```

这意味着 Docker Compose 不支持你试图定义的 `gateway` 参数。

### 2.2 原因

Docker Compose 的某些版本可能不支持特定的网络属性，如 `gateway`。这可能是由于 Docker Compose 的版本过旧或配置文件语法有误。

### 2.3 解决方法

#### 检查版本

首先，确保你的 Docker Compose 版本是最新的。要查看版本，执行以下命令：

```
docker-compose version
```

如果你的版本过旧，建议更新到最新版本。

#### 校验配置文件

验证 `docker-compose.yaml` 文件的语法是否正确。确保缩进、空格和格式都是正确的。可以使用在线 YAML 校验工具进行检查。

#### 使用其他网络配置

如果不需要特定的 `gateway` 设置，可以考虑删除或更改它。另外，如果你只是想要为容器定义一个静态 IP，可以使用 `ipv4_address` 属性。

### 2.4 调试与帮助

如果上述方法仍不能解决问题，以下是一些调试技巧和帮助指南：

#### 查看 Docker 文档

Docker 官方文档是一个宝贵的资源。确保你已经阅读了关于 [Docker Compose 文件的官方文档](https://docs.docker.com/compose/compose-file/)。

#### 使用更详细的日志

运行 `docker-compose` 时使用 `-v` 参数可以获得更详细的日志输出，这可能有助于识别问题的根源。

```
docker-compose -v up
```

#### 访问社区和论坛

Docker 有一个非常活跃的社区。如果你遇到问题，可以考虑在 [Docker 论坛](https://forums.docker.com/) 发布问题或搜索其他用户是否有相同的问题。



## 3. MySQL 连接失败

!!! 最新版本已近移除 mysql 了

在使用 Docker 运行的应用程序中，MySQL 连接失败是一个常见的问题。该问题可能由多种原因引起，以下是一份全面的指南，旨在帮助你解决 MySQL 连接问题。

### 3.1 问题描述

当你的应用程序或服务尝试连接到 MySQL 容器时，可能会遇到以下错误：

```
[error] failed to initialize database, got error dial tcp 172.28.0.2:13306: connect: connection refused
```

这意味着你的应用程序无法建立到 MySQL 的连接。

### 3.2 常见原因与解决方案

#### MySQL 容器未启动

**检查方法:** 使用 `docker ps` 命令查看所有正在运行的容器。

**解决方法:** 如果你没有看到 MySQL 容器，请确保它已经启动。

```
docker-compose up -d mysql
```

#### 错误的 MySQL 地址或端口

**检查方法:** 检查应用程序的配置文件，确保 MySQL 的地址和端口设置正确。

**解决方法:** 如果使用了默认的 Docker Compose 设置，地址应为 `mysql` (容器名) 并且默认端口是 `3306`。

#### MySQL 用户权限问题

**检查方法:** 登录 MySQL 并检查用户权限。

**解决方法:** 确保连接的 MySQL 用户具有足够的权限。你可以考虑为应用程序创建一个专用用户并授予必要的权限。

#### MySQL 的 `bind-address`

**检查方法:** 如果 MySQL 仅绑定到 `127.0.0.1`，则只能从容器内部访问它。

**解决方法:** 更改 MySQL 的 `bind-address` 到 `0.0.0.0` 以允许外部连接。

#### 网络问题

**检查方法:** 使用 `docker network inspect` 检查容器的网络设置。

**解决方法:** 确保应用程序和 MySQL 容器在同一网络上。

### 3.3 调试方法与帮助

#### 查看 MySQL 日志

查看 MySQL 容器的日志可能会提供有关连接失败的更多信息。

```
docker logs <mysql_container_name>
```

#### 使用 MySQL 客户端进行测试

使用 MySQL 客户端直接连接到数据库可以帮助定位问题。

```
mysql -h <mysql_container_ip> -P 3306 -u <username> -p
```

#### 检查防火墙设置

确保没有防火墙或网络策略阻止应用程序与 MySQL 容器之间的通信。

### 3.4 其他可能的问题

#### 使用旧版本的 Docker 或 Docker Compose

确保你使用的是 Docker 和 Docker Compose 的最新版本。旧版本可能存在已知的连接问题。

#### 数据库没有初始化

如果这是 MySQL 容器的首次启动，可能需要一些时间来初始化数据库。

#### 容器间的时间同步问题

确保所有容器的系统时间都是同步的，因为时间不同步可能会导致认证问题。




## 4. Kafka 错误

Kafka 是一个流行的消息传递系统，但与所有技术一样，你可能会遇到一些常见问题。以下是详细的指南，提供了关于 Kafka 错误的信息，以及如何解决这些问题。

### 4.1 问题描述

当尝试启动或与 Kafka 交互时，你可能会遇到以下错误：

```
Starting Kafka failed: kafka doesn't contain topic:offlineMsgToMongoMysql: 6000 ComponentStartErr
```

此错误表明 Kafka 服务没有预期的topical或组件没有正确启动。

### 4.2 常见原因与解决方案

#### Kafka 未运行或启动失败

**检查方法:** 使用 `docker ps` 或 `docker-compose ps` 查看 Kafka 容器的状态。

**解决方法:** 如果 Kafka 未运行，请确保使用正确的命令启动它。例如，使用 `docker-compose up -d kafka`。

#### topical不存在

**检查方法:** 使用 Kafka 的命令行工具查看所有可用的topical。

**解决方法:** 如果topical不存在，你需要创建它。你可以使用 `kafka-topics.sh` 脚本来创建新topical。

#### Kafka 配置问题

**检查方法:** 检查 Kafka 的配置文件，确保所有的配置项都设置正确。

**解决方法:** 根据你的需求调整 Kafka 的配置并重新启动服务。

### 4.3 调试方法与帮助

#### 查看 Kafka 日志

Kafka 容器的日志可能包含有用的信息。你可以使用以下命令查看它：

```
docker logs <kafka_container_name>
```

#### 使用 Kafka 命令行工具

Kafka 附带了一系列的命令行工具，这些工具可以帮助你管理和调试服务。确保你熟悉如何使用它们，特别是 `kafka-topics.sh` 和 `kafka-console-producer/consumer.sh`。

#### 确保 Zookeeper 正常运行

Kafka 依赖于 Zookeeper，所以确保 Zookeeper 也在正常运行。

### 4.4 其他可能的问题

#### 网络问题

确保 Kafka 和其他服务（如 Zookeeper）都在同一个 Docker 网络上，并且容器之间可以相互通信。

#### 存储问题

确保 Kafka 容器有足够的磁盘空间。如果磁盘空间不足，Kafka 可能会遇到问题。

#### 版本不兼容

确保你使用的 Kafka 客户端版本与 Kafka 服务版本兼容。



## 5. 网络错误

在使用 Docker 和容器化的应用程序时，网络问题可能是最常见的问题之一。从 IP 地址冲突到容器间连接失败，网络错误的原因和解决方案是多种多样的。

### 5.1 常见的网络错误

#### 错误 1: Invalid address

**问题描述:**

```
Error response from daemon: Invalid address 172.28.0.12: It does not belong to any of this network's subnets
```

这个错误通常意味着你试图给一个容器分配一个不属于 Docker 网络子网的 IP 地址。

出现这个问题的原因是我们已经创建过这种网络：

```bash
$ docker network ls
NETWORK ID     NAME                    DRIVER    SCOPE
1a3336557aa2   bridge                  bridge    local
0fe6ea241277   host                    host      local
d97c634eb4eb   none                    null      local
9e4540aa4961   open-im-server_server   bridge    local
```

可以看到一个类型 bridge 的 open-im-server_server 网络。

> 两种方法可以解决

**解决方案:**

1. 使用 `docker network inspect [network_name]` 检查网络的子网范围。
2. 确保为容器分配的 IP 地址在这个范围内。

#### 错误 2: Pool overlaps

**问题描述:**

```
failed to create network example_openim-server: Error response from daemon: Pool overlaps with other one on this address space
```

这意味着你试图创建一个与现有网络有重叠 IP 地址范围的新网络。

**解决方案:**

1. 更改新网络的 IP 地址范围。
2. 或者删除现有的重叠网络（在确保其不再需要的情况下）。

### 5.2 调试网络问题的方法

#### 1. `docker network ls`

列出所有的 Docker 网络，这样你可以看到是否有任何预期之外的网络或重复的网络。

#### 2. `docker network inspect [network_name]`

检查特定的 Docker 网络配置，特别是 IP 地址范围和连接到该网络的容器。

#### 3. `ping` 和 `curl`

从一个容器内部 ping 另一个容器的 IP 地址或使用 curl 尝试连接到另一个容器的服务。这可以帮助你确定网络连接问题的位置。

#### 4. 查看容器日志

使用 `docker logs [container_name]` 检查容器的日志，可能会有一些与网络相关的错误或警告。

### 5.3 其他可能的网络问题

#### DNS 解析问题

容器可能无法解析其他容器的域名。确保你的容器使用了正确的 DNS 设置，并且可以访问 DNS 服务器。

#### 端口未暴露或绑定

如果你的服务在容器内部运行，但无法从外部访问，确保你已经在 Dockerfile 中使用 `EXPOSE` 指令暴露了正确的端口，并在启动容器时绑定了这些端口。

#### 防火墙或安全组

确保任何外部的防火墙或安全组都允许必要的流量通过。



## 6. 其他问题的排查

当你使用开源项目或任何其他软件时，难免会遇到一些不可预测的问题。如何优雅地排查和解决问题是每个开发者和用户都应该掌握的重要技能。

### 6.1 明确问题描述

首先，要确保你真正理解了问题。随意地尝试各种解决方案，而不首先定义问题是一种时间浪费的策略。

- **收集错误日志**：几乎所有的应用程序或软件都有日志记录功能。始终查看日志以获取有关问题的更多详细信息。
- **重现问题**：在尝试解决问题之前，了解如何重现它是很重要的。如果一个问题不能被可靠地重现，它很难被解决。

### 6.2 分隔排除法

一种有效的故障排除策略是分隔和排除。这意味着你将系统拆分为不同的部分，并单独测试每一部分，以确定问题出在哪里。

- **单独运行组件**：例如，如果你在使用多个服务的系统中遇到问题，尝试单独运行每个服务来看哪个服务出了问题。
- **使用最小化的配置**：如果可能，使用最基本的配置启动应用程序或服务，然后逐渐添加更多的配置选项，直到你可以重现问题。

### 6.3 使用开源社区资源

- **查找已知的问题**：大多数开源项目都有一个issue跟踪器，如GitHub的Issues。首先查看那里，看看你的问题是否已经被其他人报告过。
- **提问的技巧**：如果你决定询问社区，确保你的问题是明确的、具体的，并附带足够的详细信息。包括错误消息、你的环境信息和你已经尝试过的解决方案。

### 6.4 使用调试工具

- **代码调试**：如果你对代码感到舒适，使用调试器来逐步执行代码可以帮助你更快地找到问题。
- **网络调试**：对于网络问题，工具如 `ping`, `traceroute`, `netstat` 和 `wireshark` 可以非常有用。

### 6.5 发现问题后的步骤

一旦你找到了问题，以下是一些建议的下一步：

- **查找现有的修复程序**：可能有人已经为你的问题找到了一个修复程序或解决方案。
- **修复问题**：如果你有技能和资源，你可以尝试自己修复问题。
- **报告问题**：即使你自己解决了问题，也要向开源社区报告它，这样其他人可以从你的发现中受益。

### 6.6 保持耐心

最后但同样重要的是，保持耐心和开放的心态。遇到问题是软件开发的一个普遍现象，学习如何有效地解决它们可以使你成为一个更好的开发者。

总的来说，优雅地排查和解决问题需要时间、实践和耐心，但随着时间的推移，你将发展出自己的策略和技术，使这个过程变得更加容易和直观。

### 6.7 案例

排查问题其实不仅仅是开发者的技能，很多问题很简单，在网上就能找到答案，比如说，我针对下面一个问题的解决思路：


**源码启动 -  openim-msgtransfer 报错:**

如下：

```bash
## Check OpenIM service name
!!! [0902 21:10:22] Call tree:
!!! [0902 21:10:22]  1: /root/workspaces/openim/openim-server/scripts/install/openim-msgtransfer.sh:154 openim::msgtransfer::check(...)
!!! [0902 21:10:22]  2: /root/workspaces/openim/openim-server/scripts/check-all.sh:80 source(...)
!!! Error in /root/workspaces/openim/openim-server/scripts/install/openim-msgtransfer.sh:66 
  Error in /root/workspaces/openim/openim-server/scripts/install/openim-msgtransfer.sh:66. 'PIDS=$(pgrep -f "${OPENIM_OUTPUT_HOSTBIN}/openim-msgtransfer")' exited with status 1
Call stack:
  1: /root/workspaces/openim/openim-server/scripts/install/openim-msgtransfer.sh:66 openim::msgtransfer::check(...)
  2: /root/workspaces/openim/openim-server/scripts/install/openim-msgtransfer.sh:154 source(...)
  3: /root/workspaces/openim/openim-server/scripts/check-all.sh:80 main(...)
Exiting with status 1
make[1]: *** [scripts/make-rules/golang.mk:120: go.check] Error 1
make: *** [Makefile:113: check] Error 2
```

首先，是代码层面问题，通过 dlv 调试：

```bash
root@PS2023EVRHNCXG:~/workspaces/openim/openim-server# dlv exec _output/bin/platforms/linux/amd64/openim-msgtransfer
Warning: no debug info found, some functionality will be missing such as stack traces and variable evaluation.
Type 'help' for list of commands.
(dlv) c
configFolderPath:
use config /root/workspaces/openim/openim-server/config/config.yaml
use config /root/workspaces/openim/openim-server/config/notification.yaml
mongo: mongodb://root:openIM123@172.29.0.1:37017/openim_v3?maxPoolSize=100&authSource=admin
start msg transfer prometheusPort: 0
> [unrecovered-panic] runtime.fatalpanic() /snap/go/10319/src/runtime/panic.go:1175 (hits total:2) (PC: 0x43b940)
> [unrecovered-panic] runtime.fatalpanic() /snap/go/10319/src/runtime/panic.go:1175 (hits total:2) (PC: 0x43b940)
  1170: // fatalpanic implements an unrecoverable panic. It is like fatalthrow, except
  1171: // that if msgs != nil, fatalpanic also prints panic messages and decrements
  1172: // runningPanicDefers once main is blocked from exiting.
  1173: //
  1174: //go:nosplit
=>1175: func fatalpanic(msgs *_panic) {
  1176:         pc := getcallerpc()
  1177:         sp := getcallersp()
  1178:         gp := getg()
  1179:         var docrash bool
  1180:         // Switch to the system stack to avoid any stack growth, which
(dlv) c
panic: dial tcp: lookup 314be2d68378 on 172.31.176.1:53: no such host

goroutine 252 [running]:
github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka.(*MConsumerGroup).RegisterHandleAndConsumer(0xc000026880, {0x1564d98, 0xc000274380})
        /root/workspaces/openim/openim-server/pkg/common/kafka/consumer_group.go:71 +0x128
created by github.com/OpenIMSDK/Open-IM-Server/internal/msgtransfer.(*MsgTransfer).Start in goroutine 1
        /root/workspaces/openim/openim-server/internal/msgtransfer/init.go:117 +0x18d
panic: dial tcp: lookup 314be2d68378 on 172.31.176.1:53: no such host

goroutine 253 [running]:
github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka.(*MConsumerGroup).RegisterHandleAndConsumer(0xc000388cc0, {0x1564dc8, 0xc0000172c0})
        /root/workspaces/openim/openim-server/pkg/common/kafka/consumer_group.go:71 +0x128
created by github.com/OpenIMSDK/Open-IM-Server/internal/msgtransfer.(*MsgTransfer).Start in goroutine 1
        /root/workspaces/openim/openim-server/internal/msgtransfer/init.go:118 +0x20e
Process 6315 has exited with status 2
```

判断报错，好像是 `*:53` 这个端口有问题于是，开始网上找问题：

1. 阅读 docker hub 配置

2. 阅读官方文档
3. 查找 github issue 

终于找到了一个 issue：https://github.com/segmentio/kafka-go/issues/671

但是很遗憾，依旧没有告诉我解决方案是什么，只是存在这个问题，收获到了一点，判断到了 172.31.176.1:53: no such host，之前的 `127.0.0.1:53`

虽然如此，我也知道了改成 host 模式后肯定不会有问题了，或者是全部使用 docker compose 的话也不会有问题。

```bash
RUN echo listener.security.protocol.map=INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT >> config/server.properties \
&& echo advertised.listeners=INSIDE://localhost:9092,OUTSIDE://localhost:29092 >> config/server.properties \
&& echo listeners=INSIDE://0.0.0.0:9092,OUTSIDE://0.0.0.0:29092 >> config/server.properties \
&& echo inter.broker.listener.name=INSIDE >> config/server.properties
```

于是再回头阅读官方文档，结合网络知识，用最终的配置文件，并且测试：


```bash
  kafka:
    image: 'bitnami/kafka:3.5.1'
    container_name: kafka
    user: root
    restart: always
    ports:
    - "${KAFKA_PORT}:9092"
    volumes:
      - ./scripts/create_topic.sh:/opt/bitnami/kafka/create_topic.sh
      - ${DATA_DIR}/components/kafka:/bitnami/kafka
    command: >
      bash -c "
      /opt/bitnami/scripts/kafka/run.sh & sleep 5; /opt/bitnami/kafka/create_topic.sh; wait
      "
    environment:
       - TZ=Asia/Shanghai
       - KAFKA_CFG_NODE_ID=0
       - KAFKA_CFG_PROCESS_ROLES=controller,broker
       - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@<your_host>:9093
       - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:9094
       - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092,EXTERNAL://${DOCKER_BRIDGE_GATEWAY}:${KAFKA_PORT}
       - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT
       - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
    networks:
      openim-server:
        ipv4_address: ${KAFKA_NETWORK_ADDRESS}
```

最终解决问题