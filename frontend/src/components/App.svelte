<script lang="ts">
  import Terminal from './Terminal.svelte';
  import ResumeView from './ResumeView.svelte';
  import ThemeToggle from './ThemeToggle.svelte';
  import ModeToggle from './ModeToggle.svelte';
  import DownloadButton from './DownloadButton.svelte';
  import { createTheme } from '../lib/theme';
  import { fetchResume } from '../lib/wasm';
  import type { Resume } from '../lib/schema';

  const theme = createTheme();

  let mode = $state<'terminal' | 'resume'>('terminal');
  let resume = $state<Resume | null>(null);
  let loading = $state(false);
  let error = $state<string | null>(null);

  async function ensureResume() {
    if (resume || loading) return;
    loading = true;
    error = null;
    try {
      resume = await fetchResume();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  function setMode(next: 'terminal' | 'resume') {
    mode = next;
    if (next === 'resume') {
      void ensureResume();
    }
  }
</script>

<header class="topbar">
  <div class="brand">
    <span class="brand-mark">~/</span>
    <span class="brand-name">edudiaz.dev</span>
  </div>
  <nav class="actions">
    <ModeToggle {mode} onChange={setMode} />
    <DownloadButton {resume} />
    <a class="icon-link" href="https://github.com/TrianaLab/awasm-portfolio" target="_blank" rel="noreferrer" aria-label="GitHub">
      <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor" aria-hidden="true">
        <path d="M12 .5C5.65.5.5 5.65.5 12c0 5.08 3.29 9.39 7.86 10.91.58.11.79-.25.79-.55 0-.27-.01-.99-.02-1.94-3.2.69-3.88-1.54-3.88-1.54-.52-1.32-1.27-1.67-1.27-1.67-1.04-.71.08-.7.08-.7 1.15.08 1.76 1.18 1.76 1.18 1.02 1.75 2.68 1.24 3.33.95.1-.74.4-1.24.73-1.52-2.55-.29-5.24-1.28-5.24-5.69 0-1.26.45-2.29 1.18-3.09-.12-.29-.51-1.46.11-3.04 0 0 .97-.31 3.18 1.18.92-.26 1.91-.39 2.89-.39.98 0 1.97.13 2.89.39 2.21-1.49 3.17-1.18 3.17-1.18.63 1.58.24 2.75.12 3.04.73.8 1.18 1.83 1.18 3.09 0 4.42-2.69 5.4-5.26 5.69.41.36.78 1.06.78 2.14 0 1.55-.01 2.8-.01 3.18 0 .31.21.67.8.55C20.21 21.39 23.5 17.08 23.5 12 23.5 5.65 18.35.5 12 .5Z"/>
      </svg>
    </a>
    <ThemeToggle {theme} />
  </nav>
</header>

<main class="stage">
  {#if mode === 'terminal'}
    <Terminal />
  {:else if loading}
    <div class="loading">Loading resume…</div>
  {:else if error}
    <div class="error">Failed to load resume: {error}</div>
  {:else if resume}
    <ResumeView {resume} />
  {/if}
</main>

<style>
  .topbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 2rem;
    border-bottom: 1px solid var(--color-border);
    background: var(--color-bg-elevated);
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .brand {
    display: flex;
    align-items: baseline;
    gap: 0.25rem;
    font-family: var(--font-mono);
    font-weight: 600;
    color: var(--color-text);
  }
  .brand-mark {
    color: var(--color-accent);
  }
  .brand-name {
    color: var(--color-text);
  }

  .actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .icon-link {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: var(--radius-sm);
    color: var(--color-text-muted);
    transition: background var(--transition), color var(--transition);
  }
  .icon-link:hover {
    background: var(--color-bg-subtle);
    color: var(--color-text);
  }

  .stage {
    flex: 1;
    display: flex;
    padding: 1.5rem;
    overflow: hidden;
  }

  .loading,
  .error {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }
  .error {
    color: var(--color-danger);
  }

  @media (max-width: 640px) {
    .topbar {
      padding: 0.75rem 1rem;
    }
    .stage {
      padding: 0.75rem;
    }
  }
</style>
