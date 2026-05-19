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
  // When maximized, the previous geometry is stored so the green button
  // can toggle back to the user's window arrangement.
  previousGeometry?: { x: number; y: number; w: number; h: number };
}

const MIN_W = 360;
const MIN_H = 220;

let nextId = 1;
let nextZ = 1;

function makeId(prefix: string): string {
  return `${prefix}-${nextId++}`;
}

export interface WindowManager {
  readonly windows: WindowState[];
  open(title: string): WindowState;
  close(id: string): void;
  focus(id: string): void;
  minimize(id: string): void;
  restore(id: string): void;
  move(id: string, x: number, y: number): void;
  resize(id: string, w: number, h: number): void;
  toggleMaximize(id: string, desktopW: number, desktopH: number): void;
  isMaximized(id: string): boolean;
}

export function createWindowManager(): WindowManager {
  const windows = $state<WindowState[]>([]);

  function find(id: string): WindowState | undefined {
    return windows.find((w) => w.id === id);
  }

  function open(title: string): WindowState {
    // Cascade new windows so they don't overlap perfectly.
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
    // Any user-driven resize cancels the maximized state — restoring would
    // jump back to a stale geometry the user didn't intend.
    win.previousGeometry = undefined;
  }

  function toggleMaximize(id: string, desktopW: number, desktopH: number) {
    const win = find(id);
    if (!win) return;
    if (win.previousGeometry) {
      // Currently maximized — restore.
      const prev = win.previousGeometry;
      win.x = prev.x;
      win.y = prev.y;
      win.w = prev.w;
      win.h = prev.h;
      win.previousGeometry = undefined;
      return;
    }
    // Maximize: snapshot current geometry, then fill the desktop.
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

  return {
    get windows() {
      return windows;
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
  };
}
