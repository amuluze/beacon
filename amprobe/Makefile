APP 			= amprobe
SERVER  		= ./cmd/${APP}
RELEASE_SERVER 	= release/${APP}/${APP}

GO 				= go
CONFIG_FILE := $(shell pwd)/configs/config.toml

.PHONY: dev
# run dev
dev:
	@$(GO) run $(SERVER)/main.go run -c $(CONFIG_FILE)

.PHONY: bin
# build bin
bin:
	$(GO) build -o $(APP) $(SERVER)

.PHONY: wire
# generate wire
wire:
	wire ./service/

.PHONY: amd64
# docker amd64 image
amd64:
	docker buildx build --platform=linux/amd64 -o type=docker -t amuluze/amprobe:v1.3.4 .

.PHONY: arm64
# docker arm64 image
arm64:
	docker buildx build --platform=linux/arm64 -o type=docker -t amuluze/amprobe_arm:v1.3.4 -f ./Dockerfile.arm .
