// Theme is persisted to localStorage and applied via the [data-theme]
// attribute on <html>. Svelte 5 runes give us a tiny reactive store.

const STORAGE_KEY = 'awasm.theme';
type Mode = 'dark' | 'light';

function read(): Mode {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === 'dark' || stored === 'light') return stored;
  return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
}

function apply(mode: Mode) {
  document.documentElement.dataset.theme = mode;
}

export function createTheme() {
  let mode = $state<Mode>(read());
  apply(mode);

  return {
    get mode() {
      return mode;
    },
    toggle() {
      mode = mode === 'dark' ? 'light' : 'dark';
      localStorage.setItem(STORAGE_KEY, mode);
      apply(mode);
    },
  };
}
