TARGET 			= server
SERVER  		= ./cmd/${TARGET}

GO 				= go
GOBUILD         = $(GO) build
DEV_CONFIG_FILE := $(shell pwd)/configs/config.dev.toml

.PHONY: dev
dev:
	@$(GO) run $(SERVER)/main.go run -c $(DEV_CONFIG_FILE)

.PHONY: wire
wire:
	wire ./service/

.PHONY: arm64
arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(TARGET) $(SERVER)

.PHONY: amd64
amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(TARGET) $(SERVER)
