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
	# docker build -t amuluze/amprobe:v1.3.4 .
	docker buildx build --platform linux/amd64 -t amuluze/amprobe:v1.3.4 -o type=docker .

.PHONY: arm64
# docker arm64 image
arm64:
	# docker build -t amuluze/amprobe_arm:v1.3.4 -f ./Dockerfile.arm .
	docker buildx build --platform linux/arm64 -t amuluze/amprobe_arm:v1.3.4 -o type=docker .

.PHONY: build
build:
	docker buildx build --platform linux/amd64,linux/arm64 -t amuluze/amprobe:v1.3.4 . --push
