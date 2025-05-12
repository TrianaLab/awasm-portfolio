document.addEventListener("DOMContentLoaded", () => {
    const modeToggle = document.getElementById("mode-toggle");
    const modeLabel = document.getElementById("mode-label");
    const terminal = document.getElementById("terminal");

    if (!modeToggle || !modeLabel || !terminal) {
        console.error("Required elements not found!");
        return;
    }

    // Configuración inicial: CLI visible
    terminal.style.visibility = "visible";
    terminal.style.opacity = "1";

    const worker = window.wasmWorker;
    if (!worker) {
        console.error("WebAssembly worker is not available.");
        return;
    }

    let wasmReady = false;
    const instanceCorrelationId = "mode-" + Math.random().toString(36).substr(2, 9);
    let isDownloadRequest = false;

    modeToggle.addEventListener("click", () => {
        const isUI = modeLabel.textContent === "UI";

        if (isUI) {
            // Cambiamos a modo UI, por lo que el botón mostrará CLI
            modeLabel.textContent = "CLI";
            terminal.style.visibility = "hidden";
            terminal.style.opacity = "0";
            window.resumeUtils.loadResumeComponent();
        } else {
            // Cambiamos a modo CLI, por lo que el botón mostrará UI
            modeLabel.textContent = "UI";
            terminal.style.visibility = "visible";
            terminal.style.opacity = "1";
            window.resumeUtils.unloadResumeComponent();
        }
    });

    let jsonResume = document.querySelector("json-resume");

    if (!jsonResume) {
        console.error("Required elements not found!");
        return;
    }

    // Configuración inicial: UI (json-resume) oculto
    jsonResume.style.visibility = "hidden";
    jsonResume.style.opacity = "0";

    document.addEventListener("workerMessage", (event) => {
        const { output, error, status, correlationId } = event.detail;

        if (correlationId && correlationId !== instanceCorrelationId) {
            return;
        }

        if (status === "wasm-ready") {
            wasmReady = true;
            console.log("WebAssembly module is ready.");
        } else if (error) {
            console.error("Error from WASM module:", error);
        } else if (output) {
            try {
                if (isDownloadRequest) {
                    console.log("Processing download request...");
                    const jsonOutput = JSON.parse(output);
                    const resumeData = Array.isArray(jsonOutput) && jsonOutput.length === 1 
                        ? jsonOutput[0] 
                        : jsonOutput;
                    
                    const blob = new Blob([JSON.stringify(resumeData, null, 2)], { 
                        type: 'application/json' 
                    });
                    const url = window.URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = 'resume.json';
                    document.body.appendChild(a);
                    a.click();
                    window.URL.revokeObjectURL(url);
                    document.body.removeChild(a);
                    isDownloadRequest = false;

                    fetchJsonData();
                } else {
                    const jsonData = jsyaml.load(output);
                }
            } catch (err) {
                console.error("Failed to process output:", err, output);
                isDownloadRequest = false;
            }
        }
    });

    function fetchJsonData() {
        if (!wasmReady) {
            console.warn("WASM module is not ready yet.");
            return;
        }

        console.log("Fetching YAML data...");
        worker.postMessage({ 
            type: "command", 
            command: "kubectl get profile eduardo-diaz --output json",
            correlationId: instanceCorrelationId 
        });
    }

    const downloadButton = document.getElementById("download-resume");
    if (!downloadButton) {
        console.error("Download button not found!");
        return;
    }

    downloadButton.addEventListener("click", async () => {
        if (!wasmReady) {
            console.warn("WASM module is not ready yet.");
            return;
        }

        console.log("Fetching resume data...");
        isDownloadRequest = true;
        worker.postMessage({ 
            type: "command", 
            command: "kubectl get resume eduardo-diaz -o json",
            correlationId: instanceCorrelationId 
        });
    });
});
