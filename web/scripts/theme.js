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
                if (theme === "light-theme") {
                    term.options.theme = {
                        background: '#ffffff',
                        foreground: '#000000',
                        cursor: '#000000',
                        cursorAccent: '#ffffff',
                        selectionBackground: 'rgba(0, 100, 0, 0.5)',    // darker green for better contrast
                        selectionForeground: '#ffffff',                // white text on dark selection
                    };
                } else {
                    term.options.theme = {
                        background: '#1e1e1e',
                        foreground: '#ffffff',
                        cursor: '#ffffff',
                        cursorAccent: '#000000',
                        selectionBackground: 'rgba(255, 255, 0, 0.5)',  // light yellow for good contrast with white
                        selectionForeground: '#000000',                // black text on yellow selection
                    };
                }
                term.refresh(0, term.rows - 1);
            }

            if (uiCanvas) {
                uiCanvas.style.backgroundColor = theme === "light-theme" ? "#ffffff" : "#1e1e1e";
                uiCanvas.style.color = theme === "light-theme" ? "#000000" : "#ffffff";
            }
        }

        // Añadir observer para actualizar el tema del json-resume
        const observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                if (mutation.target.classList.contains('light-theme')) {
                    document.documentElement.setAttribute('data-theme', 'light');
                } else {
                    document.documentElement.setAttribute('data-theme', 'dark');
                }
            });
        });

        observer.observe(document.body, {
            attributes: true,
            attributeFilter: ['class']
        });
    })();
});
