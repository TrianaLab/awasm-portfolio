import { describe, it, expect } from 'vitest';
import { parseHash } from './portfolio';
import {
  byUrl,
  currentRole,
  experienceGroups,
  formatYear,
  pageTitle,
  personJsonLd,
  shortOrgName,
} from './resume-select';
import type { Resume } from './schema';

const sample: Resume = {
  basics: {
    name: 'Ada Lovelace',
    label: 'Mathematician',
    email: 'ada@example.com',
    url: 'https://example.com',
    summary: 'First programmer.',
    location: { city: 'London', region: 'England', countryCode: 'GB' },
    profiles: [{ network: 'GitHub', username: 'ada', url: 'https://github.com/ada' }],
  },
  work: [
    { name: 'Analytical Co', position: 'Principal', url: 'https://analytical.example', startDate: '2024-01-01' },
    { name: 'Analytical Co', position: 'Senior', url: 'https://analytical.example', startDate: '2022-01-01', endDate: '2024-01-01' },
    { name: 'Babbage Engines', position: 'Analyst', startDate: '1843-01-01', endDate: '1850-01-01' },
  ],
  volunteer: [
    { organization: 'Org: alpha', url: 'https://example.com/alpha', summary: 'a' },
    { organization: 'beta', url: 'https://example.com/beta', summary: 'b' },
  ],
  skills: [{ name: 'Maths', keywords: ['algebra'] }],
};

describe('parseHash', () => {
  it.each([
    ['', 'home', null],
    ['#/', 'home', null],
    ['#/resume', 'resume', null],
    ['#/terminal', 'terminal', null],
    ['#/work', 'home', 'work'],
    ['#/experience', 'home', 'experience'],
    ['#/about', 'home', 'about'],
    ['#/nonsense', 'home', null],
  ])('maps %s to %s/%s', (hash, view, section) => {
    expect(parseHash(hash)).toEqual({ view, section });
  });

  it('ignores case and a trailing slash', () => {
    expect(parseHash('#/Resume/')).toEqual({ view: 'resume', section: null });
  });
});

describe('resume selectors', () => {
  it('looks entries up by their canonical url', () => {
    expect(byUrl(sample.volunteer, 'https://example.com/beta')?.organization).toBe('beta');
    expect(byUrl(sample.volunteer, 'https://example.com/missing')).toBeUndefined();
  });

  it('takes the current role from the first work entry', () => {
    expect(currentRole(sample)?.position).toBe('Principal');
    expect(currentRole(null)).toBeNull();
  });

  it('groups consecutive roles at the same employer into one company block', () => {
    const groups = experienceGroups(sample);
    expect(groups.map((g) => g.name)).toEqual(['Analytical Co', 'Babbage Engines']);
    expect(groups[0].roles.map((r) => r.position)).toEqual(['Principal', 'Senior']);
    // The group spans the earliest start of its roles through the latest end.
    expect(groups[0].startDate).toBe('2022-01-01');
    expect(groups[0].roles[0].endDate).toBeUndefined();
  });

  it('strips the org prefix from project names', () => {
    expect(shortOrgName('TrianaLab: pacto')).toBe('pacto');
    expect(shortOrgName('Docker')).toBe('Docker');
    expect(shortOrgName(undefined)).toBe('');
  });

  it('formats years without shifting across the timezone boundary', () => {
    expect(formatYear('2026-01-01')).toBe('2026');
    expect(formatYear(undefined)).toBe('');
  });

  it('builds the document title from the canonical name and current role', () => {
    expect(pageTitle(sample, 'fallback')).toBe('Ada Lovelace — Principal');
    expect(pageTitle(null, 'fallback')).toBe('fallback');
  });

  it('derives schema.org structured data from the document, not from markup', () => {
    const jsonLd = JSON.parse(personJsonLd(sample) ?? '{}');
    expect(jsonLd['@type']).toBe('Person');
    expect(jsonLd.name).toBe('Ada Lovelace');
    expect(jsonLd.jobTitle).toBe('Principal');
    expect(jsonLd.sameAs).toContain('https://github.com/ada');
    expect(personJsonLd(null)).toBeNull();
  });
});
