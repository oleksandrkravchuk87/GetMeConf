
export PDB_HOST=horton.elephantsql.com
export PDB_PORT=5432
export PDB_USER=dlxifkbx
export PDB_PASSWORD=L7Cey-ucPY4L3T6VFlFdNykNE4jO0VjV
export PDB_NAME=dlxifkbx
export MAX_OPENED_CONNECTIONS_TO_DB=5
export MAX_IDLE_CONNECTIONS_TO_DB=0
export MB_CONN_MAX_LIFETIME_MINUTES=30

export CACHE_EXPIRATION_TIME=5
export CACHE_CLEANUP_INTERVAL=10

export PORT=3000

all: dependencies build

.PHONY: build
build:
	echo "Build"
	go build -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service ./service

.PHONY: run
run:
	echo "Running service"
	echo ${PDB_HOST}
	go build -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/cService ./service
	./bin/cService

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

