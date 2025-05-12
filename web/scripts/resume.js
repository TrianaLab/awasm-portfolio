let resumeComponent = null;

function loadResumeComponent(jsonData = null) {
    // Guardar la función original de console.log
    const originalConsoleLog = console.log;
    
    // Sobrescribir console.log para filtrar los mensajes del componente
    console.log = function() {
        // Si el mensaje viene del json-resume.js, lo ignoramos
        const stack = new Error().stack;
        if (stack && !stack.includes('json-resume.js')) {
            originalConsoleLog.apply(console, arguments);
        }
    };

    if (!resumeComponent) {
        resumeComponent = document.createElement('json-resume');
        
        if (jsonData) {
            resumeComponent.resumejson = jsonData;
        }
        
        resumeComponent.style.cssText = `
            position: absolute;
            top: 80px;
            left: 0;
            width: 100%;
            height: calc(100vh - 80px);
            overflow-y: auto;
            box-sizing: border-box;
            background: transparent;
        `;

        document.body.appendChild(resumeComponent);
    } else if (jsonData) {
        resumeComponent.resumejson = jsonData;
    }

    // Restaurar la función original de console.log después de un breve delay
    setTimeout(() => {
        console.log = originalConsoleLog;
    }, 100);
}

function unloadResumeComponent() {
    if (resumeComponent && resumeComponent.parentNode) {
        resumeComponent.parentNode.removeChild(resumeComponent);
        resumeComponent = null;
    }
}

// Definir el importmap dinámicamente
const importMap = document.createElement('script');
importMap.type = 'importmap';
importMap.textContent = JSON.stringify({
    imports: {
        "lit": "https://esm.sh/lit@2.7.2?bundle",
        "lit/": "https://esm.sh/lit@2.7.2/",
        "@lit/task": "https://esm.sh/@lit/task@1.0.2?bundle",
        "@lit/task/": "https://esm.sh/@lit/task@1.0.2/"
    }
});
document.head.appendChild(importMap);

// Importar el componente json-resume pero no cargarlo automáticamente
import('https://unpkg.com/jsonresume-component');

// Exportar las funciones para usarlas en mode.js
window.resumeUtils = {
    loadResumeComponent,
    unloadResumeComponent
};