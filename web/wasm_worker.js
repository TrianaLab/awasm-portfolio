let wasmLoaded = false;
let executeCommand = null;

async function initializeWasm() {
    if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject);
    go.run(result.instance);

    executeCommand = self.executeCommand; // Cache the executeCommand function
    wasmLoaded = true; // Mark WASM as loaded
    self.postMessage({ status: "wasm-ready" }); // Notify main thread
}

self.onmessage = async (event) => {
    const { command, type } = event.data;

    if (type === "initialize") {
        if (!wasmLoaded) {
            await initializeWasm();
        }
        return;
    }

    if (!wasmLoaded) {
        await initializeWasm();
    }

    if (typeof executeCommand === "function") {
        try {
            const output = executeCommand(command.trim());
            self.postMessage({ output });
        } catch (err) {
            self.postMessage({ error: `Error executing command: ${err.message}` });
        }
    } else {
        self.postMessage({ error: "executeCommand is not available." });
    }
};
