// Load the wasm_exec.js file into the worker context
importScripts("wasm_exec.js");

let wasmLoaded = false;
let executeCommand = null;

// Flag to ensure worker initializes only once
let wasmInitialized = false;

async function initializeWasm() {
    if (!wasmInitialized) {
        if (!WebAssembly.instantiateStreaming) {
            WebAssembly.instantiateStreaming = async (resp, importObject) => {
                const source = await (await resp).arrayBuffer();
                return await WebAssembly.instantiate(source, importObject);
            };
        }

        const go = new Go(); // Go runtime is now defined because wasm_exec.js is loaded
        const result = await WebAssembly.instantiateStreaming(fetch("../assets/app.wasm"), go.importObject);
        go.run(result.instance);

        executeCommand = self.executeCommand; // Cache the executeCommand function
        wasmLoaded = true; // Mark WASM as loaded
        wasmInitialized = true; // Set the flag to true to prevent reinitialization
        self.postMessage({ status: "wasm-ready" }); // Notify the main thread
    }
}

// Automatically initialize the worker as soon as it's loaded (only once)
initializeWasm();

// Listen for incoming messages from the main thread
self.onmessage = async (event) => {
    // Include correlationId in the destructuring
    const { command, type, correlationId } = event.data;

    // Handle explicit initialization requests separately if needed
    if (type === "initialize") {
        if (!wasmLoaded) {
            await initializeWasm();
        }
        self.postMessage({ status: "wasm-ready", correlationId });
        return;
    }

    if (!wasmLoaded) {
        await initializeWasm();
    }

    if (typeof executeCommand === "function") {
        try {
            if (typeof command === "string") {
                const output = executeCommand(command.trim());
                // Include correlationId in the response
                self.postMessage({ output, correlationId });
            } else {
                self.postMessage({ error: "Command is undefined or not a string.", correlationId });
            }
        } catch (err) {
            self.postMessage({ error: `Error executing command: ${err.message}`, correlationId });
        }
    } else {
        self.postMessage({ error: "executeCommand is not available.", correlationId });
    }
};