<script lang="ts">
  import ExternalLink from './ExternalLink.svelte';
  import GitHubRepoCard from './GitHubRepoCard.svelte';
  import { GITHUB_REPO } from '../lib/config';
  import type { Resume } from '../lib/schema';

  let { resume }: { resume: Resume | null } = $props();

  const basics = $derived(resume?.basics ?? {});
</script>

<footer class="site-footer">
  <div class="inner">
    <div class="contact tap-links">
      {#if basics.email}
        <a href="mailto:{basics.email}">{basics.email}</a>
      {/if}
      <ul class="profiles">
        {#each basics.profiles ?? [] as profile (profile.network)}
          {#if profile.url}
            <li><ExternalLink href={profile.url}>{profile.network}</ExternalLink></li>
          {/if}
        {/each}
      </ul>
    </div>

    <div class="meta">
      <p class="colophon">
        Go compiled to WebAssembly, Svelte and xterm.js. This page and the PDF both come from one
        JSON Resume document, rendered at runtime.
      </p>
      <GitHubRepoCard repo={GITHUB_REPO} />
    </div>
  </div>
</footer>

<style>
  .site-footer {
    margin-top: auto;
    border-top: 1px solid var(--color-border);
    background: var(--color-bg-elevated);
  }
  .inner {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem 3rem;
    justify-content: space-between;
    max-width: var(--page-max);
    margin: 0 auto;
    padding: 2rem var(--page-pad);
  }
  .contact {
    display: flex;
    flex-direction: column;
    gap: 1.1rem;
    font-size: 0.9rem;
  }
  .profiles {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .meta {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.9rem;
    max-width: 46ch;
  }
  .colophon {
    margin: 0;
    font-size: 0.8rem;
    color: var(--color-text-muted);
  }
</style>
