EXECUTABLE := whatswatchcmd
GITVERSION := $(shell git describe --dirty --always --tags --long)
GOPATH ?= ${HOME}/go
PACKAGENAME := $(shell go list -m -f '{{.Path}}')

.PHONY: default
default: ${EXECUTABLE}

.PHONY: ${EXECUTABLE}
${EXECUTABLE}:
	# Compiling...
	go build -ldflags "-X ${PACKAGENAME}/internal/conf.Executable=${EXECUTABLE} -X ${PACKAGENAME}/internal/conf.GitVersion=${GITVERSION}" -o ${EXECUTABLE}


.PHONY: deps
deps:
	# Fetching dependancies...
	go get -d -v # Adding -u here will break CI

.PHONY: lint
lint:
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v --timeout 5m

.PHONY: hadolint
hadolint:
	docker run -it --rm -v ${PWD}/Dockerfile:/Dockerfile hadolint/hadolint:latest hadolint --ignore DL3018 Dockerfile

