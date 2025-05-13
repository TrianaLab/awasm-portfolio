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

# Version (default to "development" if not provided externally)
VERSION ?= development

# Targets

.PHONY: build run clean test install-go-test-coverage check-coverage update-readme

# General Build
build: clean ensure-deps fetch-wasm-exec
	@echo "Building WebAssembly binary..."
	GOARCH=$(GOARCH) GOOS=$(GOOS) $(GO) build \
	    -ldflags "-X 'awasm-portfolio/cmd.appVersion=$(VERSION)'" \
	    -o $(APP_WASM) main.go
	@echo "Build complete: $(APP_WASM)"

# Run Local Development Server
# Pass VERSION when calling this target to ensure correct version is built.
run: build
	@echo "Starting development server on http://127.0.0.1:8000..."
	$(PYTHON) -m http.server 8000 --bind 127.0.0.1 --directory $(WEB_DIR)

# Run Tests with Custom Coverage
test: check-coverage
	@echo "Tests completed with coverage."

# Install go-test-coverage
install-go-test-coverage:
	@echo "Installing go-test-coverage..."
	$(GO) install github.com/vladopajic/go-test-coverage/v2@latest

# Run Tests and Check Coverage
check-coverage: install-go-test-coverage
	@echo "Running tests and generating coverage report..."
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	$$(go env GOPATH)/bin/go-test-coverage --config=./.testcoverage.yml

# Extract coverage percentage and update README
update-readme: check-coverage
	@echo "Extracting coverage percentage and updating README..."
	coverage=$$(grep -Po '(?<=total:).*?\d+\.\d+%' cover.out | head -1); \
	if [ -z "$$coverage" ]; then \
		echo "Failed to extract coverage percentage"; \
		exit 1; \
	fi; \
	echo "Updating README with coverage: $$coverage"; \
	sed -i '' "s/^Coverage: .*%/Coverage: $$coverage/" README.md

# Clean Up
clean:
	@echo "Cleaning up previous builds..."
	rm -f $(APP_WASM)
	rm -f $(WASM_EXEC_JS)
	rm -f cover.out
	rm -f coverage.html

# Ensure Go Dependencies
ensure-deps:
	@echo "Tidying up Go modules..."
	$(GO) mod tidy

# Fetch wasm_exec.js from Local Go Installation
fetch-wasm-exec:
	@echo "Fetching wasm_exec.js from local Go installation..."
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WASM_EXEC_JS)

# Get resume.json
resume:
	@echo "Generating resume.json..."
	$(GO) run cli.go get resume main-resume -o json 2>&1 | jq '.[0]' > resume.json