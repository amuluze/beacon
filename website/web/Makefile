REGISTRY = registry.cn-hangzhou.aliyuncs.com
IMAGE    = ${REGISTRY}/amuluze/amprobe-offical
TAG      = v0.1.0

.PHONY: install
install:
	@pnpm install

.PHONY: clean
clean:
	rm -rf node_modules pnpm-lock.yaml

.PHONY: dev
dev: install
	@pnpm run dev

.PHONY: build
build: install
	@pnpm run build
