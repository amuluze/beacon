VECTOR 			= amvector
VECTOR_ARM 		= amvector_arm
VECTOR_SERVER   = ./cmd/${VECTOR}

.PHONY: amd64
# build bin
amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(VECTOR) $(VECTOR_SERVER)

.PHONY: arm64
# build arm
arm64:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o $(VECTOR_ARM) $(VECTOR_SERVER)

.PHONY: wire
# generate wire
wire:
	wire ./service/
