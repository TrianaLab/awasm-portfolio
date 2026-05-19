<script lang="ts">
  import type { Resume } from '../lib/schema';
  import Section from './Section.svelte';
  import TimelineEntry from './TimelineEntry.svelte';
  import SkillCloud from './SkillCloud.svelte';

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
    {#if basics.image}<img class="avatar" src={basics.image} alt={basics.name ?? 'avatar'} />{/if}
    <div class="hero-text">
      <h1 class="hero-name">{basics.name ?? 'Unnamed'}</h1>
      {#if basics.label}<p class="hero-label">{basics.label}</p>{/if}
      {#if basics.summary}<p class="hero-summary">{basics.summary}</p>{/if}
      <ul class="hero-meta">
        {#if basics.location?.city}
          <li>
            <span class="meta-icon" aria-hidden="true">📍</span>
            {basics.location.city}{basics.location.region ? `, ${basics.location.region}` : ''}
          </li>
        {/if}
        {#if basics.email}
          <li><a href="mailto:{basics.email}">{basics.email}</a></li>
        {/if}
        {#if basics.url}
          <li><a href={basics.url} target="_blank" rel="noreferrer">{basics.url.replace(/^https?:\/\//, '')}</a></li>
        {/if}
      </ul>
      {#if basics.profiles && basics.profiles.length > 0}
        <ul class="profiles">
          {#each basics.profiles as p, i (i)}
            <li>
              <a href={p.url} target="_blank" rel="noreferrer">{p.network}</a>
            </li>
          {/each}
        </ul>
      {/if}
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
        />
      {/each}
    </Section>
  {/if}

  {#if resume.volunteer && resume.volunteer.length > 0}
    <Section title="Open source & volunteering">
      {#each resume.volunteer as v, i (i)}
        <TimelineEntry
          title={v.position}
          subtitle={v.organization}
          startDate={v.startDate}
          endDate={v.endDate}
          url={v.url}
          summary={v.summary}
          highlights={v.highlights}
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
        <TimelineEntry
          title={c.name}
          subtitle={c.issuer}
          startDate={c.date}
          url={c.url}
        />
      {/each}
    </Section>
  {/if}

  {#if resume.awards && resume.awards.length > 0}
    <Section title="Awards">
      {#each resume.awards as a, i (i)}
        <TimelineEntry
          title={a.title}
          subtitle={a.awarder}
          startDate={a.date}
          summary={a.summary}
        />
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
    flex: 1;
    max-width: 920px;
    margin: 0 auto;
    padding: 0 1rem;
    overflow-y: auto;
  }

  .hero {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 2.5rem;
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
  }
  .hero-name {
    font-size: 2rem;
    margin: 0 0 0.25rem 0;
    color: var(--color-text);
  }
  .hero-label {
    font-size: 1.1rem;
    color: var(--color-accent);
    margin: 0 0 1rem 0;
    font-weight: 500;
  }
  .hero-summary {
    margin: 0 0 1rem 0;
    color: var(--color-text);
    line-height: 1.65;
  }
  .hero-meta,
  .profiles {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    font-size: 0.9rem;
  }
  .hero-meta li,
  .profiles li {
    color: var(--color-text-muted);
  }
  .meta-icon {
    margin-right: 0.25rem;
  }
  .profiles {
    margin-top: 0.5rem;
    font-family: var(--font-mono);
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
    .hero-name {
      font-size: 1.5rem;
    }
  }
</style>
