// Check if the worker is already initialized
if (!window.wasmWorker) {
    // Create a new worker and store it globally
    window.wasmWorker = new Worker("scripts/wasm_worker.js");
  }
  