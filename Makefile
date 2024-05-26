#	@rm webook || true
#	@go mod tidy
#	@# x86架构
#	@#GOOS=linux GOARCH=arm go build -o webook .
#	@GOOS=linux go build -o webook .
#	@docker rmi -f flycash/webook:v0.0.1 || true
#	@ctr -n k8s.io image remove docker.io/library/flycash/webook:v0.0.1 || true
#	@docker build -t flycash/webook:v0.0.1 .
#	@docker save -o webook.tar flycash/webook
#	@ctr -n k8s.io image import webook.tar
#	@kubectl delete deployment webook-record-deployment || true
#	@kubectl delete deployment webook-record-mysql || true
#	@kubectl delete deployment webook-record-redis || true
#	@kubectl delete deployment webook-record-etcd || true
#	@kubectl delete svc webook-record-service || true
#	@kubectl delete svc webook-record-mysql || true
#	@kubectl delete svc webook-record-redis || true
#	@kubectl delete pvc webook-mysql-pvc || true
#	@kubectl delete pv webook-mysql-pv || true
#	@#kubectl delete -f yaml/ || true
#	@#kubectl delete ingress webook-record-ingress || true
#	@kubectl apply -f yaml/webook-mysql-deployment.yaml
#	@kubectl apply -f yaml/webook-redis-deployment.yaml
#	@kubectl apply -f yaml/webook-etcd-deployment.yaml
#	@kubectl apply -f yaml/webook-mysql-service.yaml
#	@kubectl apply -f yaml/webook-redis-service.yaml
#	@kubectl apply -f yaml/webook-etcd-service.yaml
#	@kubectl apply -f yaml/webook-mysql-pv.yaml
#	@kubectl apply -f yaml/webook-mysql-pvc.yaml
#	@#kubectl apply -f yaml/ || true
#	@kubectl apply -f yaml/webook-deployment.yaml
#	@kubectl apply -f yaml/webook-service.yaml
#	@cd config && etcdctl --endpoints=100.64.1.1:30883 put /webook "$(<config.yaml)" && cd ../
