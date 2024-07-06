IMAGE_NAME=linkme-backend
build:
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