.PHONY: build
build:
	echo "Build"
	cd service; \
    go build -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service .

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