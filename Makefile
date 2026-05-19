# Project metadata
MODULE       := github.com/TrianaLab/awasm-portfolio
VERSION      ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT   := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE   := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# Directories
WEB_DIR      := web
ASSETS_DIR   := $(WEB_DIR)/assets
SCRIPTS_DIR  := $(WEB_DIR)/scripts
FRONTEND_DIR := frontend

# Files
APP_WASM     := $(ASSETS_DIR)/app.wasm
WASM_EXEC_JS := $(SCRIPTS_DIR)/wasm_exec.js

# Tools
GO           := go
NPM          := npm
PYTHON       := python3

# WASM toolchain
GOARCH       := wasm
GOOS         := js

# Docker
IMAGE        := ghcr.io/trianalab/awasm-portfolio

LDFLAGS := -X '$(MODULE)/cmd.appVersion=$(VERSION)'

.PHONY: build wasm ui run dev test lint coverage resume docker-build docker-run ensure-deps fetch-wasm-exec clean

# Composite build: produce the WASM artifact, then bundle the Svelte SPA
# around it. Order matters — Vite is configured with emptyOutDir=false so
# the WASM written by `make wasm` is preserved when Vite writes the rest.
build: wasm ui

# Go → WebAssembly. Produces web/assets/app.wasm + web/scripts/wasm_exec.js.
wasm: clean-wasm ensure-deps fetch-wasm-exec
	@mkdir -p $(ASSETS_DIR) $(SCRIPTS_DIR)
	@echo "==> Building WebAssembly binary ($(VERSION))..."
	GOARCH=$(GOARCH) GOOS=$(GOOS) $(GO) build -ldflags "$(LDFLAGS)" -o $(APP_WASM) main.go
	@echo "    built: $(APP_WASM)"

# Svelte → Vite build. Writes index.html + assets/* into web/ alongside
# the WASM produced above.
ui:
	@echo "==> Building Svelte frontend..."
	cd $(FRONTEND_DIR) && $(NPM) ci && $(NPM) run build

# Local serve: full production build, then a tiny static server.
run: build
	@echo "==> Serving production build at http://127.0.0.1:8000"
	$(PYTHON) -m http.server 8000 --bind 127.0.0.1 --directory $(WEB_DIR)

# Local dev: Vite dev server with HMR. Build the WASM once up front so
# the worker can find it, then let Vite watch the Svelte side.
dev: wasm
	@echo "==> Starting Vite dev server (HMR)..."
	cd $(FRONTEND_DIR) && $(NPM) install && $(NPM) run dev

test:
	$(GO) test ./... -v

lint:
	gofmt -s -l $(shell find . -name '*.go' -not -path './web/*' -not -path './frontend/node_modules/*')
	$(GO) vet ./...

coverage:
	$(GO) test $(shell go list ./... | grep -v /tests/) -coverprofile=coverage.out -covermode=atomic
	$(GO) tool cover -html=coverage.out -o coverage.html
	@$(GO) tool cover -func=coverage.out | tail -1

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

clean: clean-wasm
	@rm -rf $(WEB_DIR) coverage.out coverage.html cover.out

clean-wasm:
	@rm -f $(APP_WASM) $(WASM_EXEC_JS)

ensure-deps:
	@$(GO) mod tidy

fetch-wasm-exec:
	@mkdir -p $(SCRIPTS_DIR)
	@cp "$$($(GO) env GOROOT)/lib/wasm/wasm_exec.js" $(WASM_EXEC_JS)

include ci.mk
