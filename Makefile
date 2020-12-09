GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get -u
GO_RUN=$(GO_CMD) run
GO_VET=$(GO_CMD) vet
GO_GENERATE=$(GO_CMD) generate

BIN_DIR=./bin

ALLFILE=./...
MAIN_GO=$(APP_DIR)/main.go

build: clean-bin
	$(GO_BUILD) -o $(BIN_DIR)/app $(MAIN_GO)

run:
	$(GO_RUN) $(MAIN_GO)

clean: clean-bin
	$(GO_CLEAN) -cache -modcache -i -r

clean-bin:
	rm -rf $(BIN_DIR)

test:
	$(GO_TEST) $(ALLFILE)

install-tools:
	$(GO_GET) \
	github.com/kisielk/errcheck@latest

static-check:
	$(GO_VET) $(ALLFIILE); errcheck $(ALLFILE)

generate:
	$(GO_GENERATE) $(ALLFILE)
