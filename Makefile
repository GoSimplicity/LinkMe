build:
	go build -o bin/linkme cmd/main.go

run:
	go run cmd/main.go

run-dev:
	bin/linkme

dev: build run-dev

clean:
	rm -f bin/linkme

generate:
	wire ./...

docker-env-run:
	docker-compose -f docker-compose-env.yaml up -d

docker-env-down:
	docker-compose -f docker-compose-env.yaml down

docker-run: docker-build
	docker-compose -f docker-compose.yaml up -d

docker-down:
	docker-compose -f docker-compose.yaml down

docker-build:
	docker build -t linkme/gomodd:v1.22.3 .

docker-dev: docker-env-run docker-run

docker-dev-down: docker-env-down docker-down
