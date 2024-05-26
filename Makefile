linkme:
	@rm LinkMe || true
	@go mod tidy
	@# x86架构
	@#GOOS=linux GOARCH=arm go build -o webook .
	@GOOS=linux go build -o LinkMe .
	@docker rmi -f linkme:v0.0.1 || true
	@ctr -n k8s.io image remove docker.io/library/linkme:v0.0.1 || true
	@docker build -t linkme:v0.0.1 .
	@docker save -o linkme.tar linkme
	@ctr -n k8s.io image import linkme.tar
	@kubectl delete -f yaml/ || true
	@kubectl apply -f yaml/ || true
