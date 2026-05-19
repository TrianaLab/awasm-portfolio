<script lang="ts">
  import type { Snippet } from 'svelte';
  let { title, children }: { title: string; children: Snippet } = $props();
</script>

<section class="section">
  <h2 class="section-title">{title}</h2>
  <div class="section-body">
    {@render children()}
  </div>
</section>

<style>
  .section {
    margin-bottom: 2.5rem;
  }
  .section-title {
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.12em;
    color: var(--color-text-muted);
    font-weight: 600;
    margin: 0 0 1rem 0;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid var(--color-border);
  }
  .section-body {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }
</style>
