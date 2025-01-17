if (!window.wasmWorker) {
    // Create the worker and store it globally
    const worker = new Worker("scripts/wasm_worker.js");
    window.wasmWorker = worker;
    window.wasmReady = false;  // Initialize the readiness flag

    // Set up the centralized onmessage handler
    worker.onmessage = (event) => {
        const { status } = event.data;

        if (status === "wasm-ready") {
            window.wasmReady = true;
        }

        // Dispatch a custom event with all details (including correlationId if present)
        const customEvent = new CustomEvent("workerMessage", { detail: event.data });
        document.dispatchEvent(customEvent);
    };

    // Request initialization from the worker
    worker.postMessage({ type: "initialize" });
}
