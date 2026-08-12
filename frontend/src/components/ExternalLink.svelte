<script lang="ts">
  import type { Snippet } from 'svelte';

  // The visible text stays part of the accessible name (WCAG 2.5.3), and the
  // hidden suffix adds the context a screen-reader user would otherwise miss:
  // where the link goes and that it leaves the page.
  let {
    href,
    context,
    children,
  }: { href: string; context?: string; children: Snippet } = $props();
</script>

<a class="ext" {href} target="_blank" rel="noreferrer noopener">
  {@render children()}<span class="sr-only">{context ? ` — ${context}` : ''} (opens in a new tab)</span
  ><svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
    <path d="M7 17 17 7M9 7h8v8" />
  </svg>
</a>

<style>
  .ext {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    font-size: inherit;
    word-break: break-word;
  }
  .ext svg {
    flex: none;
    opacity: 0.7;
  }
</style>
