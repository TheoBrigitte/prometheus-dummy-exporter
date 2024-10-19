GOPATH     ?= $(shell go env GOPATH)
GORELEASER ?= $(GOPATH)/bin/goreleaser
VERSION    := v$(shell cat VERSION)

.PHONY: setup test build build-snapshot sync-tag release docker-build docker-release
all: setup test build build-snapshot sync-tag release docker-build docker-release

test:
	@echo '>> unit test'
	go test ./...

bench:
	@echo '>> benchmark'
	go test -bench=. -benchtime=5x ./...

build:
	@echo '>> build'
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=" \
	-X github.com/prometheus/common/version.Version=$(shell cat VERSION) \
	-X github.com/prometheus/common/version.Revision=$(shell git rev-parse HEAD) \
	-X github.com/prometheus/common/version.Branch=$(shell git rev-parse --abbrev-ref HEAD) \
	-X 'github.com/prometheus/common/version.BuildUser=$(shell git config get user.name) ($(shell git config get user.email))' \
	-X github.com/prometheus/common/version.BuildDate=$(shell date --utc +%FT%T)" \
	./

build-snapshot: $(GORELEASER)
	@echo '>> cross-build for testing'
	BUILD_BRANCH=$(shell git rev-parse --abbrev-ref HEAD) \
	BUILD_USER=$(shell whoami) \
	BUILD_HOST=$(shell hostname) \
	BUILD_DATE=$(shell date +%Y%m%d-%H:%M:%S) \
	$(GORELEASER) release --snapshot --rm-dist --debug

release: $(GORELEASER)
	@echo '>> release'
	BUILD_BRANCH=$(shell git rev-parse --abbrev-ref HEAD) \
	BUILD_USER=$(shell whoami) \
	BUILD_HOST=$(shell hostname) \
	BUILD_DATE=$(shell date +%Y%m%d-%H:%M:%S) \
	$(GORELEASER) release --rm-dist --debug

docker-build: build
	@echo '>> build docker image'
	@docker build -t theo01/prometheus-dummy-exporter:$(shell cat VERSION) .
	@docker build -t theo01/prometheus-dummy-exporter:latest .

docker-release: docker-build
	@echo '>> release docker image'
	@docker push theo01/prometheus-dummy-exporter:$(shell cat VERSION)
	@docker push theo01/prometheus-dummy-exporter:latest

$(GORELEASER):
	@wget -O - "https://github.com/goreleaser/goreleaser/releases/download/v0.98.0/goreleaser_$(shell uname -o | cut -d'/' -f2)_$(shell uname -m).tar.gz" | tar xvzf - -C /tmp
	@mv /tmp/goreleaser $(GOPATH)/bin
