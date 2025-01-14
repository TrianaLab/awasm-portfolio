# Directories
WEB_DIR := web
ASSETS_DIR := $(WEB_DIR)/assets
SCRIPTS_DIR := $(WEB_DIR)/scripts
STYLES_DIR := $(WEB_DIR)/styles

# Files
WASM_EXEC_JS := $(SCRIPTS_DIR)/wasm_exec.js
APP_WASM := $(ASSETS_DIR)/app.wasm

# Commands
GO := go
PYTHON := python3

# Flags
GOARCH := wasm
GOOS := js

# Targets

.PHONY: build-cloudflare-worker build run clean

# Build for Cloudflare Worker
build-cloudflare-worker: clean ensure-deps fetch-wasm-exec
	@echo "Building WebAssembly binary for Cloudflare Worker..."
	GOARCH=$(GOARCH) GOOS=$(GOOS) $(GO) build -o $(APP_WASM) main.go
	@echo "Build complete: $(APP_WASM)"

# General Build
build: clean ensure-deps fetch-wasm-exec
	@echo "Building WebAssembly binary..."
	GOARCH=$(GOARCH) GOOS=$(GOOS) $(GO) build -o $(APP_WASM) main.go
	@echo "Build complete: $(APP_WASM)"

# Run Local Development Server
run: build
	@echo "Starting development server on http://127.0.0.1:8000..."
	$(PYTHON) -m http.server 8000 --bind 127.0.0.1 --directory $(WEB_DIR)

# Clean Up
clean:
	@echo "Cleaning up previous builds..."
	rm -f $(APP_WASM)
	rm -f $(WASM_EXEC_JS)

# Ensure Go Dependencies
ensure-deps:
	@echo "Tidying up Go modules..."
	$(GO) mod tidy

# Fetch wasm_exec.js
fetch-wasm-exec:
	@echo "Fetching wasm_exec.js..."
	GO_VERSION=$$($(GO) version | awk '{print $$3}'); \
	WASM_URL="https://raw.githubusercontent.com/golang/go/$${GO_VERSION}/misc/wasm/wasm_exec.js"; \
	if curl --output /dev/null --silent --head --fail "$${WASM_URL}"; then \
		echo "Downloading wasm_exec.js for $${GO_VERSION}..."; \
		curl -o $(WASM_EXEC_JS) "$${WASM_URL}"; \
	else \
		echo "Error: Unable to download wasm_exec.js for $${GO_VERSION}. Please ensure the version is correct."; \
		exit 1; \
	fi
