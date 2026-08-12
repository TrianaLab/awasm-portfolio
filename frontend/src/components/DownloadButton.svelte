<script lang="ts">
  import { fetchResume } from '../lib/wasm';
  import { downloadResumePdf } from '../lib/pdf';
  import { CTA } from '../lib/portfolio';

  let { label, primary = false }: { label?: string; primary?: boolean } = $props();

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
  class="btn"
  class:btn-primary={primary}
  class:btn-icon={!label}
  onclick={handleClick}
  disabled={downloading}
  aria-label={label ? undefined : CTA.resume.label}
  title={label ? undefined : CTA.resume.label}
>
  {#if downloading}
    <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true">
      <circle cx="12" cy="12" r="9" opacity="0.25" />
      <path d="M21 12a9 9 0 0 0-9-9" class="spin" />
    </svg>
  {:else}
    <svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor" aria-hidden="true">
      <path d="M12 3v12m0 0-4-4m4 4 4-4M5 21h14" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
  {/if}
  {#if label}<span>{downloading ? 'Generating…' : label}</span>{/if}
</button>

<style>
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
