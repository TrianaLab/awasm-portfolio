<script lang="ts">
  import { onMount } from 'svelte';
  import Window from './Window.svelte';
  import Terminal from './Terminal.svelte';
  import { createWindowManager } from '../lib/windows.svelte';

  const manager = createWindowManager();
  let desktopEl = $state<HTMLDivElement | null>(null);

  onMount(() => {
    // Open one terminal automatically so the user has something to look at.
    if (manager.windows.length === 0) {
      manager.open('kubectl — terminal');
    }
  });

  function openTerminal() {
    manager.open(`kubectl — terminal #${manager.windows.length + 1}`);
  }

  export function open() {
    openTerminal();
  }
</script>

<div class="desktop" bind:this={desktopEl}>
  {#each manager.windows as win (win.id)}
    <Window {win} {manager} {desktopEl}>
      <Terminal />
    </Window>
  {/each}

  <!-- Snap-to-edge preview overlay (visible while a window is being
       dragged within ~24 px of an edge). -->
  {#if manager.snapHint.rect}
    <div
      class="snap-preview"
      aria-hidden="true"
      style:left="{manager.snapHint.rect.x}px"
      style:top="{manager.snapHint.rect.y}px"
      style:width="{manager.snapHint.rect.w}px"
      style:height="{manager.snapHint.rect.h}px"
    ></div>
  {/if}

  <!-- Dock for minimized windows. -->
  {#if manager.windows.some((w) => w.minimized)}
    <div class="dock" role="toolbar" aria-label="Minimized windows">
      {#each manager.windows.filter((w) => w.minimized) as win (win.id)}
        <button class="dock-item" onclick={() => manager.restore(win.id)} title={win.title}>
          <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor" aria-hidden="true">
            <path d="M4 5a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h16a1 1 0 0 0 1-1V6a1 1 0 0 0-1-1H4Zm2 4 3 3-3 3m4 0h6" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <span>{win.title}</span>
        </button>
      {/each}
    </div>
  {/if}

  <!-- Hint when no windows are open. -->
  {#if manager.windows.length === 0}
    <div class="empty">
      <p>No windows open.</p>
      <button class="action" onclick={openTerminal}>Open a terminal</button>
    </div>
  {/if}
</div>

<style>
  .desktop {
    position: relative;
    flex: 1;
    min-height: 0;
    overflow: hidden;
    background:
      radial-gradient(circle at 20% 30%, color-mix(in srgb, var(--color-accent) 6%, transparent), transparent 50%),
      radial-gradient(circle at 80% 70%, color-mix(in srgb, var(--color-success) 4%, transparent), transparent 50%),
      var(--color-bg);
    background-attachment: local;
  }

  .empty {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }
  .action {
    padding: 0.5rem 1rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-elevated);
    color: var(--color-text);
    cursor: pointer;
  }
  .action:hover {
    border-color: var(--color-accent);
  }

  .dock {
    position: absolute;
    left: 50%;
    bottom: 1rem;
    transform: translateX(-50%);
    display: flex;
    gap: 0.4rem;
    padding: 0.4rem;
    background: color-mix(in srgb, var(--color-bg-elevated) 90%, transparent);
    backdrop-filter: blur(8px);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-md);
    z-index: 1000;
  }
  .dock-item {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.4rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-subtle);
    color: var(--color-text);
    font-family: var(--font-mono);
    font-size: 0.75rem;
    cursor: pointer;
    transition: border-color var(--transition);
  }
  .dock-item:hover {
    border-color: var(--color-accent);
  }

  .snap-preview {
    position: absolute;
    pointer-events: none;
    background: color-mix(in srgb, var(--color-accent) 18%, transparent);
    border: 2px solid var(--color-accent);
    border-radius: var(--radius-md);
    box-shadow: 0 0 0 1px color-mix(in srgb, var(--color-accent) 30%, transparent);
    transition:
      left 100ms ease,
      top 100ms ease,
      width 100ms ease,
      height 100ms ease;
    z-index: 5;
  }
</style>
