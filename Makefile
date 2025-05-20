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
BUILD_FLAGS 	:= -ldflags "-w -s -X main.BuildStamp=$(STAMP) -X main.GitHash=$(GIT_HASH) -X main.GitBranch=$(GIT_BRANCH) -X main.Version=$(VERSION)"

GOBINDATA   = go-bindata
BINDATAPKG  = assets
BINDATAPATH = $(BINDATAPKG)/bindata.go

CONFIG_FILE := $(shell pwd)/assets/resources/configs/config.toml
MODEL_FILE := $(shell pwd)/assets/resources/configs/model.conf

.PHONY: dev
dev:
	@$(GO) run $(SERVER)/main.go run -c $(CONFIG_FILE) -m $(MODEL_FILE)

.PHONY: assets
assets:
	@$(GOBINDATA) --nomemcopy -pkg $(BINDATAPKG) -o $(BINDATAPATH) assets/resources/...

.PHONY: wire
wire:
	wire ./service/

.PHONY: bin
bin:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(BUILD_FLAGS) -o $(APP) $(SERVER)
