GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get -u
GO_RUN=$(GO_CMD) run
GO_VET=$(GO_CMD) vet
GO_GENERATE=$(GO_CMD) generate
.PHONY: build run clean clean-bin test install-tools static-check generate create-iam-role create-s3-buckets delete-s3-buckets build-lambda-function zip-lambda-function create-lambda-function update-lambda-function

BIN_DIR=./bin
APP_DIR=./cmd

ALLFILE=./...
MAIN_GO=$(APP_DIR)

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
	$(GO_VET) $(ALLFILE); errcheck $(ALLFILE)

generate:
	$(GO_GENERATE) $(ALLFILE)

SETUP_DIR=./cmd/setup

create-iam-role:
	$(GO_RUN) $(SETUP_DIR) create-iam-role

create-s3-buckets:
	$(GO_RUN) $(SETUP_DIR) create-s3-buckets

delete-s3-buckets:
	$(GO_RUN) $(SETUP_DIR) delete-s3-buckets

build-lambda-function: clean-bin
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO_BUILD) -tags lambda.norpc -o $(BIN_DIR)/bootstrap $(MAIN_GO)

zip-lambda-function: build-lambda-function
	cd $(BIN_DIR) && zip function.zip bootstrap

create-lambda-function:
	$(GO_RUN) $(SETUP_DIR) create-lambda-function

update-lambda-function:
	$(GO_RUN) $(SETUP_DIR) update-lambda-function
