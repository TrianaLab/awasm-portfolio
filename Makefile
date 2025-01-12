build:
	rm -f portfolio web/app.wasm           # Remove previous binary and Wasm file
	go mod tidy                            # Ensure dependencies are cleaned up
	GO_VERSION=$$(go version | awk '{print $$3}'); \
	WASM_URL="https://raw.githubusercontent.com/golang/go/$${GO_VERSION}/misc/wasm/wasm_exec.js"; \
	if curl --output /dev/null --silent --head --fail "$${WASM_URL}"; then \
		echo "Downloading wasm_exec.js for $${GO_VERSION}..."; \
		curl -o web/wasm_exec.js "$${WASM_URL}"; \
	else \
		echo "Error: Unable to download wasm_exec.js for $${GO_VERSION}. Please ensure the version is correct."; \
		exit 1; \
	fi
	GOARCH=wasm GOOS=js go build -o web/app.wasm main.go

run: build
	python3 -m http.server 8000 --bind 127.0.0.1 --directory web
