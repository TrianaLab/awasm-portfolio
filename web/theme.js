(function () {
    const toggleButton = document.getElementById("theme-toggle");
    const body = document.body;

    // Apply saved theme
    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
        body.classList.add(savedTheme);
        toggleButton.textContent = savedTheme === "light-theme" ? "☀️" : "🌙";
        updateTerminalTheme(savedTheme);
    }

    // Toggle theme
    toggleButton.addEventListener("click", () => {
        body.classList.toggle("light-theme");
        const newTheme = body.classList.contains("light-theme") ? "light-theme" : "dark-theme";
        toggleButton.textContent = newTheme === "light-theme" ? "☀️" : "🌙";
        localStorage.setItem("theme", newTheme);
        updateTerminalTheme(newTheme);
    });

    function updateTerminalTheme(theme) {
        // Ensure the terminal is initialized
        if (!window.termInitialized || !window.term) {
            console.error("Terminal is not initialized.");
            return;
        }
    
        const term = window.term; // Access the global terminal instance
        if (theme === "light-theme") {
            term.options.theme = {
                background: '#f5f5f5',
                foreground: '#000000',
                cursor: '#000000',
                cursorAccent: '#ffffff',
                selection: '#c7c7c7',
            };
        } else {
            term.options.theme = {
                background: '#1e1e1e',
                foreground: '#ffffff',
                cursor: '#ffffff',
                cursorAccent: '#000000',
                selection: '#555555',
            };
        }
    }    
})();
