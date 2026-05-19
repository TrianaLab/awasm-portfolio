<script lang="ts">
  import Desktop from './Desktop.svelte';
  import ResumeView from './ResumeView.svelte';
  import ThemeToggle from './ThemeToggle.svelte';
  import ModeToggle from './ModeToggle.svelte';
  import DownloadButton from './DownloadButton.svelte';
  import GitHubRepoCard from './GitHubRepoCard.svelte';
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
    <GitHubRepoCard repo="TrianaLab/awasm-portfolio" />
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
