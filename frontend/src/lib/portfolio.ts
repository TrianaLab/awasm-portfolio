// Portfolio presentation model.
//
// Boundary: resume.json (JSON Resume) owns every factual claim — biography,
// employment, education, skills, project descriptions, dates. This file owns
// presentation and editorial emphasis only: navigation, ordering, CTA copy,
// headings, terminal suggestions and SEO strings.
//
// Résumé entries are referenced by their canonical `url` (a stable key that
// already exists in the document) and never copied. resume-schema.test.ts runs
// the real CLI over the canonical document and fails if any URL or resource
// name referenced below stops resolving — that is the intended guard rail.

export type View = 'home' | 'terminal';

/** Home-page section ids that double as routes (`#/work` → home + scroll). */
export const HOME_SECTIONS = ['work', 'experience', 'education', 'about'] as const;

export type Location = { view: View; section: string | null };

/** Parses `window.location.hash` into a view + optional home section. */
export function parseHash(hash: string): Location {
  const path = hash.replace(/^#\/?/, '').replace(/\/+$/, '').toLowerCase();
  if (path === 'terminal') return { view: 'terminal', section: null };
  if ((HOME_SECTIONS as readonly string[]).includes(path)) {
    return { view: 'home', section: path };
  }
  return { view: 'home', section: null };
}

export const NAV: { href: string; label: string }[] = [
  { href: '#/work', label: 'Work' },
  { href: '#/experience', label: 'Experience' },
  { href: '#/education', label: 'Education' },
  { href: '#/about', label: 'About' },
  { href: '#/terminal', label: 'Terminal' },
];

export const CTA = {
  work: { href: '#/work', label: 'Explore my work' },
  terminal: { href: '#/terminal', label: 'Open the terminal' },
  resume: { label: 'Download résumé' },
} as const;

export type FeaturedConfig = {
  url: string;
  // Presentation, not résumé data: JSON Resume has one `url` per entry and
  // that one is the repository. The project site belongs here instead of as a
  // proprietary field in the canonical document.
  website: string;
  websiteLabel: string;
  eyebrow: string;
  thesis: string;
  problem: string;
  approach: string;
};

/**
 * The projects the home page is built around, in display order. Each is
 * referenced by the canonical URL of its JSON Resume entry; the copy below is
 * editorial framing, not a restatement of the entry's own description.
 */
export const FEATURED: FeaturedConfig[] = [
  {
    url: 'https://github.com/TrianaLab/mira',
    website: 'https://miradb.dev',
    websiteLabel: 'Visit the Mira website',
    eyebrow: 'Latest work',
    thesis:
      'Short term memory for autonomous systems. When an agent is handed a live incident, the telemetry should be something it can read and reason over directly, not a dashboard built for a human to squint at.',
    problem:
      'Every telemetry backend I have run converts OTLP into something else on the way in: rows in a column store, Parquet files, a time series index. Each conversion is a place fidelity can go missing, and most of them want a cluster and a database alongside before you get to ask the first question.',
    approach:
      'Keep the OpenTelemetry model as the on-disk layout, so there is nothing to convert. Blocks are Arrow IPC, aligned so a query reads the buffers straight out of the page cache, and the block directory is the only state there is, which is what removes the coordination layer. One binary then serves a query API, an MCP endpoint, a browser UI and a terminal UI off that same read path.',
  },
  {
    url: 'https://github.com/TrianaLab/pacto',
    website: 'https://pacto.run',
    websiteLabel: 'Visit the Pacto website',
    eyebrow: 'Featured work',
    thesis:
      'Pacto is to agents what Internal Developer Platforms are to humans. That is the direction, and a good part of it is still an open draft.',
    problem:
      'What a service needs to run, and what it promises back, is rarely written down anywhere a machine can read it. It lives in Helm values and in whoever was on call last quarter, and it goes stale as soon as the code or the cluster moves. Ask what a change will break and someone has to go and find out.',
    approach:
      'Write the contract once and push it to the same OCI registry as the image. The operator checks running workloads against it and reports what has drifted. I have a draft PR open that would take this a layer up, into a graph over services, contract revisions and deployed targets, so you can ask what a change reaches before it ships. That part is not merged yet.',
  },
];

/**
 * Projects surfaced as "Selected systems", in display order. Every entry is a
 * canonical `volunteer[].url` from the résumé document. The FEATURED urls are
 * shown separately above this list and are deliberately excluded here.
 */
export const SELECTED_SYSTEM_URLS = [
  'https://edudiaz.dev',
  'https://github.com/TrianaLab/remake',
] as const;

/**
 * Editorial "how I work" statements. These are opinions about practice, not
 * résumé facts, so they belong here rather than in the canonical document.
 */
export const PRINCIPLES: { title: string; body: string }[] = [
  {
    title: 'Rules a machine can check',
    body: 'Conventions that only live in a README go stale and nobody notices for a while. That is most of what Pacto is about: the check fails instead of a reviewer catching it.',
  },
  {
    title: 'Nothing important kept in someone’s head',
    body: 'If I had to remember a specific trick to get a deploy through, it goes into the CLI or the operator afterwards.',
  },
  {
    title: 'Boring underneath',
    body: 'The unusual part should be the thing engineers type into. Below that I would rather run stock components other people have already debugged.',
  },
];

/** Suggested commands shown in the terminal welcome and on its entry screen. */
export const TERMINAL_SUGGESTIONS: { command: string; hint: string }[] = [
  { command: 'kubectl get all', hint: 'every résumé resource at once' },
  { command: 'kubectl get work', hint: 'employment history as a table' },
  { command: 'kubectl describe volunteer volunteer-trianalab-pacto', hint: 'drill into one project' },
  { command: 'kubectl get resume main-resume -o yaml', hint: 'the raw JSON Resume document' },
];

export const SEO = {
  /** Runtime <title>: {name} and {role} are filled from the canonical résumé. */
  titleTemplate: '{name} · {role}',
} as const;
