(function () {
    const toggleButton = document.getElementById("theme-toggle");
    const body = document.body;

    // Apply saved theme
    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
        body.classList.add(savedTheme);
        updateToggleIcon(savedTheme);
        updateTerminalTheme(savedTheme);
    }

    // Toggle theme
    toggleButton.addEventListener("click", () => {
        body.classList.toggle("light-theme");
        const newTheme = body.classList.contains("light-theme") ? "light-theme" : "dark-theme";
        updateToggleIcon(newTheme);
        localStorage.setItem("theme", newTheme);
        updateTerminalTheme(newTheme);
    });

    function updateToggleIcon(theme) {
        const icon = toggleButton.querySelector("i");
        // Always use the fa-circle-half-stroke icon for theme toggle
        if (theme === "light-theme") {
            icon.classList.add("light-icon");
            icon.classList.remove("dark-icon");
        } else {
            icon.classList.add("dark-icon");
            icon.classList.remove("light-icon");
        }
    }

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
    }
})();

document.addEventListener("DOMContentLoaded", () => {
    const modeToggle = document.getElementById("mode-toggle");
    const modeLabel = document.getElementById("mode-label");
    const modeIcon = modeToggle.querySelector("i");
    const terminal = document.getElementById("terminal");
    const uiCanvas = document.getElementById("ui-canvas");

    if (!modeToggle) {
        console.error("Toggle button not found!");
        return;
    }

    modeToggle.addEventListener("click", () => {
        const isCLI = modeLabel.textContent === "CLI";
        if (isCLI) {
            // Switch to UI mode
            modeLabel.textContent = "UI";
            modeIcon.classList.replace("fa-toggle-on", "fa-toggle-off");
            terminal.style.transition = "transform 0.5s, opacity 0.5s";
            terminal.style.transform = "translateY(100%)";
            terminal.style.opacity = "0";

            setTimeout(() => {
                terminal.style.display = "none";
                uiCanvas.style.display = "flex";
                uiCanvas.style.opacity = "1";
                uiCanvas.style.transform = "translateY(0)";
            }, 500);
        } else {
            // Switch to CLI mode
            modeLabel.textContent = "CLI";
            modeIcon.classList.replace("fa-toggle-off", "fa-toggle-on");
            uiCanvas.style.transition = "transform 0.5s, opacity 0.5s";
            uiCanvas.style.opacity = "0";
            uiCanvas.style.transform = "translateY(100%)";

            setTimeout(() => {
                uiCanvas.style.display = "none";
                terminal.style.display = "block";
                terminal.style.transform = "translateY(0)";
                terminal.style.opacity = "1";
            }, 500);
        }
    });
});

