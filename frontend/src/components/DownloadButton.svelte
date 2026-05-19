<script lang="ts">
  import { fetchResume } from '../lib/wasm';
  import { downloadResumePdf } from '../lib/pdf';
  import type { Resume } from '../lib/schema';

  let { resume }: { resume: Resume | null } = $props();

  let downloading = $state(false);

  async function handleClick() {
    if (downloading) return;
    downloading = true;
    try {
      const data = resume ?? (await fetchResume());
      await downloadResumePdf(data);
    } finally {
      downloading = false;
    }
  }
</script>

<button
  class="icon-btn"
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
  .icon-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: var(--radius-sm);
    color: var(--color-text-muted);
    transition: background var(--transition), color var(--transition);
  }
  .icon-btn:hover:not(:disabled) {
    background: var(--color-bg-subtle);
    color: var(--color-text);
  }
  .icon-btn:disabled {
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
