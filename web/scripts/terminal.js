(function () {
    // Create a global terminal instance
    const term = new Terminal({
        cursorBlink: true,
        theme: {
            background: '#1e1e1e',
            foreground: '#ffffff',
            selectionBackground: 'rgba(255, 255, 255, 0.3)',
            selectionForeground: '#000000', 
        },
    });
    window.term = term; // Save the terminal instance globally
    term.open(document.getElementById("terminal"));

    const fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);
    fitAddon.fit();

    // Update terminal size on window resize
    const updateTerminalSize = () => {
        fitAddon.fit();
    };
    window.addEventListener('resize', updateTerminalSize);

    window.termInitialized = true;

    let commandHistory = [];
    let historyIndex = -1;
    let currentInput = "";
    let cursorPosition = 0;

    // Use the existing worker from the global window object
    const worker = window.wasmWorker;
    if (!worker) {
        console.error("WebAssembly worker is not available.");
        return;
    }

    let wasmReady = window.wasmReady || false;

    // Generate a unique correlation ID for this instance of terminal.js
    const terminalCorrelationId = "terminal-" + Math.random().toString(36).substr(2, 9);

    // Listen for custom events dispatched by the centralized handler
    document.addEventListener("workerMessage", (event) => {
        const { output, error, status, correlationId } = event.detail;
    
        // Filter messages to only process those matching this script's correlationId
        // Also allow global messages like 'wasm-ready' without correlationId.
        if (correlationId && correlationId !== terminalCorrelationId) {
            return;
        }
    
        if (status === "wasm-ready") {
            // Check if we've already handled wasm-ready
            if (!wasmReady) {
                wasmReady = true;
                term.write("WebAssembly module initialized.\r\n");
                writePrompt();
            }
            // Exit early to avoid processing further parts of this message
            return;
        }
    
        if (error) {
            term.write(`Error: ${error}\r\n`);
        } else if (output) {
            term.write(output.replace(/\n/g, "\r\n") + "\r\n");
        }
    
        writePrompt();
    });
    

    function writePrompt() {
        term.write("\x1b[32m$ \x1b[0m");
    }

    function showWelcomeMessage() {
        const asciiArt = [
            "\x1b[32m   █████╗ ██╗    ██╗ █████╗ ███████╗███╗   ███╗\x1b[0m",
            "\x1b[32m  ██╔══██╗██║    ██║██╔══██╗██╔════╝████╗ ████║\x1b[0m",
            "\x1b[32m  ███████║██║ █╗ ██║███████║███████╗██╔████╔██║\x1b[0m",
            "\x1b[32m  ██╔══██║██║███╗██║██╔══██║╚════██║██║╚██╔╝██║\x1b[0m",
            "\x1b[32m  ██║  ██║╚███╔███╔╝██║  ██║███████║██║ ╚═╝ ██║\x1b[0m",
            "\x1b[32m  ╚═╝  ╚═╝ ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝\x1b[0m"
        ];
    
        const welcomeMessageLines = [
            "Welcome to \x1b[32mEdu Diaz's\x1b[0m AWASM Portfolio!",
            'Type \x1b[32mkubectl describe profile\x1b[0m to get started or \x1b[32mkubectl --help\x1b[0m to explore all features.'
        ];
    
        // Print each line of ASCII art
        asciiArt.forEach(line => term.write(line + "\r\n"));
    
        // Add some spacing
        term.write("\r\n");
    
        // Print each line of the welcome message
        welcomeMessageLines.forEach(line => term.write(line + "\r\n"));
    
        // Add final spacing
        term.write("\r\n");
    }    
    
    showWelcomeMessage();
    writePrompt();

    async function processCommand(command) {
        if (command.trim() === "clear") {
            term.clear();
            showWelcomeMessage();
            writePrompt();
            currentInput = "";
            cursorPosition = 0;
            return;
        }

        if (command.trim().toLowerCase() === "triana") {
            term.write("💃🏻\r\n");
            writePrompt();
            return;
        }

        if (!wasmReady) {
            term.write("Initializing WebAssembly module...\r\n");
            // Wait until wasmReady becomes true using a promise
            await new Promise((resolve) => {
                const onReady = (event) => {
                    const { status } = event.detail;
                    if (status === "wasm-ready") {
                        wasmReady = true;
                        document.removeEventListener("workerMessage", onReady);
                        resolve();
                    }
                };
                document.addEventListener("workerMessage", onReady);
            });
        }

        // Send command to worker with correlationId
        worker.postMessage({ 
            type: "command", 
            command, 
            correlationId: terminalCorrelationId 
        });
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
