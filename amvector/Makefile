VECTOR 			= amvector
VECTOR_ARM      = amvector_arm
VECTOR_SERVER   = ./cmd/${VECTOR}

.PHONY: bin
# build bin
bin:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o $(VECTOR) $(VECTOR_SERVER)

.PHONY: arm
# build arm
arm:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o $(VECTOR_ARM) $(VECTOR_SERVER)

.PHONY: wire
# generate wire
wire:
	wire ./service/
