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

    // Añade una variable para rastrear el tipo de petición actual
    let isDownloadRequest = false;

    // Listen for custom events dispatched by the centralized handler
    document.addEventListener("workerMessage", (event) => {
        const { output, error, status, correlationId } = event.detail;

        if (correlationId && correlationId !== instanceCorrelationId) {
            return; // Ignore messages not meant for this script
        }

        if (status === "wasm-ready") {
            wasmReady = true;
            console.log("WebAssembly module is ready.");
        } else if (error) {
            console.error("Error from WASM module:", error);
        } else if (output) {
            try {
                // Maneja la respuesta según el tipo de petición
                if (isDownloadRequest) {
                    console.log("Processing download request...");
                    // Asegúrate de que la respuesta es JSON válido
                    const jsonOutput = JSON.parse(output);
                    // Si es un array con un solo elemento, toma ese elemento
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
                    isDownloadRequest = false; // Resetea el flag
                } else {
                    // Para peticiones de visualización
                    const jsonData = jsyaml.load(output);
                    renderGraph(jsonData);
                }
            } catch (err) {
                console.error("Failed to process output:", err, output);
                isDownloadRequest = false; // Resetea el flag en caso de error
            }
        }
    });

    function fetchYamlData() {
        if (!wasmReady) {
            console.warn("WASM module is not ready yet.");
            return;
        }

        console.log("Fetching YAML data...");
        worker.postMessage({ 
            type: "command", 
            command: "kubectl get all --all-namespaces --output yaml",
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

            fetchYamlData();
            updateGraphSize();
        } else {
            modeLabel.textContent = "CLI";

            graphContainer.style.visibility = "hidden";
            graphContainer.style.opacity = "0";
            terminal.style.visibility = "visible";
            terminal.style.opacity = "1";
        }
    });

    // New resizing functionality
    const updateGraphSize = () => {
        const svg = document.querySelector("#graph-container svg");
        if (svg) {
            svg.setAttribute("width", `${window.innerWidth}px`);
            svg.setAttribute("height", `${window.innerHeight}px`);
        }
    };

    const updateSizes = () => {
        if (terminal.style.visibility === "visible") {
            terminal.style.height = `${window.innerHeight * 0.8}px`;
            terminal.style.width = `${window.innerWidth * 0.8}px`;
        }
        if (graphContainer.style.visibility === "visible") {
            graphContainer.style.height = `${window.innerHeight * 0.9}px`;
            graphContainer.style.width = `${window.innerWidth * 0.9}px`;
            updateGraphSize();
        }
    };

    window.addEventListener("resize", updateSizes);
    updateSizes();

    // Download resume functionality
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
        isDownloadRequest = true; // Marca que esta es una petición de descarga
        worker.postMessage({ 
            type: "command", 
            command: "kubectl get resume eduardo-diaz -o json",
            correlationId: instanceCorrelationId 
        });
    });
});
