.PHONY: all build-cronserver build-cronserver-img run-cronserver gotool clean help

BIN_DIR=./bin
CRONSERVER=axis-cronserver
CRONSERVER_IMAGE=yintech/axis-cronserver:v1
CRONCLIENT=axis-cronclient
CRONCLIENT_IMAGE=yintech/axis-cronclient:v1

all: build-cronserver

build-cronserver:
	@if [ ! -d ${BIN_DIR} ]; then mkdir ${BIN_DIR}; fi
	go build -o ${BIN_DIR}/${CRONSERVER} cmd/mq/internal/aliyunsls.go

build-cronserver-img:
	docker build -t ${CRONSERVER_IMAGE} -f cronserver-Dockerfile .

run-cronserver:
	@go run cmd/mq/internal/aliyunsls.go -f cmd/mq/etc/axis.yaml

build-cronclient:
	@if [ ! -d ${BIN_DIR} ]; then mkdir ${BIN_DIR}; fi
	go build -o ${BIN_DIR}/${CRONCLIENT} cmd/cron/cron.go

build-cronclient-img:
	docker build -t ${CRONCLIENT_IMAGE} -f cronclient-Dockerfile .

run-cronclient:
	@go run cmd/cron/cron.go

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BIN_DIR}/${CRONSERVER} ]; then rm ${BIN_DIR}/${CRONSERVER}; fi

help:
	@echo "make build-cronserver - 编译生成定时任务服务端二进制文件"
	@echo "make build-cronserver-img - 构建定时任务服务端docker image"
	@echo "make run-cronserver - 直接运行Go定时任务服务端代码"
	@echo "make build-cronclient - 编译生成定时任务客户端二进制文件"
	@echo "make build-cronclient-img - 构建定时任务客户端docker image"
	@echo "make run-cronclient - 直接运行Go定时任务客户端代码"
	@echo "make clean - 移除所有二进制文件"
	@echo "make gotool - 运行Go工具'fmt'和'vet'"