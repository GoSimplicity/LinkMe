IMAGE_NAME=linkme/gomodd:v1.22.3

# 创建数据目录并提权
init:
	mkdir -p ./data/kafka/data && chmod -R 777 ./data/kafka
	mkdir -p ./data/es/data && chmod -R 777 ./data/es

# 启动项目依赖
env-up:
	docker-compose -f docker-compose-env.yaml up -d

# 构建 Docker 镜像
build:
	docker build -t $(IMAGE_NAME) .

# 启动项目
up: build
	docker-compose up -d

# 停止项目
down:
	docker-compose down --remove-orphans

# 重新构建并启动
rebuild: down build up

# 清理所有容器和数据
clean:
	docker-compose -f docker-compose-env.yaml down --remove-orphans
	docker-compose down --remove-orphans
	rm -rf ./data

# 一键部署 - 执行完整的部署流程
deploy: init env-up build up
	@echo "项目部署完成!"
	@echo "访问 http://localhost:8888 查看项目"

# 一键重新部署 - 清理后重新部署
redeploy: clean init deploy

# 一键更新 - 拉取最新代码并重新部署
update:
	git pull
	make redeploy

# 显示所有容器状态
status:
	docker-compose ps
	docker-compose -f docker-compose-env.yaml ps

# 查看项目日志
logs:
	docker-compose logs -f