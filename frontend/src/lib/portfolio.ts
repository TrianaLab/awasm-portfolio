// Portfolio presentation model.
//
// Boundary: resume.json (JSON Resume) owns every factual claim — biography,
// employment, education, skills, project descriptions, dates. This file owns
// presentation and editorial emphasis only: navigation, ordering, CTA copy,
// headings, terminal suggestions and SEO strings.
//
// Résumé entries are referenced by their canonical `url` (a stable key that
// already exists in the document) and never copied. If a URL below stops
// resolving, resume-select.test.ts fails — that is the intended guard rail.

export type View = 'home' | 'resume' | 'terminal';

/** Home-page section ids that double as routes (`#/work` → home + scroll). */
export const HOME_SECTIONS = ['work', 'experience', 'about'] as const;

export type Location = { view: View; section: string | null };

/** Parses `window.location.hash` into a view + optional home section. */
export function parseHash(hash: string): Location {
  const path = hash.replace(/^#\/?/, '').replace(/\/+$/, '').toLowerCase();
  if (path === 'resume') return { view: 'resume', section: null };
  if (path === 'terminal') return { view: 'terminal', section: null };
  if ((HOME_SECTIONS as readonly string[]).includes(path)) {
    return { view: 'home', section: path };
  }
  return { view: 'home', section: null };
}

export const NAV: { href: string; label: string }[] = [
  { href: '#/work', label: 'Work' },
  { href: '#/experience', label: 'Experience' },
  { href: '#/about', label: 'About' },
  { href: '#/resume', label: 'Résumé' },
  { href: '#/terminal', label: 'Terminal' },
];

export const CTA = {
  work: { href: '#/work', label: 'Explore my work' },
  terminal: { href: '#/terminal', label: 'Open the terminal' },
  resume: { label: 'Download résumé' },
} as const;

/**
 * The one project the home page is built around. Referenced by the canonical
 * URL of its JSON Resume entry; the copy below is editorial framing, not a
 * restatement of the entry's own description.
 */
export const FEATURED = {
  url: 'https://github.com/TrianaLab/pacto',
  eyebrow: 'Featured work',
  thesis:
    'Pacto is to agents what Internal Developer Platforms are to humans. That is the direction, and a good part of it is still an open draft.',
  problem:
    'What a service needs to run, and what it promises back, is rarely written down anywhere a machine can read it. It lives in Helm values and in whoever was on call last quarter, and it goes stale as soon as the code or the cluster moves. Ask what a change will break and someone has to go and find out.',
  approach:
    'Write the contract once and push it to the same OCI registry as the image. The operator checks running workloads against it and reports what has drifted. I have a draft PR open that would take this a layer up, into a graph over services, contract revisions and deployed targets, so you can ask what a change reaches before it ships. That part is not merged yet.',
} as const;

/**
 * Projects surfaced as "Selected systems", in display order. Every entry is a
 * canonical `volunteer[].url` from the résumé document. FEATURED.url is shown
 * separately above this list and is deliberately excluded here.
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
  { command: 'kubectl describe volunteer pacto', hint: 'drill into one project' },
  { command: 'kubectl get resume main-resume -o yaml', hint: 'the raw JSON Resume document' },
];

export const SEO = {
  /** Runtime <title>: {name} and {role} are filled from the canonical résumé. */
  titleTemplate: '{name} · {role}',
} as const;
