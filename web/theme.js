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
            updateTerminalTheme(savedTheme);
        }

        // Toggle theme
        toggleButton.addEventListener("click", () => {
            body.classList.toggle("light-theme");
            const newTheme = body.classList.contains("light-theme") ? "light-theme" : "dark-theme";
            updateTerminalTheme(newTheme);
            localStorage.setItem("theme", newTheme);
        });

        function updateTerminalTheme(theme) {
            if (!window.termInitialized || !window.term) {
                console.error("Terminal is not initialized.");
                return;
            }

            const term = window.term;
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

            // Refresh terminal styles
            term.refresh(0, term.rows - 1);
        }
    })();

    // Mode Toggle Functionality
    const modeToggle = document.getElementById("mode-toggle");
    const modeLabel = document.getElementById("mode-label");
    const toggleSwitch = document.querySelector(".toggle-switch");
    const terminal = document.getElementById("terminal");
    const uiCanvas = document.getElementById("ui-canvas");

    if (!modeToggle) {
        console.error("Mode toggle button not found!");
        return;
    }

    modeToggle.addEventListener("click", () => {
        const isCLI = modeLabel.textContent === "CLI";
        toggleSwitch.classList.toggle("active", !isCLI);

        if (isCLI) {
            modeLabel.textContent = "UI";
            terminal.style.transition = "transform 0.5s, opacity 0.5s";
            terminal.style.transform = "translateY(100%)";
            terminal.style.opacity = "0";

            setTimeout(() => {
                terminal.style.display = "none";
                uiCanvas.style.display = "flex";
                uiCanvas.style.opacity = "1";
            }, 500);
        } else {
            modeLabel.textContent = "CLI";
            uiCanvas.style.opacity = "0";

            setTimeout(() => {
                uiCanvas.style.display = "none";
                terminal.style.display = "block";
                terminal.style.transform = "translateY(0)";
                terminal.style.opacity = "1";
            }, 500);
        }
    });
});
