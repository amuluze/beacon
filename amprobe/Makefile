APP 			= amprobe
SERVER  		= ./cmd/${APP}
RELEASE_SERVER 	= release/${APP}/${APP}

IMAGE = amuluze/amprobe
TAG = v1.3.8

GO 				= go
CONFIG_FILE := $(shell pwd)/configs/config.toml

.PHONY: dev
# run dev
dev:
	@$(GO) run $(SERVER)/main.go run -c $(CONFIG_FILE)

.PHONY: wire
# generate wire
wire:
	wire ./service/

.PHONY: bin
# build bin
bin: wire
	$(GO) build -o $(APP) $(SERVER)


.PHONY: amd64
# docker amd64 image
amd64: bin
	# docker build -t $(IMAGE):$(TAG) -f ./Dockerfile.amd .
	docker buildx build --platform linux/amd64 -t $(IMAGE):$(TAG) -o type=docker .

.PHONY: arm64
# docker arm64 image
arm64: bin
	# docker build -t $(IMAGE):$(TAG) -f ./Dockerfile.arm .
	docker buildx build --platform linux/arm64 -t $(IMAGE):$(TAG) -o type=docker .

.PHONY: build
build:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE):$(TAG) . --push
