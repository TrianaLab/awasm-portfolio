build:
	rm -f portfolio web/app.wasm           # Remove previous binary and Wasm file
	go mod tidy                           # Ensure dependencies are cleaned up
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js web/
	GOARCH=wasm GOOS=js go build -o web/app.wasm main.go

run: build
	python3 -m http.server 8000 --bind 127.0.0.1 --directory web

