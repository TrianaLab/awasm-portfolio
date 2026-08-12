<script lang="ts">
  import DownloadButton from './DownloadButton.svelte';
  import ExternalLink from './ExternalLink.svelte';
  import type { Resume } from '../lib/schema';
  import { CTA, FEATURED, PRINCIPLES } from '../lib/portfolio';
  import {
    experienceGroups,
    featuredProject,
    formatYear,
    headlineTitle,
    selectedSystems,
    shortOrgName,
    upstreamContributions,
  } from '../lib/resume-select';

  let { resume }: { resume: Resume } = $props();

  const basics = $derived(resume.basics ?? {});
  const title = $derived(headlineTitle(resume));
  const featured = $derived(featuredProject(resume));
  const systems = $derived(selectedSystems(resume));
  const upstream = $derived(upstreamContributions(resume));
  const groups = $derived(experienceGroups(resume));

  function span(start?: string, end?: string): string {
    const from = formatYear(start);
    const to = end ? formatYear(end) : 'Present';
    return from ? `${from} – ${to}` : to;
  }
</script>

<!-- Hero ------------------------------------------------------------- -->
<section class="hero" aria-labelledby="hero-name">
  <p class="kicker mono">
    {#if basics.location?.city}{basics.location.city}, {basics.location.countryCode ?? ''}{/if}
  </p>
  <h1 id="hero-name">{basics.name}</h1>
  {#if title}
    <p class="hero-role">{title}</p>
  {/if}
  {#if basics.summary}
    <p class="hero-summary">{basics.summary}</p>
  {/if}
  <div class="hero-actions">
    <a class="btn btn-primary" href={CTA.work.href}>{CTA.work.label}</a>
    <DownloadButton label={CTA.resume.label} />
    <a class="btn" href={CTA.terminal.href}>{CTA.terminal.label}</a>
  </div>
  <p class="hero-note">
    The terminal is a Go CLI compiled to WebAssembly, reading the same résumé data as this
    page. You don't have to open it to read anything here.
  </p>
</section>

<!-- Work -------------------------------------------------------------- -->
<section id="work" class="band" tabindex="-1" aria-labelledby="work-h">
  <h2 id="work-h" class="section-h">Work</h2>

  {#if featured}
    <article class="featured" aria-labelledby="featured-h">
      <p class="kicker mono">{FEATURED.eyebrow}</p>
      <h3 id="featured-h">{shortOrgName(featured.organization)}</h3>
      <p class="thesis">{FEATURED.thesis}</p>

      <div class="featured-grid">
        <div>
          <h4>The problem</h4>
          <p>{FEATURED.problem}</p>
        </div>
        <div>
          <h4>The approach</h4>
          <p>{FEATURED.approach}</p>
        </div>
        <div class="featured-system">
          <h4>The system</h4>
          <p>{featured.summary}</p>
        </div>
      </div>

      <dl class="facts">
        {#if featured.position}
          <div><dt>Role</dt><dd>{featured.position}</dd></div>
        {/if}
        <div><dt>Since</dt><dd class="mono">{formatYear(featured.startDate)}</dd></div>
        <div><dt>Licence</dt><dd>Open source</dd></div>
      </dl>

      <p class="featured-links tap-links">
        <ExternalLink href={FEATURED.website} context="{shortOrgName(featured.organization)} project website">
          Visit the Pacto website
        </ExternalLink>
        {#if featured.url}
          <ExternalLink href={featured.url} context="{shortOrgName(featured.organization)} on GitHub">
            View the GitHub repository
          </ExternalLink>
        {/if}
      </p>
    </article>
  {/if}

  <h3 class="sub-h">Selected systems</h3>
  <ul class="cards tap-links">
    {#each systems as system (system.url)}
      <li class="card">
        <p class="kicker mono">{system.position ?? 'Open source'}</p>
        <h4>{shortOrgName(system.organization)}</h4>
        <p>{system.summary}</p>
        {#if system.url}
          <ExternalLink href={system.url}>
            {system.url.replace(/^https?:\/\//, '')}
          </ExternalLink>
        {/if}
      </li>
    {/each}
  </ul>

  {#if upstream.length > 0}
    <h3 class="sub-h">Upstream contributions</h3>
    <ul class="upstream tap-links">
      {#each upstream as item (item.url)}
        <li>
          <span class="upstream-org">
            {#if item.url}
              <ExternalLink href={item.url}>{item.organization}</ExternalLink>
            {:else}
              {item.organization}
            {/if}
          </span>
          <span class="upstream-summary">{item.summary}</span>
        </li>
      {/each}
    </ul>
  {/if}
</section>

<!-- Experience -------------------------------------------------------- -->
<section id="experience" class="band" tabindex="-1" aria-labelledby="experience-h">
  <h2 id="experience-h" class="section-h">Experience</h2>
  <ol class="timeline">
    {#each groups as group, gi (group.name + gi)}
      <li class="company" class:current={!group.roles[0].endDate}>
        <div class="company-head tap-links">
          <h3>
            {#if group.url}
              <ExternalLink href={group.url}>{group.name}</ExternalLink>
            {:else}
              {group.name}
            {/if}
          </h3>
          <span class="company-span mono">{span(group.startDate, group.roles[0].endDate)}</span>
        </div>
        <ol class="roles">
          {#each group.roles as r, ri (ri)}
            <li class="role">
              <div class="role-head">
                <h4>{r.position}</h4>
                <span class="role-span mono">{span(r.startDate, r.endDate)}</span>
              </div>
              <!-- Only the most recent role in each company block carries its
                   summary; the full history stays on the résumé page. -->
              {#if ri === 0 && r.summary}
                <p class="role-summary">{r.summary}</p>
              {/if}
            </li>
          {/each}
        </ol>
      </li>
    {/each}
  </ol>
  <p class="band-foot">
    <a class="btn" href="#/resume">Read the full résumé</a>
  </p>
</section>

<!-- About ------------------------------------------------------------- -->
<section id="about" class="band" tabindex="-1" aria-labelledby="about-h">
  <h2 id="about-h" class="section-h">How I work</h2>
  <ul class="cards">
    {#each PRINCIPLES as p (p.title)}
      <li class="card principle">
        <h3>{p.title}</h3>
        <p>{p.body}</p>
      </li>
    {/each}
  </ul>

  {#if resume.skills && resume.skills.length > 0}
    <h3 class="sub-h">Tools I reach for</h3>
    <dl class="skills">
      {#each resume.skills as skill (skill.name)}
        <div>
          <dt>{skill.name}</dt>
          <dd>{(skill.keywords ?? []).join(' · ')}</dd>
        </div>
      {/each}
    </dl>
  {/if}
</section>

<style>
  .hero {
    max-width: var(--page-max);
    margin: 0 auto;
    padding: clamp(2.5rem, 8vw, 5rem) var(--page-pad) clamp(2rem, 5vw, 3.5rem);
  }
  .kicker {
    margin: 0 0 0.6rem;
    font-size: 0.75rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  .hero h1 {
    margin: 0;
    font-size: clamp(2.25rem, 6vw, 3.75rem);
    letter-spacing: -0.02em;
  }
  .hero-role {
    margin: 0.35rem 0 1.25rem;
    font-size: clamp(1.05rem, 2.4vw, 1.4rem);
    font-weight: 550;
    color: var(--color-accent);
  }
  .hero-summary {
    margin: 0 0 1.75rem;
    max-width: 62ch;
    font-size: clamp(1rem, 1.6vw, 1.12rem);
    color: var(--color-text);
  }
  .hero-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.6rem;
  }
  .hero-note {
    margin: 1.5rem 0 0;
    max-width: 58ch;
    font-size: 0.85rem;
    color: var(--color-text-muted);
  }

  .band {
    max-width: var(--page-max);
    margin: 0 auto;
    padding: clamp(2rem, 5vw, 3.5rem) var(--page-pad);
    border-top: 1px solid var(--color-border);
  }
  /* Scroll targets are focused programmatically after in-page navigation so
     assistive tech follows along. They are not tab stops or controls, so they
     take the focus without drawing a ring around a whole page section. */
  .band:focus,
  .band:focus-visible {
    outline: none;
  }
  .section-h {
    margin: 0 0 1.75rem;
    font-size: 0.78rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: var(--color-text-muted);
    font-weight: 600;
  }
  .sub-h {
    margin: 2.75rem 0 1rem;
    font-size: 1rem;
    font-weight: 600;
  }
  .band-foot {
    margin: 2rem 0 0;
  }

  /* Featured project ------------------------------------------------- */
  .featured {
    padding: clamp(1.25rem, 3vw, 2rem);
    border: 1px solid var(--color-border);
    border-left: 3px solid var(--color-accent);
    border-radius: var(--radius-md);
    background: var(--color-bg-elevated);
  }
  .featured h3 {
    margin: 0;
    font-family: var(--font-mono);
    font-size: clamp(1.75rem, 4.5vw, 2.5rem);
    letter-spacing: -0.01em;
  }
  .thesis {
    margin: 0.5rem 0 1.75rem;
    max-width: 52ch;
    font-size: clamp(1.05rem, 2vw, 1.3rem);
    font-weight: 550;
    color: var(--color-text);
  }
  .featured-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 1.5rem;
  }
  .featured-grid h4,
  .card h4 {
    margin: 0 0 0.4rem;
    font-size: 0.72rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  .featured-grid p {
    margin: 0;
    font-size: 0.92rem;
  }
  .featured-system {
    grid-column: 1 / -1;
  }
  .featured-system p {
    max-width: 78ch;
  }
  .facts {
    display: flex;
    flex-wrap: wrap;
    gap: 1.75rem;
    margin: 1.75rem 0 0;
    padding-top: 1.25rem;
    border-top: 1px solid var(--color-border);
  }
  .facts dt {
    font-size: 0.7rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  .facts dd {
    margin: 0.15rem 0 0;
    font-size: 0.92rem;
    font-weight: 550;
  }
  .featured-links {
    display: flex;
    flex-wrap: wrap;
    /* Row gap clears the .tap-links hit boxes when these wrap on a phone. */
    gap: 0.6rem 1.5rem;
    margin: 1.25rem 0 0;
  }

  /* Card grids ------------------------------------------------------- */
  .cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 1rem;
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .card {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    padding: 1.15rem 1.25rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-elevated);
  }
  .card .kicker {
    margin: 0;
  }
  .card h4 {
    margin: 0;
    font-family: var(--font-mono);
    font-size: 1.05rem;
    letter-spacing: 0;
    text-transform: none;
    color: var(--color-text);
  }
  .card p {
    margin: 0;
    font-size: 0.9rem;
    color: var(--color-text);
  }
  .principle h3 {
    margin: 0;
    font-size: 1.02rem;
  }
  .principle p {
    color: var(--color-text-muted);
  }

  /* Upstream --------------------------------------------------------- */
  .upstream {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 0.9rem 2rem;
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .upstream li {
    padding-left: 0.9rem;
    border-left: 2px solid var(--color-border);
  }
  .upstream-org {
    display: block;
    font-weight: 600;
    font-size: 0.92rem;
  }
  .upstream-summary {
    display: block;
    /* Clears the .tap-links hit-box overlay on the organisation link above. */
    margin-top: 0.5rem;
    font-size: 0.85rem;
    color: var(--color-text-muted);
  }

  /* Experience ------------------------------------------------------- */
  .timeline {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .company {
    padding: 1.25rem 1.35rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-elevated);
  }
  .company.current {
    border-color: color-mix(in srgb, var(--color-accent) 45%, var(--color-border));
  }
  .company-head {
    display: flex;
    flex-wrap: wrap;
    align-items: baseline;
    justify-content: space-between;
    gap: 0.5rem 1rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid var(--color-border);
  }
  .company-head h3 {
    margin: 0;
    font-size: 1.15rem;
  }
  .company-span,
  .role-span {
    font-size: 0.78rem;
    color: var(--color-text-muted);
    white-space: nowrap;
  }
  .roles {
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .role {
    padding-top: 0.85rem;
  }
  .role-head {
    display: flex;
    flex-wrap: wrap;
    align-items: baseline;
    justify-content: space-between;
    gap: 0.25rem 1rem;
  }
  .role-head h4 {
    margin: 0;
    font-size: 0.97rem;
    font-weight: 600;
    text-transform: none;
    letter-spacing: 0;
    color: var(--color-text);
  }
  .role-summary {
    margin: 0.4rem 0 0;
    max-width: 78ch;
    font-size: 0.9rem;
    color: var(--color-text-muted);
  }

  /* Skills ----------------------------------------------------------- */
  .skills {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 0.9rem 2rem;
    margin: 0;
  }
  .skills dt {
    font-size: 0.9rem;
    font-weight: 600;
  }
  .skills dd {
    margin: 0.15rem 0 0;
    font-family: var(--font-mono);
    font-size: 0.8rem;
    color: var(--color-text-muted);
  }
</style>
