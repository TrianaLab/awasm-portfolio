<script lang="ts">
  type Theme = { mode: 'dark' | 'light'; toggle: () => void };
  let { theme }: { theme: Theme } = $props();
</script>

<button class="icon-btn" onclick={theme.toggle} aria-label="Toggle theme">
  {#if theme.mode === 'dark'}
    <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor" aria-hidden="true">
      <path d="M12 4V2m0 20v-2m8-8h2M2 12h2m13.66-5.66 1.41-1.41M4.93 19.07l1.41-1.41m0-11.32L4.93 4.93m14.14 14.14-1.41-1.41M12 7a5 5 0 1 0 0 10 5 5 0 0 0 0-10Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" fill="none"/>
    </svg>
  {:else}
    <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor" aria-hidden="true">
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79Z"/>
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
  .icon-btn:hover {
    background: var(--color-bg-subtle);
    color: var(--color-text);
  }
</style>
