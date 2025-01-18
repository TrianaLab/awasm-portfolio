document.addEventListener("DOMContentLoaded", () => {
    const modeToggle = document.getElementById("mode-toggle");
    const modeLabel = document.getElementById("mode-label");
    const terminal = document.getElementById("terminal");
    const graphContainer = document.getElementById("graph-container");

    if (!modeToggle || !modeLabel || !terminal || !graphContainer) {
        console.error("Required elements not found!");
        return;
    }

    // Use the existing global worker
    const worker = window.wasmWorker;
    if (!worker) {
        console.error("WebAssembly worker is not available.");
        return;
    }

    let wasmReady = false;

    // Generate a unique correlation ID for this instance of mode.js
    const instanceCorrelationId = "mode-" + Math.random().toString(36).substr(2, 9);

    // Listen for custom events dispatched by the centralized handler
    document.addEventListener("workerMessage", (event) => {
        const { output, error, status, correlationId } = event.detail;

        // Filter messages to only process those matching this script's correlationId,
        // or handle global messages like wasm-ready without a correlationId.
        if (correlationId && correlationId !== instanceCorrelationId) {
            return; // Ignore messages not meant for this script
        }

        if (status === "wasm-ready") {
            wasmReady = true;
            console.log("WebAssembly module is ready.");
        } else if (error) {
            console.error("Error from WASM module:", error);
        } else if (output) {
            // Basic check for JSON-like output
            const trimmedOutput = output.trim();
            if ((trimmedOutput.startsWith("{") && trimmedOutput.endsWith("}")) ||
                (trimmedOutput.startsWith("[") && trimmedOutput.endsWith("]"))) {
                try {
                    const jsonData = JSON.parse(output);
                    console.log("JSON data refreshed:", jsonData);
    
                    // Render the graph with the new data
                    renderGraph(jsonData);
                } catch (err) {
                    console.error("Failed to parse JSON data:", err);
                }
            } else {
                console.warn("Received non-JSON output:", output);
                // Optionally handle non-JSON output here.
            }
        }
    });

    function fetchJsonData() {
        if (!wasmReady) {
            console.warn("WASM module is not ready yet.");
            return;
        }

        console.log("Fetching JSON data...");
        // Send a message to the worker with the correlationId
        worker.postMessage({ 
            type: "command", 
            command: "kubectl get all --all-namespaces --output json",
            correlationId: instanceCorrelationId 
        });
    }

    function renderGraph(jsonData) {
        const renderEvent = new CustomEvent("render-graph", { detail: jsonData });
        document.dispatchEvent(renderEvent);
    }

    modeToggle.addEventListener("click", () => {
        const isCLI = modeLabel.textContent === "CLI";

        if (isCLI) {
            modeLabel.textContent = "UI";

            terminal.style.visibility = "hidden";
            terminal.style.opacity = "0";
            graphContainer.style.visibility = "visible";
            graphContainer.style.opacity = "1";

            // Fetch fresh data and render the graph
            fetchJsonData();
        } else {
            modeLabel.textContent = "CLI";

            graphContainer.style.visibility = "hidden";
            graphContainer.style.opacity = "0";
            terminal.style.visibility = "visible";
            terminal.style.opacity = "1";
        }
    });
});
