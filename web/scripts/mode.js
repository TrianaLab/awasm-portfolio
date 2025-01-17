document.addEventListener("DOMContentLoaded", () => {
    const modeToggle = document.getElementById("mode-toggle");
    const modeLabel = document.getElementById("mode-label");
    const terminal = document.getElementById("terminal");
    const uiCanvas = document.getElementById("ui-canvas");

    if (!modeToggle || !modeLabel || !terminal || !uiCanvas) {
        console.error("One or more elements for mode toggle not found!");
        return;
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
