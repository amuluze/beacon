TARGET 			= amvector
SERVER   = ./cmd/${TARGET}

GO         = go
GOBUILD    = $(GO) build -mod=readonly

INSTALLER  = amvector.install
OBJPATH    = amvector-obj
GOBINDATA  = go-bindata
BINDATAPKG = assets
BINDATAPATH = $(BINDATAPKG)/bindata.go

.PHONY: amd64
# build amd64
amd64:
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CGO_LDFLAGS="-static" $(GOBUILD) -ldflags="-s -w" -o $(TARGET) $(SERVER)
	#CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-gnu-g++ CGO_LDFLAGS="-static" go build -ldflags="-s -w" -o $(TARGET) $(SERVER)

.PHONY: arm64
# build arm64
arm64:
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CGO_LDFLAGS="-static" $(GOBUILD) -ldflags="-s -w" -o $(TARGET) $(SERVER)
	#CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++ CGO_LDFLAGS="-static" go build -ldflags="-s -w" -o $(TARGET) $(SERVER)

.PHONY: wire
# generate wire
wire:
	wire ./service/

.PHONY: assets
# generate assets
assets:
	@$(GOBINDATA) --nomemcopy -pkg $(BINDATAPKG) -o $(BINDATAPATH) assets/resources/...


.PHONY: installer
# generate installer
installer: assets
	@mkdir -p $(OBJPATH)
	@cp $(TARGET) $(OBJPATH)
	@strip $(OBJPATH)/$(TARGET)
	@cp ./assets/setup.sh $(OBJPATH)
	@cp ./assets/config.yml $(OBJPATH)
	@makeself --gzip $(OBJPATH) $(INSTALLER) "builder" ./setup.sh
	@rm -r $(OBJPATH)
