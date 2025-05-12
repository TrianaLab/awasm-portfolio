document.addEventListener("DOMContentLoaded", () => {
    const modeToggle = document.getElementById("mode-toggle");
    const modeLabel = document.getElementById("mode-label");
    const terminal = document.getElementById("terminal");
    const downloadButton = document.getElementById("download-resume");

    if (!modeToggle || !modeLabel || !terminal || !downloadButton) return;

    terminal.style.visibility = "visible";
    terminal.style.opacity = "1";

    const worker = window.wasmWorker;
    if (!worker) return;

    let wasmReady = false;
    const correlationId = `mode-${Math.random().toString(36).substr(2, 9)}`;
    let isDownloadRequest = false;

    modeToggle.addEventListener("click", () => {
        if (modeLabel.textContent === "UI") {
            modeLabel.textContent = "CLI";
            terminal.style.visibility = "hidden";
            fetchResumeData();
        } else {
            modeLabel.textContent = "UI";
            terminal.style.visibility = "visible";
            terminal.style.opacity = "1";
            window.resumeUtils.unloadResumeComponent();
        }
    });

    function fetchResumeData() {
        if (!wasmReady) return;
        worker.postMessage({ 
            type: "command",
            command: "kubectl get resume eduardo-diaz --output json",
            correlationId 
        });
    }

    downloadButton.addEventListener("click", () => {
        if (!wasmReady) return;
        isDownloadRequest = true;
        worker.postMessage({ 
            type: "command",
            command: "kubectl get resume eduardo-diaz -o json",
            correlationId 
        });
    });

    document.addEventListener("workerMessage", event => {
        const { output, error, status, correlationId: eventId } = event.detail;
        if (eventId && eventId !== correlationId) return;

        if (status === "wasm-ready") {
            wasmReady = true;
        } else if (error) {
            console.error("WASM module error:", error);
        } else if (output) {
            try {
                if (isDownloadRequest) {
                    const jsonData = JSON.parse(output);
                    const resumeData = Array.isArray(jsonData) && jsonData.length === 1 ? jsonData[0] : jsonData;
                    const blob = new Blob([JSON.stringify(resumeData, null, 2)], { type: 'application/json' });
                    const url = window.URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = 'resume.json';
                    document.body.appendChild(a);
                    a.click();
                    window.URL.revokeObjectURL(url);
                    document.body.removeChild(a);
                    isDownloadRequest = false;
                } else {
                    const jsonData = JSON.parse(output);
                    const resumeData = Array.isArray(jsonData) && jsonData.length === 1 ? jsonData[0] : jsonData;
                    window.resumeUtils.loadResumeComponent(resumeData);
                }
            } catch (err) {
                console.error("Processing output failed:", err, output);
                isDownloadRequest = false;
            }
        }
    });
});
