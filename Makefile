.PHONY: default deps build lint docker hadolint fmt
EXECUTABLE := whattowatchcmd
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')
GO_PACKAGES = $(shell go list ./... | grep -v /vendor/)

default: build

build:
	# Compiling...
	go build -ldflags "-X ${PACKAGENAME}/internal/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/internal/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}


deps:
	# Fetching dependancies...
	go get -d -v # Adding -u here will break CI

lint:
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v --timeout 5m

docker:
	docker build -t ${EXECUTABLE}  ./

hadolint:
	docker run -it --rm -v ${PWD}/Dockerfile:/Dockerfile hadolint/hadolint:latest hadolint --ignore DL3018 Dockerfile

fmt:
	@go fmt $(GO_PACKAGES)
