<script lang="ts">
  import type { Resume } from '../lib/schema';
  import Section from './Section.svelte';
  import TimelineEntry from './TimelineEntry.svelte';
  import SkillCloud from './SkillCloud.svelte';
  import DownloadButton from './DownloadButton.svelte';
  import ExternalLink from './ExternalLink.svelte';
  import { CTA } from '../lib/portfolio';

  let { resume }: { resume: Resume } = $props();

  const basics = $derived(resume.basics ?? {});
  const skillItems = $derived(
    (resume.skills ?? []).map((s) => ({
      primary: s.name,
      secondary: s.level,
      tags: s.keywords,
    })),
  );
  const interestItems = $derived(
    (resume.interests ?? []).map((i) => ({
      primary: i.name,
      tags: i.keywords,
    })),
  );
  const languageItems = $derived(
    (resume.languages ?? []).map((l) => ({
      primary: l.language,
      secondary: l.fluency,
    })),
  );
</script>

<div class="resume">
  <header class="hero">
    {#if basics.image}<img class="avatar" src={basics.image} alt="" />{/if}
    <div class="hero-text">
      <p class="kicker mono">Résumé</p>
      <h1 class="hero-name">{basics.name ?? 'Unnamed'}</h1>
      {#if basics.label}<p class="hero-label">{basics.label}</p>{/if}
      {#if basics.summary}<p class="hero-summary">{basics.summary}</p>{/if}
      <ul class="hero-meta">
        {#if basics.location?.city}
          <li>
            {basics.location.city}{basics.location.region ? `, ${basics.location.region}` : ''}
          </li>
        {/if}
        {#if basics.email}
          <li><a href="mailto:{basics.email}">{basics.email}</a></li>
        {/if}
        {#if basics.url}
          <li><a href={basics.url}>{basics.url.replace(/^https?:\/\//, '')}</a></li>
        {/if}
        {#each basics.profiles ?? [] as p (p.network)}
          {#if p.url}
            <li><ExternalLink href={p.url}>{p.network}</ExternalLink></li>
          {/if}
        {/each}
      </ul>
      <div class="hero-actions">
        <DownloadButton label={CTA.resume.label} primary />
        <a class="btn" href="#/">Back to overview</a>
      </div>
      <p class="hero-note">
        Generated in the browser when you click, from the same JSON Resume document this page
        renders. Vector text and no images, so applicant tracking systems can parse it.
      </p>
    </div>
  </header>

  {#if resume.work && resume.work.length > 0}
    <Section title="Experience">
      {#each resume.work as w, i (i)}
        <TimelineEntry
          title={w.position}
          subtitle={w.name}
          startDate={w.startDate}
          endDate={w.endDate}
          url={w.url}
          summary={w.summary}
          highlights={w.highlights}
        />
      {/each}
    </Section>
  {/if}

  {#if resume.volunteer && resume.volunteer.length > 0}
    <Section title="Open source & volunteering">
      {#each resume.volunteer as v, i (i)}
        <TimelineEntry
          title={v.organization}
          subtitle={v.position}
          startDate={v.startDate}
          endDate={v.endDate}
          url={v.url}
          summary={v.summary}
          highlights={v.highlights}
        />
      {/each}
    </Section>
  {/if}

  {#if resume.education && resume.education.length > 0}
    <Section title="Education">
      {#each resume.education as e, i (i)}
        <TimelineEntry
          title={e.studyType ? `${e.studyType} · ${e.area ?? ''}` : e.area}
          subtitle={e.institution}
          startDate={e.startDate}
          endDate={e.endDate}
          url={e.url}
          summary={e.score}
          highlights={e.courses}
          highlightsLabel="Coursework"
        />
      {/each}
    </Section>
  {/if}

  {#if resume.projects && resume.projects.length > 0}
    <Section title="Projects">
      {#each resume.projects as p, i (i)}
        <TimelineEntry
          title={p.name}
          startDate={p.startDate}
          endDate={p.endDate}
          url={p.url}
          summary={p.description}
          highlights={p.highlights}
        />
      {/each}
    </Section>
  {/if}

  {#if resume.certificates && resume.certificates.length > 0}
    <Section title="Certificates">
      {#each resume.certificates as c, i (i)}
        <TimelineEntry title={c.name} subtitle={c.issuer} startDate={c.date} url={c.url} point />
      {/each}
    </Section>
  {/if}

  {#if resume.awards && resume.awards.length > 0}
    <Section title="Awards">
      {#each resume.awards as a, i (i)}
        <TimelineEntry title={a.title} subtitle={a.awarder} startDate={a.date} summary={a.summary} point />
      {/each}
    </Section>
  {/if}

  {#if resume.publications && resume.publications.length > 0}
    <Section title="Publications">
      {#each resume.publications as p, i (i)}
        <TimelineEntry
          title={p.name}
          subtitle={p.publisher}
          startDate={p.releaseDate}
          url={p.url}
          summary={p.summary}
          point
        />
      {/each}
    </Section>
  {/if}

  {#if skillItems.length > 0}
    <Section title="Skills">
      <SkillCloud items={skillItems} title="Skills" />
    </Section>
  {/if}

  {#if languageItems.length > 0}
    <Section title="Languages">
      <SkillCloud items={languageItems} title="Languages" />
    </Section>
  {/if}

  {#if interestItems.length > 0}
    <Section title="Interests">
      <SkillCloud items={interestItems} title="Interests" />
    </Section>
  {/if}
</div>

<style>
  .resume {
    width: 100%;
    max-width: 900px;
    margin: 0 auto;
    padding: clamp(2rem, 5vw, 3.5rem) var(--page-pad) 4rem;
  }

  .hero {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 3rem;
    padding-bottom: 2rem;
    border-bottom: 1px solid var(--color-border);
  }

  .avatar {
    width: 96px;
    height: 96px;
    border-radius: var(--radius-lg);
    object-fit: cover;
    flex-shrink: 0;
  }

  .hero-text {
    flex: 1;
    min-width: 0;
  }
  .kicker {
    margin: 0 0 0.5rem;
    font-size: 0.72rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  .hero-name {
    font-size: clamp(1.9rem, 5vw, 2.6rem);
    letter-spacing: -0.02em;
    margin: 0 0 0.2rem 0;
    color: var(--color-text);
  }
  .hero-label {
    font-size: 1.1rem;
    color: var(--color-accent);
    margin: 0 0 1rem 0;
    font-weight: 550;
  }
  .hero-summary {
    margin: 0 0 1.25rem 0;
    max-width: 68ch;
    color: var(--color-text);
    line-height: 1.65;
  }
  .hero-meta {
    list-style: none;
    margin: 0 0 1.5rem;
    padding: 0;
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem 1.25rem;
    font-size: 0.88rem;
    color: var(--color-text-muted);
  }
  .hero-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.6rem;
  }
  .hero-note {
    margin: 1rem 0 0;
    max-width: 60ch;
    font-size: 0.8rem;
    color: var(--color-text-muted);
  }

  @media (max-width: 640px) {
    .hero {
      flex-direction: column;
      align-items: flex-start;
    }
    .avatar {
      width: 72px;
      height: 72px;
    }
  }
</style>
