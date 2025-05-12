document.addEventListener("DOMContentLoaded", () => {
    const toggleButton = document.getElementById("theme-toggle");
    const body = document.body;
    if (!toggleButton) return;
    
    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
        body.classList.add(savedTheme);
        updateTheme(savedTheme);
    }
    
    toggleButton.addEventListener("click", () => {
        body.classList.toggle("light-theme");
        const newTheme = body.classList.contains("light-theme") ? "light-theme" : "dark-theme";
        updateTheme(newTheme);
        localStorage.setItem("theme", newTheme);
    });
    
    function updateTheme(theme) {
        const term = window.term;
        if (term) {
            term.options.theme = theme === "light-theme"
              ? {
                  background: '#ffffff',
                  foreground: '#000000',
                  cursor: '#000000',
                  cursorAccent: '#ffffff',
                  selectionBackground: 'rgba(0, 100, 0, 0.5)',
                  selectionForeground: '#ffffff'
                }
              : {
                  background: '#1e1e1e',
                  foreground: '#ffffff',
                  cursor: '#ffffff',
                  cursorAccent: '#000000',
                  selectionBackground: 'rgba(255, 255, 0, 0.5)',
                  selectionForeground: '#000000'
                };
            term.refresh(0, term.rows - 1);
        }
        const resumeElem = document.querySelector("json-resume");
        if (resumeElem) {
            resumeElem.setAttribute('data-theme', theme === "light-theme" ? 'light' : 'dark');
        }
    }
});
