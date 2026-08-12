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
  thesis: 'Pacto is to service operations what OpenAPI is to HTTP APIs.',
  problem:
    'A service’s operational contract — what it needs, what it promises, how it scales, how it fails — is folklore. It lives in runbooks, Helm values and the heads of the people who were on call last quarter, and it drifts away from the artifact the moment either one changes.',
  approach:
    'Describe the contract once, ship it next to the image in the same OCI registry, and make the platform verify running workloads against it. Deployment stops being a pile of YAML and becomes an assertion the cluster can check.',
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

/** Employers whose entries carry the most weight on the home page. */
export const HEADLINE_EMPLOYER_URL = 'https://emergence.ai';

/**
 * Editorial "how I work" statements. These are opinions about practice, not
 * résumé facts, so they belong here rather than in the canonical document.
 */
export const PRINCIPLES: { title: string; body: string }[] = [
  {
    title: 'Contracts over conventions',
    body: 'A golden path that is only documented is a suggestion. Encode it so the platform can check it, and the interface becomes the guardrail.',
  },
  {
    title: 'Operational knowledge belongs in tooling',
    body: 'Anything a senior engineer has to remember at 3am is a defect in the interface. Push it into the CLI, the operator, the admission check.',
  },
  {
    title: 'Boring where it counts',
    body: 'Interesting infrastructure pages you. I spend novelty budget on the developer interface and keep everything underneath dull and observable.',
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
  titleTemplate: '{name} — {role}',
  tagline:
    'Platform engineer building contract-driven deployment systems for cloud-native infrastructure.',
} as const;
