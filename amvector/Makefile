VECTOR 			= amvector
VECTOR_ARM 		= amvector_arm
VECTOR_SERVER   = ./cmd/${VECTOR}

.PHONY: amd64
# build bin
amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-gnu-g++ CGO_LDFLAGS="-static" go build -ldflags="-s -w" -o $(VECTOR) $(VECTOR_SERVER)

.PHONY: arm64
# build arm
arm64:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++ CGO_LDFLAGS="-static" go build -ldflags="-s -w" -o $(VECTOR_ARM) $(VECTOR_SERVER)

.PHONY: wire
# generate wire
wire:
	wire ./service/
