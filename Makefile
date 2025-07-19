GO=go
BINARY_NAME=filemac
CMD_DIR=./cmd/filemac

.PHONY: all build run clean test

all: build

build:
	$(GO) build -o $(BINARY_NAME) $(CMD_DIR)

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

test:
	$(GO) test ./...

build-catalog:
	$(GO) build -o $(BINARY_NAME) $(CMD_DIR) && ./$(BINARY_NAME) catalog

build-tags:
	$(GO) build -o $(BINARY_NAME) $(CMD_DIR) && ./$(BINARY_NAME) tags

build-utils:
	$(GO) build -o $(BINARY_NAME) $(CMD_DIR) && ./$(BINARY_NAME) utils

run-catalog: build-catalog
	./$(BINARY_NAME) catalog

run-tags: build-tags
	./$(BINARY_NAME) tags

run-utils: build-utils
	./$(BINARY_NAME) utils
