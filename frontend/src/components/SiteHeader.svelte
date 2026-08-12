<script lang="ts">
  import ThemeToggle from './ThemeToggle.svelte';
  import { NAV } from '../lib/portfolio';
  import { DOMAIN } from '../lib/config';
  import type { View } from '../lib/portfolio';

  type Theme = { mode: 'dark' | 'light'; toggle: () => void };
  let { view, section, theme }: { view: View; section: string | null; theme: Theme } = $props();

  let menuOpen = $state(false);
  let toggleEl = $state<HTMLButtonElement | null>(null);

  // A nav item is current when its route matches the active view, or when it
  // points at the home section the user has navigated to.
  function isCurrent(href: string): boolean {
    const target = href.replace('#/', '');
    if (target === 'resume' || target === 'terminal') return view === target;
    return view === 'home' && section === target;
  }

  // Collapsing the menu takes the focused link out of the layout. On a real
  // route change App moves focus to the new view, but tapping the link for the
  // route you are already on fires no hashchange, so focus would fall to
  // <body>. Handing it back to the toggle covers both: App overrides it when a
  // navigation does happen.
  function closeMenu() {
    if (!menuOpen) return;
    menuOpen = false;
    toggleEl?.focus();
  }
</script>

<svelte:window
  onkeydown={(e) => {
    if (e.key === 'Escape') closeMenu();
  }}
/>

<header class="topbar">
  <a class="brand" href="#/" onclick={closeMenu}>
    <span class="brand-mark" aria-hidden="true">~/</span>
    <span class="brand-name">{DOMAIN}</span>
  </a>

  <!-- One list at every viewport. Below 780px it collapses behind a text
       "Menu" button rather than an ambiguous glyph. -->
  <nav class="nav" class:open={menuOpen} id="primary-nav" aria-label="Primary">
    <ul>
      {#each NAV as item (item.href)}
        <li>
          <a
            href={item.href}
            aria-current={isCurrent(item.href) ? 'page' : undefined}
            onclick={closeMenu}
          >
            {item.label}
          </a>
        </li>
      {/each}
    </ul>
  </nav>

  <div class="tools">
    <ThemeToggle {theme} />
    <button
      type="button"
      class="btn menu-toggle"
      aria-expanded={menuOpen}
      aria-controls="primary-nav"
      bind:this={toggleEl}
      onclick={() => (menuOpen = !menuOpen)}
    >
      {menuOpen ? 'Close' : 'Menu'}
    </button>
  </div>
</header>

<style>
  .topbar {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    min-height: var(--header-h);
    padding: 0 var(--page-pad);
    border-bottom: 1px solid var(--color-border);
    background: color-mix(in srgb, var(--color-bg) 88%, transparent);
    backdrop-filter: blur(10px);
    position: sticky;
    top: 0;
    z-index: 20;
  }

  .brand {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    min-height: var(--tap);
    font-family: var(--font-mono);
    font-weight: 600;
    font-size: 0.95rem;
    color: var(--color-text);
    text-decoration: none;
    white-space: nowrap;
  }
  .brand:hover {
    opacity: 1;
    text-decoration: none;
  }
  .brand-mark {
    color: var(--color-accent);
  }

  .nav {
    margin-right: auto;
  }
  .nav ul {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .nav a {
    display: inline-flex;
    align-items: center;
    min-height: var(--tap);
    padding: 0 0.7rem;
    border-radius: var(--radius-sm);
    color: var(--color-text-muted);
    font-size: 0.9rem;
    font-weight: 500;
    text-decoration: none;
  }
  .nav a:hover {
    color: var(--color-text);
    background: var(--color-bg-subtle);
    opacity: 1;
    text-decoration: none;
  }
  .nav a[aria-current='page'] {
    color: var(--color-text);
    background: var(--color-bg-subtle);
    box-shadow: inset 0 -2px 0 var(--color-accent);
  }

  .tools {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin-left: auto;
  }

  .menu-toggle {
    display: none;
    padding: 0 0.9rem;
    font-size: 0.85rem;
  }

  @media (max-width: 780px) {
    .topbar {
      flex-wrap: wrap;
      gap: 0.5rem;
      padding: 0.5rem 1rem;
    }
    .menu-toggle {
      display: inline-flex;
    }
    /* Collapsed: the whole list leaves the layout and the a11y tree. */
    .nav {
      display: none;
      order: 3;
      flex-basis: 100%;
      margin: 0 0 0.25rem;
      border-top: 1px solid var(--color-border);
      padding-top: 0.4rem;
    }
    .nav.open {
      display: block;
    }
    .nav ul {
      flex-direction: column;
      align-items: stretch;
      gap: 0.15rem;
    }
    .nav a {
      font-size: 1rem;
      padding: 0 0.75rem;
    }
    .nav a[aria-current='page'] {
      box-shadow: none;
      background: var(--color-accent-soft);
      color: var(--color-accent);
    }
  }
</style>
