let resumeComponent = null;

function loadResumeComponent(jsonData = null) {
    const originalConsoleLog = console.log;
    
    console.log = function() {
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
            height: calc(100dvh - 80px);
            overflow-y: auto;
            box-sizing: border-box;
            background: transparent;
        `;

        document.body.appendChild(resumeComponent);
    } else if (jsonData) {
        resumeComponent.resumejson = jsonData;
    }

    setTimeout(() => {
        console.log = originalConsoleLog;
    }, 100);
}

function unloadResumeComponent() {
    resumeComponent?.parentNode?.removeChild(resumeComponent);
    resumeComponent = null;
}

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

import('https://unpkg.com/jsonresume-component');

window.resumeUtils = {
    loadResumeComponent,
    unloadResumeComponent
};