document.addEventListener("DOMContentLoaded", () => {
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
});
