APP 			= amprobe
SERVER  		= ./cmd/${APP}

STAMP    = $(shell date +%s)
VERSION = $(shell cat VERSION)
GIT_HASH  = $(shell git rev-parse --short=8 HEAD)
GIT_TAG   = $(shell git describe --tags --abbrev=0)
GIT_BRANCH = $(shell which git >/dev/null 2>&1 && git rev-parse --abbrev-ref HEAD)

IMAGE = amuluze/amprobe
TAG = $(VERSION)

GO 				= go
CONFIG_FILE := $(shell pwd)/configs/config.toml
MODEL_FILE := $(shell pwd)/configs/model.conf
DEV_CONFIG_FILE := $(shell pwd)/configs/config.dev.toml
BUILD_FLAGS 	:= -ldflags "-w -s -X main.BuildStamp=$(STAMP) -X main.GitHash=$(GIT_HASH) -X main.GitBranch=$(GIT_BRANCH) -X main.Version=$(VERSION)"

.PHONY: dev
# run dev
dev:
	@$(GO) run $(SERVER)/main.go run -c $(DEV_CONFIG_FILE) -m $(MODEL_FILE)

.PHONY: wire
# generate wire
wire:
	wire ./service/

.PHONY: bin
# build bin
bin:
	$(GO) build $(BUILD_FLAGS) -o $(APP) $(SERVER)

.PHONY: amd64
# docker amd64 image
amd64: bin
	docker buildx build --platform linux/amd64 -t $(IMAGE):$(TAG) -o type=docker .

.PHONY: arm64
# docker arm64 image
arm64: bin
	docker buildx build --platform linux/arm64 -t $(IMAGE):$(TAG) -o type=docker .

.PHONY: build
build:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE):$(TAG) . --push
