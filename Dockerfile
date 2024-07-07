# 使用官方的 Golang 镜像作为构建阶段
FROM golang:1.22.3 AS build-stage
# 设置 Go 代理为国内源
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN go mod tidy
# 静态编译可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
# 验证可执行文件是否被正确生成

# 使用 Alpine 镜像作为生产阶段
FROM alpine:latest
WORKDIR /root/
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
# 复制构建阶段的可执行文件到生产阶段
COPY --from=build-stage /app/main .
COPY --from=build-stage /app/config ./config
# 确保可执行文件具有执行权限
RUN chmod +x ./main
# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
# 设置时区（以 Asia/Shanghai 为例）
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
# 设置编码
ENV LANG C.UTF-8
EXPOSE 9999
CMD ["./main"]
