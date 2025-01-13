(function () {
    // Check if the terminal is already initialized
    if (window.term) {
        console.warn("Terminal is already initialized.");
        return;
    }

    // Create a global terminal instance
    const term = new Terminal({
        cursorBlink: true,
        theme: {
            background: '#1e1e1e',
            foreground: '#ffffff',
        },
    });
    window.term = term; // Save the terminal instance globally
    term.open(document.getElementById("terminal"));

    const fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);
    fitAddon.fit();

    // Terminal initialization is complete
    window.termInitialized = true;

    let commandHistory = [];
    let historyIndex = -1;
    let currentInput = "";
    let cursorPosition = 0;

    // Initialize the worker
    const worker = new Worker("wasm_worker.js");
    let wasmReady = false;

    worker.onmessage = (event) => {
        const { output, error, status } = event.data;

        if (status === "wasm-ready") {
            wasmReady = true;
            term.write("WebAssembly module initialized.\r\n");
            writePrompt();
            return;
        }

        if (error) {
            term.write(`Error: ${error}\r\n`);
        } else if (output) {
            term.write(output.replace(/\n/g, "\r\n") + "\r\n");
        }

        writePrompt();
    };

    worker.postMessage({ type: "initialize" }); // Trigger WASM initialization in the worker

    function writePrompt() {
        term.write("$ ");
    }

    function showWelcomeMessage() {
        const welcomeMessage = `
Welcome to TrianaLab AWASM Portfolio! Type "kubectl --help" to get started.
        `;
        term.write(welcomeMessage + "\r\n\r\n");
    }

    showWelcomeMessage();
    writePrompt();

    async function processCommand(command) {
        if (command.trim() === "clear") {
            term.clear();
            writePrompt();
            currentInput = "";
            cursorPosition = 0;
            return;
        }

        if (!wasmReady) {
            term.write("Initializing WebAssembly module...\r\n");
            await new Promise((resolve) => {
                const interval = setInterval(() => {
                    if (wasmReady) {
                        clearInterval(interval);
                        resolve();
                    }
                }, 50);
            });
        }

        worker.postMessage({ type: "command", command });
    }

    function updateInput() {
        term.write("\u001b[2K\r"); // Clear the current line
        writePrompt();
        term.write(currentInput);
        if (cursorPosition < currentInput.length) {
            term.write(`\u001b[${currentInput.length - cursorPosition}D`); // Move cursor back
        }
    }

    term.onKey(function (e) {
        const char = e.key;

        if (char === "\r") { // Enter key
            if (currentInput.trim() === "") {
                term.write("\r\n");
                writePrompt();
            } else {
                if (commandHistory.length === 0 || commandHistory[commandHistory.length - 1] !== currentInput.trim()) {
                    commandHistory.push(currentInput.trim());
                }
                historyIndex = commandHistory.length;
                term.write("\r\n");
                processCommand(currentInput.trim());
            }
            currentInput = "";
            cursorPosition = 0;
        } else if (char === "\u007F") { // Backspace key
            if (cursorPosition > 0) {
                currentInput = currentInput.slice(0, cursorPosition - 1) + currentInput.slice(cursorPosition);
                cursorPosition--;
                updateInput();
            }
        } else if (char === "\u001b[A") { // Up arrow (history)
            if (historyIndex > 0) {
                historyIndex--;
                currentInput = commandHistory[historyIndex];
                cursorPosition = currentInput.length;
                updateInput();
            }
        } else if (char === "\u001b[B") { // Down arrow (history)
            if (historyIndex < commandHistory.length - 1) {
                historyIndex++;
                currentInput = commandHistory[historyIndex];
                cursorPosition = currentInput.length;
                updateInput();
            } else if (historyIndex === commandHistory.length - 1) {
                historyIndex++;
                currentInput = "";
                cursorPosition = 0;
                updateInput();
            }
        } else if (char === "\u001b[D") { // Left arrow
            if (cursorPosition > 0) {
                cursorPosition--;
                term.write("\u001b[D");
            }
        } else if (char === "\u001b[C") { // Right arrow
            if (cursorPosition < currentInput.length) {
                cursorPosition++;
                term.write("\u001b[C");
            }
        } else if (char === '\u0003') { // Ctrl+C
            term.write("^C\r\n"); // Display "^C" and move to a new line
            currentInput = ""; // Clear the current input
            cursorPosition = 0; // Reset cursor position
            writePrompt(); // Display a new prompt
        } else if (char.length === 1) { // Regular character input
            currentInput = currentInput.slice(0, cursorPosition) + char + currentInput.slice(cursorPosition);
            cursorPosition++;
            updateInput();
        }
    });

    term.textarea.addEventListener("paste", function (event) {
        const pasteText = event.clipboardData.getData("text");
        currentInput = currentInput.slice(0, cursorPosition) + pasteText + currentInput.slice(cursorPosition);
        cursorPosition += pasteText.length;
        updateInput();
        event.preventDefault();
    });
})();
