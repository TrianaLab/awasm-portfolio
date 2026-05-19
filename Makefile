# Project metadata
MODULE       := github.com/TrianaLab/awasm-portfolio
VERSION      ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT   := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE   := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# Directories
WEB_DIR      := web
ASSETS_DIR   := $(WEB_DIR)/assets
SCRIPTS_DIR  := $(WEB_DIR)/scripts

# Files
APP_WASM     := $(ASSETS_DIR)/app.wasm
WASM_EXEC_JS := $(SCRIPTS_DIR)/wasm_exec.js

# Tools
GO           := go
PYTHON       := python3

# WASM toolchain
GOARCH       := wasm
GOOS         := js

# Docker
IMAGE        := ghcr.io/trianalab/awasm-portfolio

LDFLAGS := -X '$(MODULE)/cmd.appVersion=$(VERSION)'

.PHONY: build run clean test lint coverage resume docker-build docker-run ensure-deps fetch-wasm-exec

build: clean ensure-deps fetch-wasm-exec
	@echo "==> Building WebAssembly binary ($(VERSION))..."
	GOARCH=$(GOARCH) GOOS=$(GOOS) $(GO) build -ldflags "$(LDFLAGS)" -o $(APP_WASM) main.go
	@echo "    built: $(APP_WASM)"

run: build
	@echo "==> Starting dev server at http://127.0.0.1:8000"
	$(PYTHON) -m http.server 8000 --bind 127.0.0.1 --directory $(WEB_DIR)

test:
	go test ./... -v

lint:
	gofmt -s -l $(shell find . -name '*.go' -not -path './web/*')
	go vet ./...

coverage:
	go test $(shell go list ./... | grep -v /tests/) -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | tail -1

resume:
	@echo "==> Generating resume.json..."
	$(GO) run cli.go get resume main-resume -o json | jq '.[0]' > resume.json

docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(IMAGE):$(VERSION) .

docker-run: docker-build
	docker run --rm -p 8000:80 $(IMAGE):$(VERSION)

clean:
	@rm -f $(APP_WASM) $(WASM_EXEC_JS) coverage.out coverage.html cover.out

ensure-deps:
	@$(GO) mod tidy

fetch-wasm-exec:
	@cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(WASM_EXEC_JS)

include ci.mk
