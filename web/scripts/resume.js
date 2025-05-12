let resumeComponent = null;

function loadResumeComponent(jsonData = null) {
    if (!resumeComponent) {
        resumeComponent = document.createElement('json-resume');
        
        // Si tenemos datos JSON, los usamos directamente
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
        // Actualizar datos si el componente ya existe
        resumeComponent.resumejson = jsonData;
    }
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