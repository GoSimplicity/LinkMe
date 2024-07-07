IMAGE_NAME=linkme-backend
# 拉取最新代码
pull:
	git pull origin main
# 构建 Docker 镜像
build: pull
	docker build -t $(IMAGE_NAME) .
# 启动容器
run: build
	docker run -d -p 9999:9999 --name $(IMAGE_NAME) $(IMAGE_NAME)
# 停止并删除容器
clean:
	docker stop $(IMAGE_NAME) || true
	docker rm $(IMAGE_NAME) || true
# 重新构建并启动容器
rebuild: clean build run