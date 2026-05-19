<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { Terminal } from '@xterm/xterm';
  import { FitAddon } from '@xterm/addon-fit';
  import { runCommand } from '../lib/wasm';
  import '@xterm/xterm/css/xterm.css';

  let container: HTMLDivElement;
  let term: Terminal;
  let fit: FitAddon;
  let resizeObserver: ResizeObserver | null = null;

  const PROMPT = '\x1b[32medudiaz.dev\x1b[0m\x1b[34m$\x1b[0m ';
  const WELCOME = [
    '\x1b[36m  __ ___      ____ _ ___ _ __ ___    \x1b[0m',
    '\x1b[36m / _` \\ \\ /\\ / / _` / __| \'_ ` _ \\   \x1b[0m',
    '\x1b[36m| (_| |\\ V  V / (_| \\__ \\ | | | | |  \x1b[0m',
    '\x1b[36m \\__,_| \\_/\\_/ \\__,_|___/_| |_| |_|  \x1b[0m',
    '',
    'Welcome — type \x1b[33mkubectl --help\x1b[0m to get started.',
    '',
  ];

  let input = '';
  let cursor = 0;
  const history: string[] = [];
  let historyIndex = -1;
  let running = false;

  function write(s: string) {
    term.write(s);
  }

  function prompt() {
    input = '';
    cursor = 0;
    term.write(`\r\n${PROMPT}`);
    // Force the viewport to follow the new prompt — xterm doesn't always
    // scroll on its own when output overflows the visible region.
    term.scrollToBottom();
  }

  function redrawInput() {
    // Erase current line after the prompt and reprint with cursor positioning.
    term.write('\r\x1b[K' + PROMPT + input);
    if (cursor < input.length) {
      term.write(`\x1b[${input.length - cursor}D`);
    }
  }

  async function execute(line: string) {
    running = true;
    try {
      const trimmed = line.trim();
      if (trimmed === 'clear') {
        term.reset();
        return;
      }
      if (trimmed === '') return;
      const output = await runCommand(trimmed);
      if (output) {
        term.write('\r\n' + output.replace(/\n/g, '\r\n'));
      }
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err);
      term.write(`\r\n\x1b[31mError: ${msg}\x1b[0m`);
    } finally {
      running = false;
    }
  }

  function handleKey({ key, domEvent }: { key: string; domEvent: KeyboardEvent }) {
    if (running) return;
    const code = domEvent.keyCode;
    const ctrl = domEvent.ctrlKey;

    if (code === 13) {
      // Enter
      const line = input;
      if (line.trim()) {
        history.push(line);
      }
      historyIndex = history.length;
      void (async () => {
        await execute(line);
        prompt();
      })();
    } else if (code === 8) {
      // Backspace
      if (cursor === 0) return;
      input = input.slice(0, cursor - 1) + input.slice(cursor);
      cursor -= 1;
      redrawInput();
    } else if (code === 46) {
      // Delete
      if (cursor === input.length) return;
      input = input.slice(0, cursor) + input.slice(cursor + 1);
      redrawInput();
    } else if (code === 37) {
      // Left arrow
      if (cursor === 0) return;
      cursor -= 1;
      term.write('\x1b[D');
    } else if (code === 39) {
      // Right arrow
      if (cursor === input.length) return;
      cursor += 1;
      term.write('\x1b[C');
    } else if (code === 38) {
      // Up — previous history
      if (history.length === 0) return;
      historyIndex = Math.max(0, historyIndex - 1);
      input = history[historyIndex];
      cursor = input.length;
      redrawInput();
    } else if (code === 40) {
      // Down — next history
      if (historyIndex < history.length - 1) {
        historyIndex += 1;
        input = history[historyIndex];
      } else {
        historyIndex = history.length;
        input = '';
      }
      cursor = input.length;
      redrawInput();
    } else if (code === 36) {
      // Home
      const delta = cursor;
      cursor = 0;
      if (delta > 0) term.write(`\x1b[${delta}D`);
    } else if (code === 35) {
      // End
      const delta = input.length - cursor;
      cursor = input.length;
      if (delta > 0) term.write(`\x1b[${delta}C`);
    } else if (ctrl && key === '\x03') {
      // Ctrl+C
      term.write('^C');
      prompt();
    } else if (key.length === 1 && key.charCodeAt(0) >= 32) {
      // Printable
      input = input.slice(0, cursor) + key + input.slice(cursor);
      cursor += 1;
      redrawInput();
    }
  }

  function handlePaste(text: string) {
    if (running) return;
    // Single-line paste only — strip newlines and treat the rest as input.
    const clean = text.replace(/\r?\n/g, ' ');
    input = input.slice(0, cursor) + clean + input.slice(cursor);
    cursor += clean.length;
    redrawInput();
  }

  onMount(() => {
    term = new Terminal({
      fontFamily: 'JetBrains Mono, SF Mono, Consolas, monospace',
      fontSize: 14,
      cursorBlink: true,
      convertEol: true,
      theme: themeColors(),
    });
    fit = new FitAddon();
    term.loadAddon(fit);
    term.open(container);
    fit.fit();

    resizeObserver = new ResizeObserver(() => fit.fit());
    resizeObserver.observe(container);

    for (const line of WELCOME) {
      term.writeln(line);
    }
    write(PROMPT);

    term.onKey(handleKey);
    term.onData((data) => {
      // Paste comes through onData rather than onKey when length > 1.
      if (data.length > 1 && !data.startsWith('\x1b')) {
        handlePaste(data);
      }
    });

    const updateTheme = () => {
      term.options.theme = themeColors();
    };
    const observer = new MutationObserver(updateTheme);
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] });
    return () => observer.disconnect();
  });

  onDestroy(() => {
    resizeObserver?.disconnect();
    term?.dispose();
  });

  function themeColors() {
    const styles = getComputedStyle(document.documentElement);
    const get = (v: string, fallback: string) => styles.getPropertyValue(v).trim() || fallback;
    return {
      background: get('--color-bg-elevated', '#161b22'),
      foreground: get('--color-text', '#e6edf3'),
      cursor: get('--color-accent', '#58a6ff'),
      selectionBackground: get('--color-accent-soft', 'rgba(88, 166, 255, 0.3)'),
    };
  }
</script>

<div class="terminal-wrap">
  <!-- xterm's height: 100% sizes against the parent's content box, so
       padding inside .terminal lets xterm spill into the padding area.
       Use sibling flex spacers to create visual breathing room INSTEAD of
       padding inside .terminal, which gives xterm a clean integer pixel
       box to size into. -->
  <div class="pad-top"></div>
  <div class="row">
    <div class="pad-side"></div>
    <div class="terminal" bind:this={container}></div>
    <div class="pad-side"></div>
  </div>
  <div class="pad-bottom"></div>
</div>

<style>
  .terminal-wrap {
    flex: 1;
    min-height: 0;
    min-width: 0;
    display: flex;
    flex-direction: column;
    background: var(--color-bg-elevated);
    overflow: hidden;
  }

  .row {
    flex: 1;
    min-height: 0;
    display: flex;
  }

  .terminal {
    flex: 1;
    min-width: 0;
    min-height: 0;
    overflow: hidden;
  }

  .pad-top {
    flex: none;
    height: 0.6rem;
  }
  .pad-bottom {
    flex: none;
    height: 1.4rem;
  }
  .pad-side {
    flex: none;
    width: 0.6rem;
  }

  .terminal :global(.xterm) {
    height: 100%;
  }
</style>
