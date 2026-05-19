<script lang="ts">
  import { onMount } from 'svelte';

  let { repo }: { repo: string } = $props();

  let tag = $state<string | null>(null);
  let stars = $state<number | null>(null);
  let forks = $state<number | null>(null);

  const CACHE_KEY = `awasm.gh.${repo}`;
  const CACHE_TTL_MS = 15 * 60 * 1000;

  type Cached = { tag: string | null; stars: number | null; forks: number | null; at: number };

  function readCache(): Cached | null {
    try {
      const raw = localStorage.getItem(CACHE_KEY);
      if (!raw) return null;
      const parsed = JSON.parse(raw) as Cached;
      if (Date.now() - parsed.at > CACHE_TTL_MS) return null;
      return parsed;
    } catch {
      return null;
    }
  }

  function writeCache(data: Cached) {
    try {
      localStorage.setItem(CACHE_KEY, JSON.stringify(data));
    } catch {
      /* quota / private mode — best-effort cache */
    }
  }

  async function fetchJson<T>(url: string): Promise<T | null> {
    try {
      const r = await fetch(url, { headers: { Accept: 'application/vnd.github+json' } });
      if (!r.ok) return null;
      return (await r.json()) as T;
    } catch {
      return null;
    }
  }

  onMount(async () => {
    const cached = readCache();
    if (cached) {
      tag = cached.tag;
      stars = cached.stars;
      forks = cached.forks;
      return;
    }

    type RepoMeta = { stargazers_count?: number; forks_count?: number };
    type Release = { tag_name?: string };

    const [meta, release] = await Promise.all([
      fetchJson<RepoMeta>(`https://api.github.com/repos/${repo}`),
      fetchJson<Release>(`https://api.github.com/repos/${repo}/releases/latest`),
    ]);
    stars = meta?.stargazers_count ?? null;
    forks = meta?.forks_count ?? null;
    tag = release?.tag_name ?? null;
    writeCache({ tag, stars, forks, at: Date.now() });
  });

  const compact = $derived(new Intl.NumberFormat('en', { notation: 'compact' }));
</script>

<a class="card" href="https://github.com/{repo}" target="_blank" rel="noreferrer" aria-label="GitHub repository {repo}">
  <span class="icon" aria-hidden="true">
    <svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor">
      <path d="M12 2A10 10 0 0 0 2 12c0 4.42 2.87 8.17 6.84 9.5.5.08.66-.23.66-.5v-1.69c-2.77.6-3.36-1.34-3.36-1.34-.46-1.16-1.11-1.47-1.11-1.47-.91-.62.07-.6.07-.6 1 .07 1.53 1.03 1.53 1.03.87 1.52 2.34 1.07 2.91.83.09-.65.35-1.09.63-1.34-2.22-.25-4.55-1.11-4.55-4.94 0-1.11.38-2 1.03-2.71-.1-.25-.45-1.29.1-2.68 0 0 .84-.27 2.75 1.02.79-.22 1.65-.33 2.5-.33s1.71.11 2.5.33c1.91-1.29 2.75-1.02 2.75-1.02.55 1.39.2 2.43.1 2.68.65.71 1.03 1.6 1.03 2.71 0 3.84-2.34 4.68-4.57 4.93.36.31.69.92.69 1.85V21c0 .27.16.59.67.5A10 10 0 0 0 22 12 10 10 0 0 0 12 2Z"/>
    </svg>
  </span>
  <span class="body">
    <span class="name">{repo}</span>
    <span class="meta">
      {#if tag}
        <span class="meta-item" title="Latest release">
          <svg viewBox="0 0 24 24" width="11" height="11" fill="currentColor" aria-hidden="true">
            <path d="M21.41 11.58 12.41 2.58A2 2 0 0 0 11 2H4a2 2 0 0 0-2 2v7a2 2 0 0 0 .59 1.42l9 9a2 2 0 0 0 2.83 0l7-7a2 2 0 0 0 0-2.84zM6.5 8A1.5 1.5 0 1 1 8 6.5 1.5 1.5 0 0 1 6.5 8z"/>
          </svg>
          {tag}
        </span>
      {/if}
      {#if stars !== null}
        <span class="meta-item" title="Stars">
          <svg viewBox="0 0 24 24" width="11" height="11" fill="currentColor" aria-hidden="true">
            <path d="m12 17.27 6.18 3.73-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"/>
          </svg>
          {compact.format(stars)}
        </span>
      {/if}
      {#if forks !== null}
        <span class="meta-item" title="Forks">
          <svg viewBox="0 0 16 16" width="11" height="11" fill="currentColor" aria-hidden="true">
            <path d="M5 5.372v.878c0 .414.336.75.75.75h4.5a.75.75 0 0 0 .75-.75v-.878a2.25 2.25 0 1 1 1.5 0v.878a2.25 2.25 0 0 1-2.25 2.25h-1.5v2.128a2.251 2.251 0 1 1-1.5 0V8.5h-1.5A2.25 2.25 0 0 1 3.5 6.25v-.878a2.25 2.25 0 1 1 1.5 0ZM5 3.25a.75.75 0 1 0-1.5 0 .75.75 0 0 0 1.5 0Zm6.75.75a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Zm-3 8.75a.75.75 0 1 0-1.5 0 .75.75 0 0 0 1.5 0Z"/>
          </svg>
          {compact.format(forks)}
        </span>
      {/if}
    </span>
  </span>
</a>

<style>
  .card {
    display: inline-flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.35rem 0.75rem 0.35rem 0.6rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-subtle);
    color: var(--color-text);
    text-decoration: none;
    transition: border-color var(--transition), background var(--transition);
  }
  .card:hover {
    border-color: var(--color-accent);
    background: var(--color-bg-elevated);
    text-decoration: none;
  }
  .icon {
    display: inline-flex;
    align-items: center;
    color: var(--color-text-muted);
  }
  .body {
    display: inline-flex;
    flex-direction: column;
    gap: 0.1rem;
  }
  .name {
    font-family: var(--font-mono);
    font-size: 0.78rem;
    font-weight: 600;
    color: var(--color-text);
    white-space: nowrap;
  }
  .meta {
    display: inline-flex;
    align-items: center;
    gap: 0.65rem;
    color: var(--color-text-muted);
    font-size: 0.7rem;
    font-family: var(--font-mono);
  }
  .meta-item {
    display: inline-flex;
    align-items: center;
    gap: 0.2rem;
  }

  @media (max-width: 720px) {
    .body {
      display: none;
    }
  }
</style>
