document.addEventListener("DOMContentLoaded", () => {
    // Theme Toggle Functionality
    (function () {
        const toggleButton = document.getElementById("theme-toggle");
        const body = document.body;

        if (!toggleButton) {
            console.error("Theme toggle button not found!");
            return;
        }

        // Apply saved theme
        const savedTheme = localStorage.getItem("theme");
        if (savedTheme) {
            body.classList.add(savedTheme);
            updateTheme(savedTheme);
        }

        // Toggle theme
        toggleButton.addEventListener("click", () => {
            body.classList.toggle("light-theme");
            const newTheme = body.classList.contains("light-theme") ? "light-theme" : "dark-theme";
            updateTheme(newTheme);
            localStorage.setItem("theme", newTheme);
        });

        function updateTheme(theme) {
            if (!window.termInitialized || !window.term) {
                console.error("Terminal is not initialized.");
                return;
            }

            const term = window.term;
            const uiCanvas = document.getElementById("ui-canvas");

            term.options.theme = theme === "light-theme" ? {
                background: '#ffffff',
                foreground: '#000000',
                cursor: '#000000',
                cursorAccent: '#ffffff',
                selection: '#c7c7c7',
            } : {
                background: '#1e1e1e',
                foreground: '#ffffff',
                cursor: '#ffffff',
                cursorAccent: '#000000',
                selection: '#555555',
            };

            // Update UI Canvas styles
            if (uiCanvas) {
                if (theme === "light-theme") {
                    uiCanvas.style.backgroundColor = "#ffffff";
                    uiCanvas.style.color = "#000000";
                } else {
                    uiCanvas.style.backgroundColor = "#1e1e1e";
                    uiCanvas.style.color = "#ffffff";
                }
            }

            term.refresh(0, term.rows - 1);
        }
    })();

    // Mode Toggle Functionality
    const modeToggle = document.getElementById("mode-toggle");
    const modeLabel = document.getElementById("mode-label");
    const terminal = document.getElementById("terminal");
    const uiCanvas = document.getElementById("ui-canvas");

    if (!modeToggle || !modeLabel || !terminal || !uiCanvas) {
        console.error("One or more elements for mode toggle not found!");
        return;
    }

    modeToggle.addEventListener("click", () => {
        const isCLI = modeLabel.textContent === "CLI";

        if (isCLI) {
            // Switch to UI mode
            modeLabel.textContent = "UI";
            terminal.style.transition = "transform 0.3s, opacity 0.3s";
            terminal.style.transform = "translateY(100%)";
            terminal.style.opacity = "0";

            setTimeout(() => {
                terminal.style.display = "none";
                uiCanvas.style.display = "flex";
                uiCanvas.style.transform = "scale(0.9)";
                uiCanvas.style.opacity = "0";

                setTimeout(() => {
                    uiCanvas.style.transition = "transform 0.3s, opacity 0.3s";
                    uiCanvas.style.transform = "scale(1)";
                    uiCanvas.style.opacity = "1";
                }, 50);
            }, 300);
        } else {
            // Switch to CLI mode
            modeLabel.textContent = "CLI";
            uiCanvas.style.transition = "transform 0.3s, opacity 0.3s";
            uiCanvas.style.transform = "scale(0.9)";
            uiCanvas.style.opacity = "0";

            setTimeout(() => {
                uiCanvas.style.display = "none";
                terminal.style.display = "block";
                terminal.style.transition = "transform 0.3s, opacity 0.3s";
                terminal.style.transform = "translateY(100%)";
                terminal.style.opacity = "0";

                setTimeout(() => {
                    terminal.style.transform = "translateY(0)";
                    terminal.style.opacity = "1";
                }, 50);
            }, 300);
        }
    });
});
