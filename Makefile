BUILD_DATE := $(shell date +'%d.%m.%Y %H:%M:%S')
BUILD_COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: all
all: clean generate build test

.PHONY: generate
generate:
	@echo "\n### $@"
	go generate ./...

.PHONY: build
build: build_bot build_server

.PHONY: build_bot
build_bot:
	@echo "\n### $@"
	@mkdir -p ./bin
	@cd cmd/myhealthbot && \
	go build \
	-ldflags "-X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" \
	-o ../../bin/myhealthbot .	 

.PHONY: build_server
build_server:
	@echo "\n### $@"
	@mkdir -p ./bin
	@cd cmd/myhealthserver && \
	go build \
	-ldflags "-X 'main.buildDate=$(BUILD_DATE)' -X main.buildCommit=$(BUILD_COMMIT)" \
	-o ../../bin/myhealthserver .

.PHONY: test
test:
	@echo "\n### $@"
	go test ./... -v --count 1

.PHONY: clean
clean:
	@echo "\n### $@"
	@rm -rf ./bin		 