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

  const MIN_W = 360;
  const MIN_H = 220;

  function maximize() {
    const w = desktopEl?.clientWidth ?? window.innerWidth;
    const h = desktopEl?.clientHeight ?? window.innerHeight;
    manager.toggleMaximize(win.id, w, h);
  }

  // ─── Drag from the chrome bar ───────────────────────────────────────
  let dragOffset: { x: number; y: number } | null = null;
  let lastChromeClickAt = 0;
  const DBL_CLICK_MS = 350;

  function desktopRect(): { left: number; top: number; w: number; h: number } | null {
    if (!desktopEl) return null;
    const r = desktopEl.getBoundingClientRect();
    return { left: r.left, top: r.top, w: desktopEl.clientWidth, h: desktopEl.clientHeight };
  }

  function onChromePointerDown(event: PointerEvent) {
    if ((event.target as HTMLElement).closest('.traffic')) return; // dot click
    manager.focus(win.id);
    // Manual double-click detection: setPointerCapture below swallows
    // the native dblclick event, so we track click cadence ourselves
    // (matching macOS / Windows: double-click chrome → toggle maximize).
    const now = performance.now();
    if (now - lastChromeClickAt < DBL_CLICK_MS) {
      lastChromeClickAt = 0;
      maximize();
      return;
    }
    lastChromeClickAt = now;
    dragOffset = { x: event.clientX - win.x, y: event.clientY - win.y };
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
  }

  function onChromePointerMove(event: PointerEvent) {
    if (!dragOffset) return;
    manager.move(win.id, event.clientX - dragOffset.x, event.clientY - dragOffset.y);
    const desk = desktopRect();
    if (desk) {
      manager.updateSnapHint(event.clientX - desk.left, event.clientY - desk.top, desk.w, desk.h);
    }
  }

  function onChromePointerUp(event: PointerEvent) {
    dragOffset = null;
    (event.currentTarget as HTMLElement).releasePointerCapture(event.pointerId);
    const desk = desktopRect();
    if (desk) {
      manager.commitSnap(win.id, desk.w, desk.h);
    } else {
      manager.clearSnapHint();
    }
  }

  // ─── Resize from any edge or corner ─────────────────────────────────
  type Edge = 'n' | 's' | 'e' | 'w' | 'ne' | 'nw' | 'se' | 'sw';
  type ResizeState = {
    edge: Edge;
    startPointerX: number;
    startPointerY: number;
    startX: number;
    startY: number;
    startW: number;
    startH: number;
  };
  let resizeState: ResizeState | null = null;

  function onResizePointerDown(edge: Edge, event: PointerEvent) {
    event.stopPropagation();
    manager.focus(win.id);
    resizeState = {
      edge,
      startPointerX: event.clientX,
      startPointerY: event.clientY,
      startX: win.x,
      startY: win.y,
      startW: win.w,
      startH: win.h,
    };
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
  }

  function onResizePointerMove(event: PointerEvent) {
    if (!resizeState) return;
    const dx = event.clientX - resizeState.startPointerX;
    const dy = event.clientY - resizeState.startPointerY;
    const { edge, startX, startY, startW, startH } = resizeState;

    let newX = startX;
    let newY = startY;
    let newW = startW;
    let newH = startH;

    // Horizontal axis
    if (edge.includes('e')) {
      newW = Math.max(MIN_W, startW + dx);
    } else if (edge.includes('w')) {
      // Clamp so the right edge stays put; the left edge can't push past
      // the min-width frontier.
      const proposedW = Math.max(MIN_W, startW - dx);
      const consumedDx = startW - proposedW;
      newX = startX + consumedDx;
      newW = proposedW;
    }

    // Vertical axis
    if (edge.includes('s')) {
      newH = Math.max(MIN_H, startH + dy);
    } else if (edge.includes('n')) {
      const proposedH = Math.max(MIN_H, startH - dy);
      const consumedDy = startH - proposedH;
      newY = Math.max(0, startY + consumedDy);
      newH = proposedH;
    }

    if (newX !== startX || newY !== startY) {
      manager.move(win.id, newX, newY);
    }
    if (newW !== startW || newH !== startH) {
      manager.resize(win.id, newW, newH);
    }
  }

  function onResizePointerUp(event: PointerEvent) {
    resizeState = null;
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
        <button class="dot dot-red" aria-label="Close window" onclick={() => manager.close(win.id)}></button>
        <button class="dot dot-yellow" aria-label="Minimize window" onclick={() => manager.minimize(win.id)}></button>
        <button class="dot dot-green" aria-label="Maximize window" onclick={maximize}></button>
      </div>
      <span class="title">{win.title}</span>
    </div>

    <div class="content">{@render children()}</div>

    <!-- 4 edges + 4 corners; each one is a thin transparent strip with
         the appropriate resize cursor and a pointer-handler that drives
         the same generic resize logic. -->
    {#each ['n', 's', 'e', 'w', 'ne', 'nw', 'se', 'sw'] as const as edge (edge)}
      <div
        class="rz rz-{edge}"
        role="separator"
        aria-label="Resize {edge}"
        aria-orientation={edge === 'n' || edge === 's' ? 'horizontal' : 'vertical'}
        onpointerdown={(e) => onResizePointerDown(edge, e)}
        onpointermove={onResizePointerMove}
        onpointerup={onResizePointerUp}
        onpointercancel={onResizePointerUp}
      ></div>
    {/each}
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
    cursor: default;
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
  .traffic:hover .dot::before {
    opacity: 1;
  }
  .dot {
    position: relative;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    border: none;
    padding: 0;
    cursor: pointer;
    transition: filter var(--transition);
  }
  .dot::before {
    content: '';
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: 'JetBrains Mono', monospace;
    font-size: 9px;
    font-weight: 700;
    color: rgba(0, 0, 0, 0.55);
    line-height: 1;
    opacity: 0;
    transition: opacity 120ms ease;
    pointer-events: none;
  }
  .dot:hover {
    filter: brightness(1.05);
  }
  .dot-red {
    background: #ff5f56;
  }
  .dot-red::before {
    content: '×';
    font-size: 12px;
  }
  .dot-yellow {
    background: #ffbd2e;
  }
  .dot-yellow::before {
    content: '−';
    font-size: 11px;
  }
  .dot-green {
    background: #27c93f;
  }
  .dot-green::before {
    content: '⛶';
    font-size: 9px;
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
    margin-right: 60px;
  }

  .content {
    flex: 1;
    min-height: 0;
    display: flex;
    overflow: hidden;
  }

  /* ─── Resize handles ───────────────────────────────────────────────
     Hit areas are 6 px thick along edges, 12 px at corners. Corners
     sit on top of edges so diagonal cursors win at the intersection. */
  .rz {
    position: absolute;
    touch-action: none;
    background: transparent;
    z-index: 2;
  }
  .rz-n {
    top: 0;
    left: 0;
    right: 0;
    height: 6px;
    cursor: ns-resize;
  }
  .rz-s {
    bottom: 0;
    left: 0;
    right: 0;
    height: 6px;
    cursor: ns-resize;
  }
  .rz-e {
    top: 0;
    right: 0;
    bottom: 0;
    width: 6px;
    cursor: ew-resize;
  }
  .rz-w {
    top: 0;
    left: 0;
    bottom: 0;
    width: 6px;
    cursor: ew-resize;
  }
  .rz-nw {
    top: 0;
    left: 0;
    width: 12px;
    height: 12px;
    cursor: nwse-resize;
    z-index: 3;
  }
  .rz-ne {
    top: 0;
    right: 0;
    width: 12px;
    height: 12px;
    cursor: nesw-resize;
    z-index: 3;
  }
  .rz-sw {
    bottom: 0;
    left: 0;
    width: 12px;
    height: 12px;
    cursor: nesw-resize;
    z-index: 3;
  }
  .rz-se {
    bottom: 0;
    right: 0;
    width: 14px;
    height: 14px;
    cursor: nwse-resize;
    z-index: 3;
    /* Keep the classic visible grip on the SE corner so users know
       the window is resizable. */
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
  }
  .rz-se:hover {
    opacity: 0.9;
  }
</style>
