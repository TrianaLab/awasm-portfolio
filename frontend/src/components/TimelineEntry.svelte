<script lang="ts">
  let {
    title,
    subtitle,
    startDate,
    endDate,
    url,
    summary,
    highlights,
  }: {
    title?: string;
    subtitle?: string;
    startDate?: string;
    endDate?: string;
    url?: string;
    summary?: string;
    highlights?: string[];
  } = $props();

  const dateRange = $derived(formatRange(startDate, endDate));

  function formatRange(start?: string, end?: string): string {
    if (!start && !end) return '';
    const left = formatDate(start);
    const right = end ? formatDate(end) : 'Present';
    return `${left} → ${right}`;
  }

  function formatDate(d?: string): string {
    if (!d) return '';
    const parsed = new Date(d);
    if (Number.isNaN(parsed.getTime())) return d;
    return parsed.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
  }
</script>

<article class="entry">
  <header class="entry-head">
    <div class="entry-title-row">
      <h3 class="entry-title">{title ?? ''}</h3>
      {#if url}
        <a class="entry-link" href={url} target="_blank" rel="noreferrer" aria-label="External link">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M7 17 17 7M9 7h8v8" />
          </svg>
        </a>
      {/if}
    </div>
    {#if subtitle}<p class="entry-subtitle">{subtitle}</p>{/if}
    {#if dateRange}<p class="entry-dates">{dateRange}</p>{/if}
  </header>
  {#if summary}<p class="entry-summary">{summary}</p>{/if}
  {#if highlights && highlights.length > 0}
    <ul class="entry-highlights">
      {#each highlights as h (h)}
        <li>{h}</li>
      {/each}
    </ul>
  {/if}
</article>

<style>
  .entry {
    padding: 1rem 1.25rem;
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    transition: border-color var(--transition), transform var(--transition);
  }
  .entry:hover {
    border-color: var(--color-accent);
    transform: translateY(-1px);
  }
  .entry-head {
    margin-bottom: 0.5rem;
  }
  .entry-title-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .entry-title {
    font-size: 1rem;
    font-weight: 600;
    margin: 0;
    color: var(--color-text);
  }
  .entry-link {
    color: var(--color-text-muted);
    transition: color var(--transition);
  }
  .entry-link:hover {
    color: var(--color-accent);
  }
  .entry-subtitle {
    margin: 0.15rem 0 0 0;
    color: var(--color-accent);
    font-size: 0.9rem;
    font-weight: 500;
  }
  .entry-dates {
    margin: 0.15rem 0 0 0;
    color: var(--color-text-muted);
    font-family: var(--font-mono);
    font-size: 0.8rem;
  }
  .entry-summary {
    margin: 0;
    color: var(--color-text);
    font-size: 0.92rem;
  }
  .entry-highlights {
    margin: 0.5rem 0 0 0;
    padding-left: 1.2rem;
    color: var(--color-text);
    font-size: 0.92rem;
  }
  .entry-highlights li {
    margin-bottom: 0.2rem;
  }
</style>
