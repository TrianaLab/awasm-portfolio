<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { WindowState, WindowManager } from '../lib/windows.svelte';

  let {
    win,
    manager,
    desktopEl,
    children,
  }: {
    win: WindowState;
    manager: WindowManager;
    desktopEl: HTMLElement | null;
    children: Snippet;
  } = $props();

  function maximize() {
    const w = desktopEl?.clientWidth ?? window.innerWidth;
    const h = desktopEl?.clientHeight ?? window.innerHeight;
    manager.toggleMaximize(win.id, w, h);
  }

  // Drag from the chrome bar.
  let dragOffset: { x: number; y: number } | null = null;

  function onChromePointerDown(event: PointerEvent) {
    if ((event.target as HTMLElement).closest('.traffic')) return; // dot click
    manager.focus(win.id);
    dragOffset = { x: event.clientX - win.x, y: event.clientY - win.y };
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
  }

  function onChromePointerMove(event: PointerEvent) {
    if (!dragOffset) return;
    manager.move(win.id, event.clientX - dragOffset.x, event.clientY - dragOffset.y);
  }

  function onChromePointerUp(event: PointerEvent) {
    dragOffset = null;
    (event.currentTarget as HTMLElement).releasePointerCapture(event.pointerId);
  }

  // Resize from the SE corner.
  let resizeOffset: { x: number; y: number; w: number; h: number } | null = null;

  function onResizePointerDown(event: PointerEvent) {
    event.stopPropagation();
    manager.focus(win.id);
    resizeOffset = { x: event.clientX, y: event.clientY, w: win.w, h: win.h };
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
  }

  function onResizePointerMove(event: PointerEvent) {
    if (!resizeOffset) return;
    manager.resize(
      win.id,
      resizeOffset.w + (event.clientX - resizeOffset.x),
      resizeOffset.h + (event.clientY - resizeOffset.y),
    );
  }

  function onResizePointerUp(event: PointerEvent) {
    resizeOffset = null;
    (event.currentTarget as HTMLElement).releasePointerCapture(event.pointerId);
  }
</script>

{#if !win.minimized}
  <div
    class="window"
    role="dialog"
    aria-label={win.title}
    tabindex="-1"
    style:left="{win.x}px"
    style:top="{win.y}px"
    style:width="{win.w}px"
    style:height="{win.h}px"
    style:z-index={win.z}
    onpointerdown={() => manager.focus(win.id)}
  >
    <div
      class="chrome"
      role="toolbar"
      aria-label="Window controls and drag handle"
      tabindex="0"
      onpointerdown={onChromePointerDown}
      onpointermove={onChromePointerMove}
      onpointerup={onChromePointerUp}
      onpointercancel={onChromePointerUp}
    >
      <div class="traffic">
        <button
          class="dot dot-red"
          aria-label="Close window"
          onclick={() => manager.close(win.id)}
        ></button>
        <button
          class="dot dot-yellow"
          aria-label="Minimize window"
          onclick={() => manager.minimize(win.id)}
        ></button>
        <button
          class="dot dot-green"
          aria-label="Maximize window"
          onclick={maximize}
        ></button>
      </div>
      <span class="title">{win.title}</span>
    </div>

    <div class="content">{@render children()}</div>

    <div
      class="resize-handle"
      role="separator"
      aria-label="Resize"
      aria-orientation="horizontal"
      onpointerdown={onResizePointerDown}
      onpointermove={onResizePointerMove}
      onpointerup={onResizePointerUp}
      onpointercancel={onResizePointerUp}
    ></div>
  </div>
{/if}

<style>
  .window {
    position: absolute;
    display: flex;
    flex-direction: column;
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-lg);
    overflow: hidden;
    min-width: 360px;
    min-height: 220px;
  }

  .chrome {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.55rem 0.8rem;
    background: var(--color-bg-subtle);
    border-bottom: 1px solid var(--color-border);
    cursor: grab;
    user-select: none;
    touch-action: none;
  }
  .chrome:active {
    cursor: grabbing;
  }

  .traffic {
    display: flex;
    gap: 0.4rem;
  }
  .dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    border: none;
    padding: 0;
    cursor: pointer;
    transition: filter var(--transition);
  }
  .dot:hover {
    filter: brightness(1.15);
  }
  .dot-red {
    background: #ff5f56;
  }
  .dot-yellow {
    background: #ffbd2e;
  }
  .dot-green {
    background: #27c93f;
  }

  .title {
    flex: 1;
    text-align: center;
    font-family: var(--font-mono);
    font-size: 0.8rem;
    color: var(--color-text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    /* Compensate for the traffic lights on the left so the title is visually centred. */
    margin-right: 60px;
  }

  .content {
    flex: 1;
    min-height: 0;
    display: flex;
    overflow: hidden;
  }

  .resize-handle {
    position: absolute;
    right: 0;
    bottom: 0;
    width: 14px;
    height: 14px;
    cursor: nwse-resize;
    background: linear-gradient(
      135deg,
      transparent 50%,
      var(--color-text-muted) 50%,
      var(--color-text-muted) 60%,
      transparent 60%,
      transparent 70%,
      var(--color-text-muted) 70%,
      var(--color-text-muted) 80%,
      transparent 80%
    );
    opacity: 0.45;
    touch-action: none;
  }
  .resize-handle:hover {
    opacity: 0.9;
  }
</style>
