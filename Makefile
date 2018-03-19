include .env
export

all: dependencies build

.PHONY: build
build:
	echo "Build"
	go build -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service ./service

.PHONY: run
run:
	echo "Running service"
	go build -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service ./service
	./bin/service

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	dep ensure

install dep:
	echo    "Installing dep"
	curl    https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

.PHONY: tests
tests:
	echo "Tests"
	go test ./service
	go test ./repository

docker-build:
	docker build -t configservice . && docker run --net=${DOCKER_NET_NAME} -p ${SERVICE_PORT}:${SERVICE_PORT} --env-file .env configservice