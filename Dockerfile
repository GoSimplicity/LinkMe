# 使用官方的 Golang 镜像作为构建阶段
FROM golang:1.18 AS build-stage
# 设置 Go 代理为国内源
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .

# 使用一个小的基础镜像作为生产阶段
FROM alpine:latest
WORKDIR /root/
COPY --from=build-stage /app/main .
EXPOSE 9999
CMD ["./main"]