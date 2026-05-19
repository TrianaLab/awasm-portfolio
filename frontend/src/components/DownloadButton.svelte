<script lang="ts">
  import { fetchResume } from '../lib/wasm';
  import { downloadResumePdf } from '../lib/pdf';

  let downloading = $state(false);

  // Always refetch on click so the PDF reflects the in-WASM state at
  // download time (the user can mutate it via the terminal between views).
  async function handleClick() {
    if (downloading) return;
    downloading = true;
    try {
      const fresh = await fetchResume();
      await downloadResumePdf(fresh);
    } finally {
      downloading = false;
    }
  }
</script>

<button
  type="button"
  class="pill"
  onclick={handleClick}
  disabled={downloading}
  aria-label="Download resume as PDF"
  title="Download resume as PDF"
>
  {#if downloading}
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true">
      <circle cx="12" cy="12" r="9" opacity="0.25" />
      <path d="M21 12a9 9 0 0 0-9-9" class="spin" />
    </svg>
  {:else}
    <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor" aria-hidden="true">
      <path d="M12 3v12m0 0-4-4m4 4 4-4M5 21h14" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
  {/if}
</button>

<style>
  .pill {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: var(--radius-sm);
    border: 1px solid transparent;
    background: transparent;
    color: var(--color-text-muted);
    cursor: pointer;
    transition: background var(--transition), color var(--transition);
  }
  .pill:hover:not(:disabled) {
    background: var(--color-bg-subtle);
    color: var(--color-text);
  }
  .pill:disabled {
    cursor: progress;
  }
  .spin {
    transform-origin: center;
    animation: rotate 0.9s linear infinite;
  }
  @keyframes rotate {
    to {
      transform: rotate(360deg);
    }
  }
</style>
