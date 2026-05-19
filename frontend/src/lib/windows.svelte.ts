// Window-manager state. A single rune store owns every window's
// geometry, z-stack, minimized flag. Pure data — the Desktop and
// Window components render against it.

export interface WindowState {
  id: string;
  title: string;
  x: number;
  y: number;
  w: number;
  h: number;
  z: number;
  minimized: boolean;
  // When maximized OR snapped, the previous geometry is stored so the
  // user can restore to their original window arrangement.
  previousGeometry?: { x: number; y: number; w: number; h: number };
}

export type SnapZone = 'left' | 'right' | 'top' | null;

export interface SnapHint {
  zone: SnapZone;
  rect: { x: number; y: number; w: number; h: number } | null;
}

const MIN_W = 360;
const MIN_H = 220;
const SNAP_THRESHOLD = 24;

let nextId = 1;
let nextZ = 1;

function makeId(prefix: string): string {
  return `${prefix}-${nextId++}`;
}

function snapRectFor(
  zone: SnapZone,
  desktopW: number,
  desktopH: number,
): { x: number; y: number; w: number; h: number } | null {
  switch (zone) {
    case 'left':
      return { x: 0, y: 0, w: Math.floor(desktopW / 2), h: desktopH };
    case 'right':
      return { x: Math.ceil(desktopW / 2), y: 0, w: Math.floor(desktopW / 2), h: desktopH };
    case 'top':
      return { x: 0, y: 0, w: desktopW, h: desktopH };
    default:
      return null;
  }
}

export interface WindowManager {
  readonly windows: WindowState[];
  readonly snapHint: SnapHint;
  open(title: string): WindowState;
  close(id: string): void;
  focus(id: string): void;
  minimize(id: string): void;
  restore(id: string): void;
  move(id: string, x: number, y: number): void;
  resize(id: string, w: number, h: number): void;
  toggleMaximize(id: string, desktopW: number, desktopH: number): void;
  isMaximized(id: string): boolean;
  updateSnapHint(pointerX: number, pointerY: number, desktopW: number, desktopH: number): void;
  clearSnapHint(): void;
  commitSnap(id: string, desktopW: number, desktopH: number): boolean;
}

export function createWindowManager(): WindowManager {
  const windows = $state<WindowState[]>([]);
  const snapHint = $state<SnapHint>({ zone: null, rect: null });

  function find(id: string): WindowState | undefined {
    return windows.find((w) => w.id === id);
  }

  function open(title: string): WindowState {
    const offset = (windows.length % 6) * 28;
    const win: WindowState = {
      id: makeId('term'),
      title,
      x: 40 + offset,
      y: 40 + offset,
      w: 720,
      h: 440,
      z: ++nextZ,
      minimized: false,
    };
    windows.push(win);
    return win;
  }

  function close(id: string) {
    const i = windows.findIndex((w) => w.id === id);
    if (i !== -1) windows.splice(i, 1);
  }

  function focus(id: string) {
    const win = find(id);
    if (!win) return;
    win.z = ++nextZ;
    win.minimized = false;
  }

  function minimize(id: string) {
    const win = find(id);
    if (win) win.minimized = true;
  }

  function restore(id: string) {
    focus(id);
  }

  function move(id: string, x: number, y: number) {
    const win = find(id);
    if (!win) return;
    win.x = Math.max(0, x);
    win.y = Math.max(0, y);
  }

  function resize(id: string, w: number, h: number) {
    const win = find(id);
    if (!win) return;
    win.w = Math.max(MIN_W, w);
    win.h = Math.max(MIN_H, h);
    // Any user-driven resize cancels the snapped/maximized state — a
    // future restore would jump back to a stale geometry the user didn't
    // intend.
    win.previousGeometry = undefined;
  }

  function toggleMaximize(id: string, desktopW: number, desktopH: number) {
    const win = find(id);
    if (!win) return;
    if (win.previousGeometry) {
      const prev = win.previousGeometry;
      win.x = prev.x;
      win.y = prev.y;
      win.w = prev.w;
      win.h = prev.h;
      win.previousGeometry = undefined;
      return;
    }
    win.previousGeometry = { x: win.x, y: win.y, w: win.w, h: win.h };
    win.x = 0;
    win.y = 0;
    win.w = Math.max(MIN_W, desktopW);
    win.h = Math.max(MIN_H, desktopH);
  }

  function isMaximized(id: string): boolean {
    const win = find(id);
    return !!win?.previousGeometry;
  }

  function detectZone(pointerX: number, pointerY: number, desktopW: number): SnapZone {
    if (pointerY <= SNAP_THRESHOLD) return 'top';
    if (pointerX <= SNAP_THRESHOLD) return 'left';
    if (pointerX >= desktopW - SNAP_THRESHOLD) return 'right';
    return null;
  }

  function updateSnapHint(pointerX: number, pointerY: number, desktopW: number, desktopH: number) {
    const zone = detectZone(pointerX, pointerY, desktopW);
    snapHint.zone = zone;
    snapHint.rect = snapRectFor(zone, desktopW, desktopH);
  }

  function clearSnapHint() {
    snapHint.zone = null;
    snapHint.rect = null;
  }

  // Apply the active snap (if any) to the window. Returns true when a
  // snap was committed so the caller can suppress the normal release.
  function commitSnap(id: string, desktopW: number, desktopH: number): boolean {
    const win = find(id);
    if (!win || !snapHint.zone) {
      clearSnapHint();
      return false;
    }
    const rect = snapRectFor(snapHint.zone, desktopW, desktopH);
    clearSnapHint();
    if (!rect) return false;
    // Snapshot the pre-snap geometry so toggleMaximize-style restore works.
    win.previousGeometry = { x: win.x, y: win.y, w: win.w, h: win.h };
    win.x = rect.x;
    win.y = rect.y;
    win.w = rect.w;
    win.h = rect.h;
    return true;
  }

  return {
    get windows() {
      return windows;
    },
    get snapHint() {
      return snapHint;
    },
    open,
    close,
    focus,
    minimize,
    restore,
    move,
    resize,
    toggleMaximize,
    isMaximized,
    updateSnapHint,
    clearSnapHint,
    commitSnap,
  };
}
