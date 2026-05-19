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
      <!-- Material Design Icons / Octicons github — solid silhouette, same shape mkdocs Material renders. -->
      <svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor" aria-hidden="true">
        <path d="M12 2A10 10 0 0 0 2 12c0 4.42 2.87 8.17 6.84 9.5.5.08.66-.23.66-.5v-1.69c-2.77.6-3.36-1.34-3.36-1.34-.46-1.16-1.11-1.47-1.11-1.47-.91-.62.07-.6.07-.6 1 .07 1.53 1.03 1.53 1.03.87 1.52 2.34 1.07 2.91.83.09-.65.35-1.09.63-1.34-2.22-.25-4.55-1.11-4.55-4.94 0-1.11.38-2 1.03-2.71-.1-.25-.45-1.29.1-2.68 0 0 .84-.27 2.75 1.02.79-.22 1.65-.33 2.5-.33s1.71.11 2.5.33c1.91-1.29 2.75-1.02 2.75-1.02.55 1.39.2 2.43.1 2.68.65.71 1.03 1.6 1.03 2.71 0 3.84-2.34 4.68-4.57 4.93.36.31.69.92.69 1.85V21c0 .27.16.59.67.5A10 10 0 0 0 22 12 10 10 0 0 0 12 2Z"/>
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
