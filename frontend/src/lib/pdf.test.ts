import { describe, it, expect } from 'vitest';
import { buildResumeDocDef } from './pdf';
import type { Resume } from './schema';

const sample: Resume = {
  basics: {
    name: 'Ada Lovelace',
    label: 'Mathematician',
    email: 'ada@example.com',
    phone: '+44 0',
    url: 'https://example.com',
    location: { city: 'London', region: 'UK' },
    profiles: [{ network: 'GitHub', username: 'ada', url: 'https://github.com/ada' }],
    summary: 'First programmer.',
  },
  work: [
    {
      name: 'Babbage Engines',
      position: 'Analyst',
      startDate: '1843-01-01',
      endDate: '1850-01-01',
      summary: 'Wrote the first algorithm.',
      highlights: ['Note G'],
    },
  ],
  education: [
    {
      institution: 'Self-taught',
      area: 'Mathematics',
      studyType: 'Independent',
      startDate: '1820-01-01',
    },
  ],
  skills: [{ name: 'Algorithms', level: 'Expert', keywords: ['recursion', 'tables'] }],
  languages: [{ language: 'English', fluency: 'Native' }],
};

describe('buildResumeDocDef', () => {
  it('sets the document metadata from basics', () => {
    const doc = buildResumeDocDef(sample);
    expect(doc.info?.author).toBe('Ada Lovelace');
    expect(doc.info?.subject).toBe('Mathematician');
  });

  it('paginates as A4 with consistent margins', () => {
    const doc = buildResumeDocDef(sample);
    expect(doc.pageSize).toBe('A4');
    expect(doc.pageMargins).toEqual([40, 50, 40, 50]);
  });

  it('renders the name and label as the first content blocks', () => {
    const doc = buildResumeDocDef(sample);
    const content = doc.content as { text?: string; style?: string }[];
    expect(content[0]?.text).toBe('Ada Lovelace');
    expect(content[0]?.style).toBe('name');
    expect(content[1]?.text).toBe('Mathematician');
    expect(content[1]?.style).toBe('label');
  });

  it('includes every populated section title', () => {
    const doc = buildResumeDocDef(sample);
    const titles = (doc.content as { text?: string; style?: string }[])
      .filter((c) => c.style === 'sectionTitle')
      .map((c) => c.text);
    expect(titles).toContain('EXPERIENCE');
    expect(titles).toContain('EDUCATION');
    expect(titles).toContain('SKILLS');
    expect(titles).toContain('LANGUAGES');
    expect(titles).not.toContain('PROJECTS'); // not populated
    expect(titles).not.toContain('CERTIFICATES');
  });

  it('handles an empty resume without crashing', () => {
    const doc = buildResumeDocDef({});
    expect(doc.pageSize).toBe('A4');
    expect((doc.content as unknown[]).length).toBeGreaterThan(0);
  });

  // ── ATS-friendliness assertions ────────────────────────────────────
  // These guard the choices made for ATS parsers: no decorated section
  // titles, no multi-column rows in scannable sections, no non-ASCII
  // date separators, plain " | " contact joins.

  it('section titles have no decoration (decorated text confuses ATS)', () => {
    const doc = buildResumeDocDef(sample);
    const sectionTitleStyle = doc.styles?.sectionTitle as Record<string, unknown> | undefined;
    expect(sectionTitleStyle).toBeDefined();
    expect(sectionTitleStyle?.decoration).toBeUndefined();
  });

  it('skills render as a single linear text run, not a two-column row', () => {
    const doc = buildResumeDocDef(sample);
    const content = doc.content as unknown[];
    const titles = content
      .map((c) => (c as { text?: string; style?: string }).text)
      .map((t, i) => ({ t, i }))
      .filter((x) => x.t === 'SKILLS');
    expect(titles).toHaveLength(1);
    // The next entry after the SKILLS title is the first skill row.
    const skillsRow = content[titles[0].i + 1] as Record<string, unknown>;
    expect(skillsRow.columns).toBeUndefined();
    expect(Array.isArray(skillsRow.text)).toBe(true);
  });

  it('date ranges use ASCII hyphen, not the "→" arrow', () => {
    const doc = buildResumeDocDef(sample);
    const flat = JSON.stringify(doc.content);
    expect(flat).not.toContain('→');
    expect(flat).toMatch(/\d{4}.+ - .+\d{4}|\d{4}.+ - Present/);
  });

  it('contact line uses " | " as a separator (parsed by ATS readers)', () => {
    const doc = buildResumeDocDef(sample);
    const content = doc.content as { text?: string; style?: string }[];
    const contact = content.find((c) => c.style === 'contact' && c.text?.includes('ada@example.com'));
    expect(contact?.text).toContain(' | ');
  });

  it('profile lines include the full URL, not just the username', () => {
    const doc = buildResumeDocDef(sample);
    const content = doc.content as { text?: string; style?: string }[];
    const profiles = content.find((c) => c.style === 'contact' && c.text?.includes('GitHub'));
    expect(profiles?.text).toContain('https://github.com/ada');
  });
});
