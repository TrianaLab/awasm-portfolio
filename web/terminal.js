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

    let executeCommand = null;

    // Polyfill for older browsers
    if (!WebAssembly.instantiateStreaming) {
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    // Load the WebAssembly module
    const wasmLoaded = new Promise((resolve, reject) => {
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
            .then((result) => {
                go.run(result.instance);
                executeCommand = window.executeCommand; // Set global function
                if (typeof executeCommand === "function") {
                    resolve();
                } else {
                    reject(new Error("executeCommand not found after WASM initialization"));
                }
            })
            .catch(reject);
    });

    function writePrompt() {
        term.write("$ ");
    }

    writePrompt();

    async function processCommand(command) {
        if (command.trim() === "clear") {
            term.clear();
            writePrompt();
            currentInput = "";
            cursorPosition = 0;
            return;
        }

        await wasmLoaded; // Ensure WASM is loaded
        if (typeof executeCommand === "function") {
            const output = executeCommand(command.trim());
            term.write(output.replace(/\n/g, "\r\n") + "\r\n");
        } else {
            term.write("Error: executeCommand is not available.\r\n");
        }

        writePrompt();
        currentInput = "";
        cursorPosition = 0;
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
