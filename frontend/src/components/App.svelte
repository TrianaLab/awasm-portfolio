<script lang="ts">
  import Desktop from './Desktop.svelte';
  import ResumeView from './ResumeView.svelte';
  import ThemeToggle from './ThemeToggle.svelte';
  import ModeToggle from './ModeToggle.svelte';
  import DownloadButton from './DownloadButton.svelte';
  import { createTheme } from '../lib/theme.svelte';
  import { fetchResume } from '../lib/wasm';
  import type { Resume } from '../lib/schema';

  const theme = createTheme();

  let mode = $state<'terminal' | 'resume'>('terminal');
  let resume = $state<Resume | null>(null);
  let loading = $state(false);
  let error = $state<string | null>(null);

  // Always refetch: the user can mutate the in-WASM resume from the
  // terminal (delete + create) and the UI must reflect the new state.
  async function refreshResume() {
    if (loading) return;
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
      void refreshResume();
    }
  }

  let desktop = $state<Desktop | undefined>();
</script>

<header class="topbar">
  <div class="brand">
    <span class="brand-mark">~/</span>
    <span class="brand-name">edudiaz.dev</span>
  </div>
  <nav class="actions">
    {#if mode === 'terminal'}
      <button
        class="icon-link"
        onclick={() => desktop?.open()}
        aria-label="Open a new terminal"
        title="Open a new terminal"
      >
        <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="M12 5v14M5 12h14" />
        </svg>
      </button>
    {/if}
    <ModeToggle {mode} onChange={setMode} />
    <DownloadButton />
    <a class="icon-link" href="https://github.com/TrianaLab/awasm-portfolio" target="_blank" rel="noreferrer" aria-label="GitHub">
      <!-- Font Awesome free brands/github (CC BY 4.0) — same glyph mkdocs Material uses. -->
      <svg viewBox="0 0 496 512" width="20" height="20" fill="currentColor" aria-hidden="true">
        <path d="M165.9 397.4c0 2-2.3 3.6-5.2 3.6-3.3 .3-5.6-1.3-5.6-3.6 0-2 2.3-3.6 5.2-3.6 3-.3 5.6 1.3 5.6 3.6zm-31.1-4.5c-.7 2 1.3 4.3 4.3 4.9 2.6 1 5.6 0 6.2-2s-1.3-4.3-4.3-5.2c-2.6-.7-5.5 .3-6.2 2.3zm44.2-1.7c-2.9 .7-4.9 2.6-4.6 4.9 .3 2 2.9 3.3 5.9 2.6 2.9-.7 4.9-2.6 4.6-4.6-.3-1.9-3-3.2-5.9-2.9zM244.8 8C106.1 8 0 113.3 0 252c0 110.9 69.8 205.8 169.5 239.2 12.8 2.3 17.3-5.6 17.3-12.1 0-6.2-.3-40.4-.3-61.4 0 0-70 15-84.7-29.8 0 0-11.4-29.1-27.8-36.6 0 0-22.9-15.7 1.6-15.4 0 0 24.9 2 38.6 25.8 21.9 38.6 58.6 27.5 72.9 20.9 2.3-16 8.8-27.1 16-33.7-55.9-6.2-112.3-14.3-112.3-110.5 0-27.5 7.6-41.3 23.6-58.9-2.6-6.5-11.1-33.3 2.6-67.9 20.9-6.5 69 27 69 27 20-5.6 41.5-8.5 62.8-8.5s42.8 2.9 62.8 8.5c0 0 48.1-33.6 69-27 13.7 34.7 5.2 61.4 2.6 67.9 16 17.7 25.8 31.5 25.8 58.9 0 96.5-58.9 104.2-114.8 110.5 9.2 7.9 17 22.9 17 46.4 0 33.7-.3 75.4-.3 83.6 0 6.5 4.6 14.4 17.3 12.1C428.2 457.8 496 362.9 496 252 496 113.3 383.5 8 244.8 8zM97.2 352.9c-1.3 1-1 3.3 .7 5.2 1.6 1.6 3.9 2.3 5.2 1 1.3-1 1-3.3-.7-5.2-1.6-1.6-3.9-2.3-5.2-1zm-10.8-8.1c-.7 1.3 .3 2.9 2.3 3.9 1.6 1 3.6 .7 4.3-.7 .7-1.3-.3-2.9-2.3-3.9-2-.6-3.6-.3-4.3 .7zm32.4 35.6c-1.6 1.3-1 4.3 1.3 6.2 2.3 2.3 5.2 2.6 6.5 1 1.3-1.3 .7-4.3-1.3-6.2-2.2-2.3-5.2-2.6-6.5-1zm-11.4-14.7c-1.6 1-1.6 3.6 0 5.9 1.6 2.3 4.3 3.3 5.6 2.3 1.6-1.3 1.6-3.9 0-6.2-1.4-2.3-4-3.3-5.6-2z"/>
      </svg>
    </a>
    <ThemeToggle {theme} />
  </nav>
</header>

<main class="stage" class:terminal-mode={mode === 'terminal'}>
  {#if mode === 'terminal'}
    <Desktop bind:this={desktop} />
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
    min-height: 0; /* break the flex shrink hold so children can scroll */
    display: flex;
    padding: 1.5rem;
    overflow: hidden;
  }
  /* No padding in terminal mode so the desktop fills the viewport. */
  .stage.terminal-mode {
    padding: 0;
  }

  /* Plain <button> reset for the icon-buttons added beside icon-links. */
  button.icon-link {
    background: transparent;
    border: none;
    cursor: pointer;
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
