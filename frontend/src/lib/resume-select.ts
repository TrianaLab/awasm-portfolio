// View-model selectors: they combine the canonical JSON Resume document with
// the presentation config in portfolio.ts. Pure functions, no résumé facts of
// their own — components render what these return.

import type { Resume, Volunteer, Work } from './schema';
import { FEATURED, SELECTED_SYSTEM_URLS, SEO } from './portfolio';

/** Looks an entry up by its canonical URL — the stable key across the doc. */
export function byUrl<T extends { url?: string }>(items: T[] | undefined, url: string): T | undefined {
  return (items ?? []).find((item) => item.url === url);
}

/** The most recent role. Used for `worksFor`, not for the headline. */
export function currentRole(resume: Resume | null): Work | null {
  return resume?.work?.[0] ?? null;
}

/**
 * The headline title is `basics.label`, not the current job title: the
 * standing professional identity outlives any one employer, and the most
 * recent position is already visible in the experience timeline.
 */
export function headlineTitle(resume: Resume | null): string | undefined {
  return resume?.basics?.label ?? currentRole(resume)?.position ?? undefined;
}

export function pageTitle(resume: Resume | null, fallback: string): string {
  const name = resume?.basics?.name;
  const role = headlineTitle(resume);
  if (!name) return fallback;
  if (!role) return name;
  return SEO.titleTemplate.replace('{name}', name).replace('{role}', role);
}

export function featuredProject(resume: Resume | null): Volunteer | undefined {
  return byUrl(resume?.volunteer, FEATURED.url);
}

/** Selected systems, in the order configured, minus any that no longer exist. */
export function selectedSystems(resume: Resume | null): Volunteer[] {
  return SELECTED_SYSTEM_URLS.map((url) => byUrl(resume?.volunteer, url)).filter(
    (v): v is Volunteer => v !== undefined,
  );
}

/** Everything else in `volunteer` — the compact upstream-contributions list. */
export function upstreamContributions(resume: Resume | null): Volunteer[] {
  const featuredUrls = new Set<string>([FEATURED.url, ...SELECTED_SYSTEM_URLS]);
  return (resume?.volunteer ?? []).filter((v) => !v.url || !featuredUrls.has(v.url));
}

export type CompanyGroup = {
  name: string;
  url?: string;
  startDate?: string;
  endDate?: string;
  roles: Work[];
};

/**
 * Groups consecutive `work` entries by employer so four Appian roles read as
 * one company timeline. Order is preserved; only the presentation changes.
 */
export function experienceGroups(resume: Resume | null): CompanyGroup[] {
  const groups: CompanyGroup[] = [];
  for (const role of resume?.work ?? []) {
    const name = role.name ?? '';
    const last = groups[groups.length - 1];
    if (last && last.name === name) {
      last.roles.push(role);
      last.startDate = role.startDate ?? last.startDate;
    } else {
      groups.push({
        name,
        url: role.url,
        startDate: role.startDate,
        endDate: role.endDate,
        roles: [role],
      });
    }
  }
  return groups;
}

/** Short organisation label: "TrianaLab: pacto" → "pacto". */
export function shortOrgName(organization?: string): string {
  if (!organization) return '';
  const colon = organization.indexOf(':');
  return colon >= 0 ? organization.slice(colon + 1).trim() : organization;
}

export function formatYear(date?: string): string {
  if (!date) return '';
  const parsed = new Date(date);
  return Number.isNaN(parsed.getTime()) ? date : String(parsed.getUTCFullYear());
}

/**
 * schema.org Person built from the canonical document — structured data with
 * no second copy of the facts living in markup.
 */
export function personJsonLd(resume: Resume | null): string | null {
  const basics = resume?.basics;
  if (!basics?.name) return null;
  return JSON.stringify({
    '@context': 'https://schema.org',
    '@type': 'Person',
    name: basics.name,
    jobTitle: headlineTitle(resume),
    description: basics.summary,
    email: basics.email ? `mailto:${basics.email}` : undefined,
    url: basics.url,
    address: basics.location?.city
      ? {
          '@type': 'PostalAddress',
          addressLocality: basics.location.city,
          addressRegion: basics.location.region,
          addressCountry: basics.location.countryCode,
        }
      : undefined,
    worksFor: currentRole(resume)?.name
      ? { '@type': 'Organization', name: currentRole(resume)?.name, url: currentRole(resume)?.url }
      : undefined,
    sameAs: (basics.profiles ?? []).map((p) => p.url).filter(Boolean),
    knowsAbout: (resume?.skills ?? []).flatMap((s) => s.keywords ?? []),
  });
}
