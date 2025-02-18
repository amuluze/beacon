TARGET 	 = amvector
SERVER   = ./cmd/${TARGET}

GO         = go
GOBUILD    = $(GO) build

INSTALLER   = amvector.install
OBJPATH     = amvector-obj
GOBINDATA   = go-bindata
BINDATAPKG  = assets
BINDATAPATH = $(BINDATAPKG)/bindata.go

.PHONY: assets
# generate assets
assets:
	@$(GOBINDATA) --nomemcopy -pkg $(BINDATAPKG) -o $(BINDATAPATH) assets/resources/...

.PHONY: wire
# generate wire
wire: assets
	wire ./service/

.PHONY: amd64
# build amd64
amd64: assets
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(TARGET) $(SERVER)

.PHONY: arm64
# build arm64
arm64: wire assets
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(TARGET) $(SERVER)

.PHONY: installer
# generate installer
installer:
	@mkdir -p $(OBJPATH)
	@cp $(TARGET) $(OBJPATH)
	@strip $(OBJPATH)/$(TARGET)
	@cp ./assets/setup.sh $(OBJPATH)
	@cp ./assets/config.yml $(OBJPATH)
	@makeself --gzip $(OBJPATH) $(INSTALLER) "amvector" ./setup.sh
	@rm -r $(OBJPATH)
	@rm $(TARGET)
