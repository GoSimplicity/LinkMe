# OpenIM Docker 使用说明 📘

> **文档资源** 📚

+ [官方部署指南](https://docs.openim.io/guides/gettingstarted/dockercompose)

## :busts_in_silhouette: Community

+ 💬 [关注推特](https://twitter.com/founder_im63606)
+ 🚀 [进slack频道](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)
+ :eyes: [进微信群](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## 环境准备 🌍

- 在服务器上安装带有 Compose 插件的 Docker 或 docker-compose。安装详情请访问 [Docker Compose 安装指南](https://docs.docker.com/compose/install/linux/)。

## 仓库克隆 🗂️

```bash
git clone https://github.com/openimsdk/openim-docker
```

## 配置修改 🔧

- 修改 `.env` 文件，配置外网 IP。如果使用域名，需配置 Nginx。

  ```plaintext
  # 设置 MinIO 服务的外网访问地址（IP或域名）
  MINIO_EXTERNAL_ADDRESS="http://external_ip:10005" 
  ```

- 其他配置请参考 .env 中的注释

## 服务启动 🚀

- 启动服务：
```bash
docker compose up -d
```

- 停止服务：
```bash
docker compose down
```

- 查看日志：
```bash
docker logs -f openim-server
docker logs -f openim-chat
```

## 快速体验 ⚡

快速体验 OpenIM 服务，请访问 [快速测试服务器指南](https://docs.openim.io/guides/gettingStarted/quickTestServer)。

