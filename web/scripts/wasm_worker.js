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
    const { command, type } = event.data;

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
