<script lang="ts">
  let {
    title,
    items,
  }: {
    title?: string;
    items?: { primary?: string; secondary?: string; tags?: string[] }[];
  } = $props();
</script>

{#if items && items.length > 0}
  <div class="cloud">
    {#each items as item (item.primary)}
      <div class="cloud-item">
        <div class="cloud-head">
          <span class="primary">{item.primary ?? ''}</span>
          {#if item.secondary}<span class="secondary">{item.secondary}</span>{/if}
        </div>
        {#if item.tags && item.tags.length > 0}
          <div class="tags">
            {#each item.tags as tag (tag)}
              <span class="tag">{tag}</span>
            {/each}
          </div>
        {/if}
      </div>
    {/each}
  </div>
{:else if title}
  <p class="empty">No {title.toLowerCase()} listed.</p>
{/if}

<style>
  .cloud {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 0.75rem;
  }
  .cloud-item {
    padding: 0.9rem 1rem;
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
  }
  .cloud-head {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    margin-bottom: 0.4rem;
  }
  .primary {
    font-weight: 600;
    color: var(--color-text);
  }
  .secondary {
    color: var(--color-text-muted);
    font-size: 0.8rem;
    font-family: var(--font-mono);
  }
  .tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.3rem;
  }
  .tag {
    padding: 0.15rem 0.5rem;
    background: var(--color-accent-soft);
    color: var(--color-accent);
    border-radius: 999px;
    font-size: 0.75rem;
    font-weight: 500;
  }
  .empty {
    color: var(--color-text-muted);
    font-style: italic;
  }
</style>
