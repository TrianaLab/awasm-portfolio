document.addEventListener("DOMContentLoaded", () => {
    const modeToggle = document.getElementById("mode-toggle");
    const modeLabel = document.getElementById("mode-label");
    const terminal = document.getElementById("terminal");
    const uiCanvas = document.getElementById("ui-canvas");

    if (!modeToggle || !modeLabel || !terminal || !uiCanvas) {
        console.error("One or more elements for mode toggle not found!");
        return;
    }

    // Initialize the worker
    const worker = new Worker("scripts/wasm_worker.js");
    let wasmReady = false;

    worker.onmessage = (event) => {
        const { output, error, status } = event.data;

        if (status === "wasm-ready") {
            wasmReady = true;
            console.log("WebAssembly module is ready.");
        }

        if (error) {
            console.error("Error from WASM module:", error);
        } else if (output) {
            console.log("Response from WASM module:", output);
        }
    };

    worker.postMessage({ type: "initialize" }); // Initialize the WASM module

    // Function to call the WebAssembly module
    function callWasmCommand(command) {
        if (!wasmReady) {
            console.warn("WASM module is not ready yet.");
            return;
        }
        worker.postMessage({ type: "command", command });
    }

    // Set initial states
    uiCanvas.style.transform = "translateY(100%)";
    uiCanvas.style.opacity = "0";
    uiCanvas.style.visibility = "hidden";

    modeToggle.addEventListener("click", () => {
        const isCLI = modeLabel.textContent === "CLI";

        if (isCLI) {
            // Switch to UI mode
            modeLabel.textContent = "UI";

            terminal.style.transition = "transform 0.3s ease-in-out, opacity 0.3s";
            terminal.style.transform = "translateY(-100%)";
            terminal.style.opacity = "0";

            setTimeout(() => {
                terminal.style.visibility = "hidden";
                uiCanvas.style.visibility = "visible";
                uiCanvas.style.display = "flex";
                uiCanvas.style.transition = "transform 0.3s ease-in-out, opacity 0.3s";
                uiCanvas.style.transform = "translateY(0)";
                uiCanvas.style.opacity = "1";

                // Call the WebAssembly module and log the response
                console.log("Switching to UI mode...");
                callWasmCommand("kubectl get all --all-namespaces --output json");
            }, 300); // Match animation duration
        } else {
            // Switch to CLI mode
            modeLabel.textContent = "CLI";

            uiCanvas.style.transition = "transform 0.3s ease-in-out, opacity 0.3s";
            uiCanvas.style.transform = "translateY(100%)";
            uiCanvas.style.opacity = "0";

            setTimeout(() => {
                uiCanvas.style.visibility = "hidden";
                terminal.style.visibility = "visible";
                terminal.style.display = "block";
                terminal.style.transition = "transform 0.3s ease-in-out, opacity 0.3s";
                terminal.style.transform = "translateY(0)";
                terminal.style.opacity = "1";
            }, 300); // Match animation duration
        }
    });
});
