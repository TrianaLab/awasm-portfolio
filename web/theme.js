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
            const term = window.term;
            const uiCanvas = document.getElementById("ui-canvas");

            if (term) {
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
                term.refresh(0, term.rows - 1);
            }

            if (uiCanvas) {
                uiCanvas.style.backgroundColor = theme === "light-theme" ? "#ffffff" : "#1e1e1e";
                uiCanvas.style.color = theme === "light-theme" ? "#000000" : "#ffffff";
            }
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

            terminal.style.transition = "transform 0.5s ease-in-out, opacity 0.5s";
            terminal.style.transform = "translateY(-100%)";
            terminal.style.opacity = "0";

            setTimeout(() => {
                terminal.style.visibility = "hidden";
                uiCanvas.style.visibility = "visible";
                uiCanvas.style.display = "flex";
                uiCanvas.style.transition = "transform 0.5s ease-in-out, opacity 0.5s";
                uiCanvas.style.transform = "translateY(0)";
                uiCanvas.style.opacity = "1";
            }, 500);
        } else {
            // Switch to CLI mode
            modeLabel.textContent = "CLI";

            uiCanvas.style.transition = "transform 0.5s ease-in-out, opacity 0.5s";
            uiCanvas.style.transform = "translateY(100%)";
            uiCanvas.style.opacity = "0";

            setTimeout(() => {
                uiCanvas.style.visibility = "hidden";
                terminal.style.visibility = "visible";
                terminal.style.display = "block";
                terminal.style.transition = "transform 0.5s ease-in-out, opacity 0.5s";
                terminal.style.transform = "translateY(0)";
                terminal.style.opacity = "1";
            }, 500);
        }
    });
});
