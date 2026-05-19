<script lang="ts">
  import Desktop from './Desktop.svelte';
  import ResumeView from './ResumeView.svelte';
  import ThemeToggle from './ThemeToggle.svelte';
  import ModeToggle from './ModeToggle.svelte';
  import DownloadButton from './DownloadButton.svelte';
  import GitHubRepoCard from './GitHubRepoCard.svelte';
  import { createTheme } from '../lib/theme.svelte';
  import { fetchResume } from '../lib/wasm';
  import { DOMAIN, GITHUB_REPO } from '../lib/config';
  import type { Resume } from '../lib/schema';

  const theme = createTheme();

  let mode = $state<'terminal' | 'resume'>('terminal');
  let resume = $state<Resume | null>(null);
  let loading = $state(false);
  let error = $state<string | null>(null);

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
    <span class="brand-name">{DOMAIN}</span>
  </div>

  <nav class="actions" aria-label="App actions">
    {#if mode === 'terminal'}
      <button
        type="button"
        class="pill icon-only"
        onclick={() => desktop?.open()}
        aria-label="Open a new terminal"
        title="Open a new terminal"
      >
        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="M12 5v14M5 12h14" />
        </svg>
      </button>
    {/if}

    <ModeToggle {mode} onChange={setMode} />
    <DownloadButton />
    <ThemeToggle {theme} />
    <span class="divider" aria-hidden="true"></span>
    <GitHubRepoCard repo={GITHUB_REPO} />
  </nav>
</header>

<main class="stage">
  <!-- Both layers stay mounted so terminal state (xterm buffers, command
       history, window arrangement, scroll position) survives the round
       trip terminal → resume → terminal. Only the active layer is
       interactive; the inactive one fades out + ignores pointer events. -->
  <div class="layer terminal-layer" class:hidden={mode !== 'terminal'} aria-hidden={mode !== 'terminal'}>
    <Desktop bind:this={desktop} />
  </div>
  <div class="layer resume-layer" class:hidden={mode !== 'resume'} aria-hidden={mode !== 'resume'}>
    {#if loading}
      <div class="loading">Loading resume…</div>
    {:else if error}
      <div class="error">Failed to load resume: {error}</div>
    {:else if resume}
      <ResumeView {resume} />
    {/if}
  </div>
</main>

<style>
  .topbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.6rem 1.25rem;
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
    gap: 0.4rem;
  }

  /* Every interactive item in the topbar normalises to the same height
     so the action bar reads as a single uniform group. The actual icon
     buttons sit inside their own components and follow the same
     dimensions via the .pill helper. */
  .actions :global(.pill),
  .actions :global(.segmented),
  .actions :global(.card) {
    height: 36px;
    box-sizing: border-box;
  }
  .actions :global(.segmented) {
    padding: 0 2px;
  }

  .pill {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    border-radius: var(--radius-sm);
    border: 1px solid transparent;
    background: transparent;
    color: var(--color-text-muted);
    cursor: pointer;
    transition: background var(--transition), color var(--transition), border-color var(--transition);
  }
  .pill:hover {
    background: var(--color-bg-subtle);
    color: var(--color-text);
  }

  .divider {
    width: 1px;
    height: 22px;
    background: var(--color-border);
    margin: 0 0.2rem;
  }

  .stage {
    position: relative;
    flex: 1;
    min-height: 0;
    overflow: hidden;
  }

  .layer {
    position: absolute;
    inset: 0;
    display: flex;
    transition: opacity 180ms ease;
  }
  .resume-layer {
    overflow-y: auto;
    padding: 1.5rem;
  }
  .layer.hidden {
    opacity: 0;
    pointer-events: none;
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

  @media (max-width: 720px) {
    .topbar {
      padding: 0.5rem 0.75rem;
      gap: 0.5rem;
    }
    .brand-name {
      display: none;
    }
    .divider {
      display: none;
    }
  }
</style>
