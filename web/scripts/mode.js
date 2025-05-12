document.addEventListener("DOMContentLoaded", () => {
    const modeToggle = document.getElementById("mode-toggle");
    const modeLabel = document.getElementById("mode-label");
    const terminal = document.getElementById("terminal");
    const downloadButton = document.getElementById("download-resume");

    if (!modeToggle || !modeLabel || !terminal || !downloadButton) {
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
            // Cambiamos a modo UI
            modeLabel.textContent = "CLI";
            terminal.style.visibility = "hidden";
            terminal.style.opacity = "0";
            // Obtener datos antes de cargar el componente
            fetchResumeData();
        } else {
            // Cambiamos a modo CLI
            modeLabel.textContent = "UI";
            terminal.style.visibility = "visible";
            terminal.style.opacity = "1";
            window.resumeUtils.unloadResumeComponent();
        }
    });

    // Función para obtener los datos del resume
    function fetchResumeData() {
        if (!wasmReady) {
            console.warn("WASM module is not ready yet.");
            return;
        }

        console.log("Fetching resume data...");
        worker.postMessage({ 
            type: "command", 
            command: "kubectl get resume eduardo-diaz --output json",
            correlationId: instanceCorrelationId 
        });
    }

    // Eliminar la verificación del json-resume aquí
    // let jsonResume = document.querySelector("json-resume");
    // if (!jsonResume) {
    //     console.error("Required elements not found!");
    //     return;
    // }

    // Configuración inicial: UI (json-resume) oculto
    let jsonResume = document.querySelector("json-resume");
    if (jsonResume) {
        jsonResume.style.visibility = "hidden";
        jsonResume.style.opacity = "0";
    }

    // Modificar el event listener del worker
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
                    const jsonData = JSON.parse(output);
                    const resumeData = Array.isArray(jsonData) && jsonData.length === 1 
                        ? jsonData[0] 
                        : jsonData;
                    window.resumeUtils.loadResumeComponent(resumeData);
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
            command: "kubectl get resume eduardo-diaz --output json",
            correlationId: instanceCorrelationId 
        });
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
