<script module lang="ts">
  // Module-scoped so the auto-typed welcome demo runs exactly once per
  // page load — subsequent terminals opened via the "+" button (or
  // re-mounted after a mode switch) skip the demo.
  let demoHasRun = false;
</script>

<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { Terminal } from '@xterm/xterm';
  import { FitAddon } from '@xterm/addon-fit';
  import { runCommand, completeLine } from '../lib/wasm';
  import '@xterm/xterm/css/xterm.css';

  let container: HTMLDivElement;
  let term: Terminal;
  let fit: FitAddon;
  let resizeObserver: ResizeObserver | null = null;

  const PROMPT = '\x1b[32medudiaz.dev\x1b[0m\x1b[34m$\x1b[0m ';
  // Color codes used in WELCOME (kept short for readability):
  //   \x1b[36m cyan (logo)   \x1b[33m yellow (callout)
  //   \x1b[90m grey (hint)   \x1b[0m  reset
  const WELCOME = [
    '\x1b[36m  __ ___      ____ _ ___ _ __ ___    \x1b[0m',
    '\x1b[36m / _` \\ \\ /\\ / / _` / __| \'_ ` _ \\   \x1b[0m',
    '\x1b[36m| (_| |\\ V  V / (_| \\__ \\ | | | | |  \x1b[0m',
    '\x1b[36m \\__,_| \\_/\\_/ \\__,_|___/_| |_| |_|  \x1b[0m',
    '',
    '\x1b[90m# Welcome. Eduardo\'s resume, exposed as kubectl resources.\x1b[0m',
    '\x1b[90m# Press Tab to discover commands, or just watch…\x1b[0m',
    '',
  ];
  const DEMO_COMMAND = 'kubectl get all';

  let input = '';
  let cursor = 0;
  const history: string[] = [];
  let historyIndex = -1;
  let running = false;
  // Auto-typed demo: when true, character input from the welcome demo
  // loop is being written into `input`; the first user keystroke aborts
  // the demo via `demoAbort`.
  let demoActive = false;
  let demoAbort: (() => void) | null = null;

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

  // Active completion cycle (zsh/fish-style). When a Tab finds multiple
  // candidates, we list them, substitute the first into the buffer, and
  // keep this state so the next Tab swaps in the next candidate. Any
  // edit that isn't another Tab cancels the cycle. Enter while cycling
  // accepts the current selection and adds a space — it does NOT submit
  // the line; the user can then Tab again at the next position.
  type CycleState = {
    candidates: string[];
    index: number;
    anchorStart: number;
    anchorEnd: number;
  };
  let cycle: CycleState | null = null;
  // True while a completion round-trip to the WASM worker is in flight.
  // Used to ignore further keystrokes (Tab, Enter, printable) so the
  // user can't race past the async boundary and accidentally submit
  // the line they were about to complete.
  let completing = false;

  function currentWord(line: string, pos: number): string {
    let start = pos;
    while (start > 0 && line[start - 1] !== ' ') start -= 1;
    return line.slice(start, pos);
  }

  function applyCycleSelection() {
    if (!cycle) return;
    const candidate = cycle.candidates[cycle.index];
    input = input.slice(0, cycle.anchorStart) + candidate + input.slice(cycle.anchorEnd);
    cycle.anchorEnd = cycle.anchorStart + candidate.length;
    cursor = cycle.anchorEnd;
    redrawInput();
  }

  function cancelCycle() {
    cycle = null;
  }

  async function handleTab() {
    // Cycling: each Tab swaps in the next candidate at the same anchor.
    if (cycle) {
      cycle.index = (cycle.index + 1) % cycle.candidates.length;
      applyCycleSelection();
      return;
    }

    // Cobra-side completion operates on "kubectl <args>", so feed it the
    // full line (or the implicit `kubectl ` if the user hasn't typed it).
    const lineForCompletion = input.length === 0 ? 'kubectl ' : input.slice(0, cursor);
    let candidates: string[];
    completing = true;
    try {
      candidates = await completeLine(lineForCompletion);
    } catch {
      completing = false;
      return;
    } finally {
      completing = false;
    }
    if (candidates.length === 0) return;

    const word = currentWord(input, cursor);
    // Cobra returns candidates with the same prefix as the partial token
    // (it does its own filtering), so each candidate already starts
    // with `word`. Be defensive in case that ever changes.
    const matching = candidates.filter((c) => c.startsWith(word));
    const pool = (matching.length > 0 ? matching : candidates).slice().sort();

    if (pool.length === 1) {
      const insertion = pool[0].slice(word.length) + ' ';
      input = input.slice(0, cursor) + insertion + input.slice(cursor);
      cursor += insertion.length;
      redrawInput();
      return;
    }

    // Multiple candidates → list them once, then enter cycle mode at
    // index 0. Subsequent Tabs advance the cycle; Enter commits.
    term.write('\r\n' + formatCandidates(pool) + '\r\n' + PROMPT + input);
    if (cursor < input.length) {
      term.write(`\x1b[${input.length - cursor}D`);
    }
    cycle = {
      candidates: pool,
      index: 0,
      anchorStart: cursor - word.length,
      anchorEnd: cursor,
    };
    applyCycleSelection();
  }

  function commitCycleAndAdvance() {
    // Accept the currently-selected candidate, add the separating
    // space, drop the cycle, and stay on the prompt (do NOT submit).
    // The user can Tab again at the new position.
    cancelCycle();
    input = input.slice(0, cursor) + ' ' + input.slice(cursor);
    cursor += 1;
    redrawInput();
  }

  // Plays the auto-typed welcome demo: types DEMO_COMMAND character by
  // character with a touch of cadence variation, pauses briefly, then
  // executes it. Resolves once finished or aborted. Aborts on the first
  // user keystroke (handled in handleKey).
  async function runWelcomeDemo() {
    demoActive = true;
    let aborted = false;
    demoAbort = () => {
      aborted = true;
    };
    const sleep = (ms: number) =>
      new Promise<void>((resolve) => setTimeout(resolve, ms));

    for (const ch of DEMO_COMMAND) {
      if (aborted) break;
      input = input.slice(0, cursor) + ch + input.slice(cursor);
      cursor += 1;
      redrawInput();
      await sleep(45 + Math.random() * 55);
    }
    if (!aborted) {
      await sleep(450);
    }
    if (aborted) {
      demoActive = false;
      demoAbort = null;
      return;
    }
    const line = input;
    history.push(line);
    historyIndex = history.length;
    await execute(line);
    demoActive = false;
    demoAbort = null;
    prompt();
  }

  function formatCandidates(items: string[]): string {
    // Two-space gap, wrap at the terminal width (with a 40-col floor
    // so very narrow phone viewports still produce readable rows).
    const lines: string[] = [];
    const wrap = Math.max(40, term?.cols ?? 80);
    let row = '';
    for (const item of items) {
      const sep = row.length === 0 ? '' : '  ';
      if (row.length + sep.length + item.length > wrap) {
        lines.push(row);
        row = item;
      } else {
        row += sep + item;
      }
    }
    if (row.length > 0) lines.push(row);
    return lines.join('\r\n');
  }

  function handleKey({ key, domEvent }: { key: string; domEvent: KeyboardEvent }) {
    if (running) return;
    // First user keystroke during the welcome demo aborts it: clear
    // the partial demo command, then fall through so the keystroke
    // itself is processed as the first character of the user's own
    // input (otherwise it would be silently dropped, eating the first
    // letter of whatever they meant to type).
    // demoActive is set false synchronously here so subsequent
    // keystrokes (which run before the async demo loop has finished
    // its cleanup) skip this branch instead of re-clearing the input.
    if (demoActive) {
      demoActive = false;
      demoAbort?.();
      demoAbort = null;
      input = '';
      cursor = 0;
      redrawInput();
    }
    // Drop keystrokes while a completion round-trip is pending so the
    // user can't submit the line by hitting Enter before cycle state
    // has had a chance to populate.
    if (completing) {
      domEvent.preventDefault();
      return;
    }
    const code = domEvent.keyCode;
    const ctrl = domEvent.ctrlKey;

    if (code === 9) {
      // Tab — completion (fresh or cycling). xterm's onKey already
      // prevents the default focus-shift; we just consume it here.
      domEvent.preventDefault();
      void handleTab();
      return;
    }

    // Enter while cycling commits the current selection and advances
    // to the next position instead of submitting the line.
    if (code === 13 && cycle) {
      commitCycleAndAdvance();
      return;
    }

    // Any other key while cycling cancels cycle state and falls through
    // to the regular handler (so typing a character mid-cycle behaves
    // normally — the currently-substituted candidate stays in place and
    // the new character is inserted at the cursor).
    if (cycle) cancelCycle();

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

    resizeObserver = new ResizeObserver(() => {
      fit.fit();
      // fit() reflows xterm asynchronously — defer scrollToBottom to the
      // next animation frame so it sees the new buffer geometry. Without
      // this, maximize after long output leaves the prompt off-screen.
      requestAnimationFrame(() => term.scrollToBottom());
    });
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

    // Welcome demo: type out a sample command after a short beat so the
    // user can register the welcome message first. Runs once per page
    // load — subsequent terminals skip it (the module-scoped flag).
    if (!demoHasRun) {
      demoHasRun = true;
      setTimeout(() => {
        void runWelcomeDemo();
      }, 700);
    }

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
