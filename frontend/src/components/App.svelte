<script lang="ts">
  import { onMount, tick } from 'svelte';
  import SiteHeader from './SiteHeader.svelte';
  import SiteFooter from './SiteFooter.svelte';
  import HomeView from './HomeView.svelte';
  import { createTheme } from '../lib/theme.svelte';
  import { fetchResume } from '../lib/wasm';
  import { PAGE_TITLE } from '../lib/config';
  import { parseHash } from '../lib/portfolio';
  import { pageTitle, personJsonLd } from '../lib/resume-select';
  import type { Resume } from '../lib/schema';

  const theme = createTheme();

  let loc = $state(parseHash(window.location.hash));
  let resume = $state<Resume | null>(null);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let mainEl = $state<HTMLElement | null>(null);
  // The terminal (xterm + the window manager) is code-split and only pulled
  // in the first time someone asks for it. Once loaded it stays mounted so
  // buffers, history and window layout survive navigating away and back.
  let terminalModule = $state<Promise<typeof import('./TerminalPane.svelte')> | null>(null);

  async function refreshResume() {
    loading = true;
    error = null;
    try {
      resume = await fetchResume();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    void refreshResume();
    const onHash = () => {
      loc = parseHash(window.location.hash);
    };
    window.addEventListener('hashchange', onHash);
    return () => window.removeEventListener('hashchange', onHash);
  });

  // Route → document state: the view class drives the viewport-locking rules
  // in app.css, and the terminal bundle is requested on first use.
  $effect(() => {
    document.documentElement.dataset.view = loc.view;
    if (loc.view === 'terminal' && !terminalModule) {
      terminalModule = import('./TerminalPane.svelte');
    }
  });

  // Coming back from the terminal, re-read the document so mutations made
  // there (create/delete) show up on the résumé and home pages.
  let prevView: string | null = null;
  $effect(() => {
    const view = loc.view;
    if (prevView === 'terminal' && view !== 'terminal') void refreshResume();
    prevView = view;
  });

  // Focus + scroll management on navigation. A hash change does not move
  // focus on its own, so keyboard and screen-reader users would otherwise
  // stay parked in the nav after switching views.
  //
  // A section target only exists once the résumé has rendered, so a cold load
  // of #/experience runs this before there is anything to scroll to; tracking
  // readiness makes it re-run when the document lands. `lastNav` keeps that
  // second run from also firing on a plain load of #/, where focus and scroll
  // position must be left exactly where the browser put them.
  const resumeReady = $derived(resume !== null);
  let lastNav: string | null = null;
  $effect(() => {
    const { view, section } = loc;
    const ready = resumeReady;
    const key = `${view}:${section ?? ''}`;
    if (lastNav === null) {
      if (!section) {
        lastNav = key;
        return;
      }
    } else if (lastNav === key) {
      return; // only readiness changed
    }
    if (section && !ready) return;
    lastNav = key;
    void tick().then(() => {
      const target = section ? document.getElementById(section) : null;
      // preventScroll: `scroll-behavior: smooth` makes scrollIntoView an
      // animation, and focus() would start a competing block:'nearest' scroll
      // that lands the user somewhere in the middle of the section.
      if (target) {
        target.scrollIntoView({ block: 'start' });
        target.focus({ preventScroll: true });
      } else {
        window.scrollTo({ top: 0 });
        mainEl?.focus({ preventScroll: true });
      }
    });
  });

  // The skip link cannot navigate: `href="#main"` sets location.hash, which
  // parseHash reads as the home route, so the one control keyboard users are
  // guaranteed to hit would throw them off the view they were skipping into.
  function skipToMain(event: MouseEvent) {
    event.preventDefault();
    mainEl?.focus();
  }

  $effect(() => {
    document.title = pageTitle(resume, PAGE_TITLE);
  });

  // schema.org Person, generated from the canonical document rather than
  // hand-written into index.html (which would be a second copy of the facts).
  $effect(() => {
    const json = personJsonLd(resume);
    if (!json) return;
    const id = 'person-jsonld';
    const el = document.getElementById(id) ?? document.createElement('script');
    el.id = id;
    (el as HTMLScriptElement).type = 'application/ld+json';
    el.textContent = json;
    if (!el.parentNode) document.head.appendChild(el);
  });
</script>

<a class="skip-link" href="#main" onclick={skipToMain}>Skip to main content</a>

<SiteHeader view={loc.view} section={loc.section} {theme} />

<main id="main" class="main" tabindex="-1" bind:this={mainEl}>
  {#if loc.view !== 'terminal'}
    {#if loading && !resume}
      <p class="status" role="status">Loading résumé…</p>
    {:else if error && !resume}
      <p class="status error" role="alert">
        Could not load the résumé data: {error}
      </p>
    {:else if resume}
      <HomeView {resume} />
    {/if}
  {/if}

  <!-- Kept mounted once loaded; `hidden` takes it out of the layout, the
       accessibility tree and the tab order so only one view is ever live. -->
  {#if terminalModule}
    <div class="terminal-layer" hidden={loc.view !== 'terminal'}>
      {#await terminalModule then Module}
        <Module.default />
      {:catch}
        <p class="status error" role="alert">
          The terminal failed to load. Everything else on this site still works.
        </p>
      {/await}
    </div>
  {/if}
</main>

{#if loc.view !== 'terminal'}
  <SiteFooter {resume} />
{/if}

<style>
  .main {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }
  .main:focus {
    outline: none;
  }

  .terminal-layer {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }

  .status {
    margin: 0;
    padding: 4rem var(--page-pad);
    text-align: center;
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }
  .status.error {
    color: var(--color-danger);
  }
</style>
