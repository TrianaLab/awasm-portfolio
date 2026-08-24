<script lang="ts">
  import ExternalLink from './ExternalLink.svelte';
  import { formatMonthYear } from '../lib/resume-select';

  let {
    title,
    subtitle,
    startDate,
    endDate,
    url,
    summary,
    highlights,
    /** Point-in-time entries (certificates, awards) show one date, never a
        range — otherwise an entry with no end date reads as "→ Present". */
    point = false,
    /** When set, highlights collapse behind a disclosure with this label. */
    highlightsLabel,
  }: {
    title?: string;
    subtitle?: string;
    startDate?: string;
    endDate?: string;
    url?: string;
    summary?: string;
    highlights?: string[];
    point?: boolean;
    highlightsLabel?: string;
  } = $props();

  const dateRange = $derived(
    point ? formatMonthYear(startDate ?? endDate) : formatRange(startDate, endDate),
  );

  function formatRange(start?: string, end?: string): string {
    if (!start && !end) return '';
    const left = formatMonthYear(start);
    const right = end ? formatMonthYear(end) : 'Present';
    return `${left} → ${right}`;
  }
</script>

<article class="entry">
  <header class="entry-head">
    <div class="entry-title-row">
      <h3 class="entry-title">
        {#if url}
          <ExternalLink href={url}>{title ?? ''}</ExternalLink>
        {:else}
          {title ?? ''}
        {/if}
      </h3>
      {#if dateRange}<p class="entry-dates mono">{dateRange}</p>{/if}
    </div>
    {#if subtitle}<p class="entry-subtitle">{subtitle}</p>{/if}
  </header>
  {#if summary}<p class="entry-summary">{summary}</p>{/if}
  {#if highlights && highlights.length > 0}
    {#if highlightsLabel}
      <details class="entry-details">
        <summary>{highlightsLabel} ({highlights.length})</summary>
        <ul class="entry-highlights">
          {#each highlights as h (h)}
            <li>{h}</li>
          {/each}
        </ul>
      </details>
    {:else}
      <ul class="entry-highlights">
        {#each highlights as h (h)}
          <li>{h}</li>
        {/each}
      </ul>
    {/if}
  {/if}
</article>

<style>
  .entry {
    padding: 1.1rem 1.25rem;
    background: var(--color-bg-elevated);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    transition: border-color var(--transition);
  }
  .entry:hover {
    border-color: var(--color-border-strong);
  }
  .entry-head {
    margin-bottom: 0.5rem;
  }
  .entry-title-row {
    display: flex;
    flex-wrap: wrap;
    align-items: baseline;
    justify-content: space-between;
    gap: 0.25rem 1rem;
  }
  .entry-title {
    font-size: 1.02rem;
    font-weight: 600;
    margin: 0;
    color: var(--color-text);
  }
  .entry-dates {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 0.78rem;
    white-space: nowrap;
  }
  .entry-subtitle {
    margin: 0.2rem 0 0 0;
    color: var(--color-accent);
    font-size: 0.9rem;
    font-weight: 500;
  }
  .entry-summary {
    margin: 0;
    max-width: 80ch;
    color: var(--color-text);
    font-size: 0.92rem;
  }
  .entry-details {
    margin-top: 0.6rem;
    font-size: 0.88rem;
  }
  .entry-details summary {
    display: inline-flex;
    align-items: center;
    min-height: 32px;
    color: var(--color-text-muted);
    cursor: pointer;
  }
  .entry-details summary:hover {
    color: var(--color-text);
  }
  .entry-highlights {
    margin: 0.5rem 0 0 0;
    padding-left: 1.2rem;
    color: var(--color-text);
    font-size: 0.9rem;
  }
  .entry-highlights li {
    margin-bottom: 0.2rem;
  }
</style>
