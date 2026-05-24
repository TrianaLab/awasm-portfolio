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
  // True for one render tick after a programmatic geometry change so
  // the Window component can apply CSS transitions (snap, maximize)
  // without lagging the drag/resize path.
  animating: boolean;
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
  open(title: string, desktopW?: number, desktopH?: number): WindowState;
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
  clearAnimating(id: string): void;
}

export function createWindowManager(): WindowManager {
  const windows = $state<WindowState[]>([]);
  const snapHint = $state<SnapHint>({ zone: null, rect: null });

  function find(id: string): WindowState | undefined {
    return windows.find((w) => w.id === id);
  }

  // initialGeometry sizes a fresh window relative to the desktop. On
  // wide viewports it lands ~70% × 70% centered; on phone-class
  // viewports it covers most of the screen so it's actually usable.
  // Subsequent opens cascade by `offset` so a second window isn't
  // perfectly overlapping the first.
  function initialGeometry(
    desktopW: number,
    desktopH: number,
    offset: number,
  ): { x: number; y: number; w: number; h: number } {
    const isPhone = desktopW < 640;
    const ratioW = isPhone ? 0.96 : 0.7;
    const ratioH = isPhone ? 0.88 : 0.7;
    const minW = isPhone ? 320 : 600;
    const minH = isPhone ? 320 : 380;
    const w = Math.max(minW, Math.round(desktopW * ratioW));
    const h = Math.max(minH, Math.round(desktopH * ratioH));
    const x = Math.max(0, Math.round((desktopW - w) / 2) + offset);
    const y = Math.max(0, Math.round((desktopH - h) / 2) + offset);
    return { x, y, w, h };
  }

  function open(title: string, desktopW = 1280, desktopH = 800): WindowState {
    const offset = (windows.length % 6) * 28;
    const geom = initialGeometry(desktopW, desktopH, offset);
    const win: WindowState = {
      id: makeId('term'),
      title,
      ...geom,
      z: ++nextZ,
      minimized: false,
      animating: true,
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
    win.animating = true;
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
    win.animating = true;
    win.x = rect.x;
    win.y = rect.y;
    win.w = rect.w;
    win.h = rect.h;
    return true;
  }

  function clearAnimating(id: string) {
    const win = find(id);
    if (win) win.animating = false;
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
    clearAnimating,
  };
}
