REGISTRY = registry.cn-hangzhou.aliyuncs.com
IMAGE    = ${REGISTRY}/amuluze/amprobe-offical
TAG      = v0.1.7

.PHONY: amd64
# docker amd64 image
amd64:
	@docker buildx build --platform linux/amd64 -t $(IMAGE):$(TAG) -o type=docker .

.PHONY: arm64
# docker arm64 image
arm64:
	DOCKER_BUILDKIT=0 @docker buildx build --platform linux/arm64 -t $(IMAGE):$(TAG) -o type=docker .

.PHONY: push
push:
	@docker buildx build --platform linux/amd64 -t ${IMAGE}:${TAG} . --push
