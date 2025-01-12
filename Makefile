build:
	rm -f portfolio web/app.wasm           # Remove previous binary and Wasm file
	go mod tidy                            # Ensure dependencies are cleaned up
	WASM_EXEC_JS=$$(find $(shell go env GOROOT) -type f -name "wasm_exec.js" 2>/dev/null); \
	if [ -z "$$WASM_EXEC_JS" ]; then \
		echo "Error: wasm_exec.js not found. Please ensure your Go installation includes it."; \
		exit 1; \
	fi; \
	cp $$WASM_EXEC_JS web/
	GOARCH=wasm GOOS=js go build -o web/app.wasm main.go

run: build
	python3 -m http.server 8000 --bind 127.0.0.1 --directory web
