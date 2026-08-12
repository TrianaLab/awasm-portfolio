<script lang="ts">
  import Desktop from './Desktop.svelte';
  import { TERMINAL_SUGGESTIONS } from '../lib/portfolio';

  let desktop = $state<Desktop | undefined>();
  let copied = $state<string | null>(null);

  async function copy(command: string) {
    try {
      await navigator.clipboard.writeText(command);
      copied = command;
      setTimeout(() => {
        if (copied === command) copied = null;
      }, 1600);
    } catch {
      /* clipboard blocked — the command is still selectable as text */
    }
  }
</script>

<div class="pane">
  <!-- The visible content of this view is a terminal, so there is nowhere
       natural to hang a page heading; without this the document has no h1 at
       all and its outline starts at level 2. -->
  <h1 class="sr-only">Terminal</h1>

  <div class="bar">
    <p class="lede">
      <span class="mono prompt" aria-hidden="true">$</span>
      A Go CLI compiled to WebAssembly, serving the résumé as <code>kubectl</code>-style
      resources. Nothing is pre-recorded: the completion and the error messages come from the
      binary itself.
    </p>
    <button type="button" class="btn" onclick={() => desktop?.open()}>New terminal</button>
  </div>

  <div class="suggestions">
    <h2 class="sr-only">Suggested commands</h2>
    <ul>
      {#each TERMINAL_SUGGESTIONS as s (s.command)}
        <li>
          <button
            type="button"
            class="chip"
            onclick={() => copy(s.command)}
            aria-label="Copy command: {s.command}, {s.hint}"
            title={s.hint}
          >
            <code>{s.command}</code>
            <span class="chip-state" aria-hidden="true">{copied === s.command ? 'copied' : 'copy'}</span>
          </button>
        </li>
      {/each}
    </ul>
    <!-- The chip's own label is static and its 'copied' badge is decorative,
         so this is the only confirmation a screen-reader user gets. -->
    <p class="sr-only" role="status">{copied ? `Copied: ${copied}` : ''}</p>
  </div>

  <Desktop bind:this={desktop} />
</div>

<style>
  .pane {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }

  .bar {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem 1.5rem;
    padding: 0.6rem var(--page-pad);
    border-bottom: 1px solid var(--color-border);
  }
  .lede {
    margin: 0;
    max-width: 72ch;
    font-size: 0.85rem;
    color: var(--color-text-muted);
  }
  .prompt {
    color: var(--color-success);
    margin-right: 0.35rem;
  }
  .lede code {
    font-family: var(--font-mono);
    font-size: 0.8rem;
    color: var(--color-text);
  }

  .suggestions {
    padding: 0.5rem var(--page-pad);
    border-bottom: 1px solid var(--color-border);
    overflow-x: auto;
  }
  .suggestions ul {
    display: flex;
    gap: 0.4rem;
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .chip {
    display: inline-flex;
    align-items: center;
    gap: 0.55rem;
    min-height: var(--tap);
    padding: 0 0.7rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-elevated);
    color: var(--color-text);
    white-space: nowrap;
    transition: border-color var(--transition);
  }
  .chip:hover {
    border-color: var(--color-accent);
  }
  .chip code {
    font-family: var(--font-mono);
    font-size: 0.78rem;
  }
  .chip-state {
    font-size: 0.65rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }

  @media (max-width: 780px) {
    .bar,
    .suggestions {
      padding-left: 1rem;
      padding-right: 1rem;
    }
  }

  /* This view pins html/body to the viewport, so the chrome above the desktop
     is space the terminal window cannot get back by scrolling. At 200% zoom on
     a laptop the lede and the chip row leave less than the window's 220px
     minimum and its bottom is clipped off unreachably, so they stand down.
     Both are redundant there: the same suggestions are printed in the terminal
     welcome banner. */
  @media (max-height: 560px) {
    .lede,
    .suggestions {
      display: none;
    }
  }
</style>
