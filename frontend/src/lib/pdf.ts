// Browser-side PDF generation. Maps the JSON Resume structure to a
// pdfmake document definition. The output is vector text (selectable,
// searchable, ATS-parseable) — no rasterisation, no HTML capture.
//
// Layout: A4 portrait, single column, classic resume.

import type { TDocumentDefinitions, Content, StyleDictionary } from 'pdfmake/interfaces';
import type { Resume } from './schema';

export function buildResumeDocDef(resume: Resume): TDocumentDefinitions {
  const content: Content[] = [];
  pushBasics(content, resume);
  pushSection(content, 'Experience', buildWork(resume));
  pushSection(content, 'Education', buildEducation(resume));
  pushSection(content, 'Open source & volunteering', buildVolunteer(resume));
  pushSection(content, 'Projects', buildProjects(resume));
  pushSection(content, 'Certificates', buildCertificates(resume));
  pushSection(content, 'Awards', buildAwards(resume));
  pushSection(content, 'Publications', buildPublications(resume));
  pushSection(content, 'Skills', buildSkills(resume));
  pushSection(content, 'Languages', buildLanguages(resume));
  pushSection(content, 'Interests', buildInterests(resume));

  return {
    info: {
      title: `${resume.basics?.name ?? 'Resume'} — Resume`,
      author: resume.basics?.name,
      subject: resume.basics?.label,
    },
    pageSize: 'A4',
    pageMargins: [40, 50, 40, 50],
    defaultStyle: { fontSize: 10, lineHeight: 1.35, color: '#2c3e50' },
    styles: stylesheet,
    content,
    footer: (current, total) => ({
      columns: [
        { text: resume.basics?.name ?? '', alignment: 'left', style: 'footer' },
        { text: `${current} / ${total}`, alignment: 'right', style: 'footer' },
      ],
      margin: [40, 0, 40, 20],
    }),
  };
}

const stylesheet: StyleDictionary = {
  name: { fontSize: 22, bold: true, color: '#1a202c' },
  label: { fontSize: 12, color: '#0969da', margin: [0, 2, 0, 8] },
  contact: { fontSize: 9, color: '#59636e', margin: [0, 0, 0, 6] },
  summary: { fontSize: 10, color: '#2c3e50', margin: [0, 4, 0, 0] },
  sectionTitle: {
    fontSize: 11,
    bold: true,
    color: '#1a202c',
    margin: [0, 14, 0, 6],
    decoration: 'underline',
    decorationColor: '#d0d7de',
  },
  entryTitle: { fontSize: 10.5, bold: true, color: '#1a202c' },
  entrySubtitle: { fontSize: 10, color: '#0969da' },
  entryDates: { fontSize: 9, italics: true, color: '#59636e' },
  entrySummary: { fontSize: 10, color: '#2c3e50', margin: [0, 2, 0, 0] },
  footer: { fontSize: 8, color: '#59636e' },
  tag: { fontSize: 9, color: '#0969da' },
};

function pushBasics(content: Content[], resume: Resume): void {
  const b = resume.basics ?? {};
  content.push({ text: b.name ?? '', style: 'name' });
  if (b.label) content.push({ text: b.label, style: 'label' });

  const contactParts: string[] = [];
  if (b.email) contactParts.push(b.email);
  if (b.phone) contactParts.push(b.phone);
  if (b.url) contactParts.push(b.url);
  if (b.location?.city) {
    const region = b.location.region ? `, ${b.location.region}` : '';
    contactParts.push(`${b.location.city}${region}`);
  }
  if (contactParts.length > 0) {
    content.push({ text: contactParts.join('  ·  '), style: 'contact' });
  }

  if (b.profiles && b.profiles.length > 0) {
    const profileLine = b.profiles.map((p) => `${p.network}: ${p.username}`).join('  ·  ');
    content.push({ text: profileLine, style: 'contact' });
  }

  if (b.summary) content.push({ text: b.summary, style: 'summary' });
}

function pushSection(content: Content[], title: string, body: Content[]): void {
  if (body.length === 0) return;
  content.push({ text: title.toUpperCase(), style: 'sectionTitle' });
  content.push(...body);
}

function buildWork(resume: Resume): Content[] {
  return (resume.work ?? []).map((w) =>
    entry({
      title: w.position,
      subtitle: w.name,
      dates: formatRange(w.startDate, w.endDate),
      summary: w.summary,
      highlights: w.highlights,
    }),
  );
}

function buildEducation(resume: Resume): Content[] {
  return (resume.education ?? []).map((e) =>
    entry({
      title: e.studyType ? `${e.studyType} · ${e.area ?? ''}` : (e.area ?? ''),
      subtitle: e.institution,
      dates: formatRange(e.startDate, e.endDate),
      summary: e.score,
      highlights: e.courses,
    }),
  );
}

function buildVolunteer(resume: Resume): Content[] {
  return (resume.volunteer ?? []).map((v) =>
    entry({
      title: v.position,
      subtitle: v.organization,
      dates: formatRange(v.startDate, v.endDate),
      summary: v.summary,
      highlights: v.highlights,
    }),
  );
}

function buildProjects(resume: Resume): Content[] {
  return (resume.projects ?? []).map((p) =>
    entry({
      title: p.name,
      dates: formatRange(p.startDate, p.endDate),
      summary: p.description,
      highlights: p.highlights,
    }),
  );
}

function buildCertificates(resume: Resume): Content[] {
  return (resume.certificates ?? []).map((c) =>
    entry({
      title: c.name,
      subtitle: c.issuer,
      dates: formatDate(c.date),
    }),
  );
}

function buildAwards(resume: Resume): Content[] {
  return (resume.awards ?? []).map((a) =>
    entry({
      title: a.title,
      subtitle: a.awarder,
      dates: formatDate(a.date),
      summary: a.summary,
    }),
  );
}

function buildPublications(resume: Resume): Content[] {
  return (resume.publications ?? []).map((p) =>
    entry({
      title: p.name,
      subtitle: p.publisher,
      dates: formatDate(p.releaseDate),
      summary: p.summary,
    }),
  );
}

function buildSkills(resume: Resume): Content[] {
  return (resume.skills ?? []).map((s) => ({
    columns: [
      {
        width: 'auto',
        text: [
          { text: s.name ?? '', style: 'entryTitle' },
          s.level ? { text: `  (${s.level})`, style: 'entrySubtitle' } : '',
        ],
      },
      {
        width: '*',
        text: (s.keywords ?? []).join(', '),
        style: 'tag',
        alignment: 'right',
      },
    ],
    margin: [0, 0, 0, 4],
  }));
}

function buildLanguages(resume: Resume): Content[] {
  return (resume.languages ?? []).map((l) => ({
    columns: [
      { width: 'auto', text: l.language ?? '', style: 'entryTitle' },
      { width: '*', text: l.fluency ?? '', style: 'entrySubtitle', alignment: 'right' },
    ],
    margin: [0, 0, 0, 4],
  }));
}

function buildInterests(resume: Resume): Content[] {
  return (resume.interests ?? []).map((i) => ({
    columns: [
      { width: 'auto', text: i.name ?? '', style: 'entryTitle' },
      { width: '*', text: (i.keywords ?? []).join(', '), style: 'tag', alignment: 'right' },
    ],
    margin: [0, 0, 0, 4],
  }));
}

function entry(opts: {
  title?: string;
  subtitle?: string;
  dates?: string;
  summary?: string;
  highlights?: string[];
}): Content {
  const stack: Content[] = [];
  stack.push({
    columns: [
      {
        width: '*',
        stack: [
          { text: opts.title ?? '', style: 'entryTitle' },
          opts.subtitle ? { text: opts.subtitle, style: 'entrySubtitle' } : '',
        ],
      },
      opts.dates ? { width: 'auto', text: opts.dates, style: 'entryDates' } : { width: 'auto', text: '' },
    ],
  });
  if (opts.summary) {
    stack.push({ text: opts.summary, style: 'entrySummary' });
  }
  if (opts.highlights && opts.highlights.length > 0) {
    stack.push({
      ul: opts.highlights,
      style: 'entrySummary',
      margin: [12, 2, 0, 0],
    });
  }
  return { stack, margin: [0, 0, 0, 8] };
}

function formatRange(start?: string, end?: string): string {
  if (!start && !end) return '';
  return `${formatDate(start)} → ${end ? formatDate(end) : 'Present'}`;
}

function formatDate(d?: string): string {
  if (!d) return '';
  const parsed = new Date(d);
  if (Number.isNaN(parsed.getTime())) return d;
  return parsed.toLocaleDateString('en-US', { year: 'numeric', month: 'short' });
}

/**
 * Generates the PDF in the browser and triggers a download as
 * `<basics.name>.pdf` (or resume.pdf if no name).
 */
export async function downloadResumePdf(resume: Resume): Promise<void> {
  // pdfmake is ~1.4 MB so it's lazy-imported on download. The companion
  // vfs_fonts module ships the embedded Roboto fonts that pdfmake uses
  // by default.
  const pdfMakeMod = await import('pdfmake/build/pdfmake');
  const vfsMod = await import('pdfmake/build/vfs_fonts');

  // pdfmake's vfs export shape differs across versions/bundlers:
  //   - newer CJS:  the module IS the vfs object
  //   - older:      { vfs: {...} }
  //   - legacy:     { pdfMake: { vfs: {...} } }
  // Search every plausible location for the dict containing Roboto.
  const looksLikeVfs = (v: unknown): v is Record<string, string> =>
    typeof v === 'object' && v !== null && 'Roboto-Regular.ttf' in v;

  const candidates: unknown[] = [
    (vfsMod as { default?: unknown }).default,
    (vfsMod as { vfs?: unknown }).vfs,
    (vfsMod as { pdfMake?: { vfs?: unknown } }).pdfMake?.vfs,
    (vfsMod as { default?: { vfs?: unknown } }).default?.vfs,
    (vfsMod as { default?: { pdfMake?: { vfs?: unknown } } }).default?.pdfMake?.vfs,
    vfsMod,
  ];
  const vfs = candidates.find(looksLikeVfs);
  if (!vfs) {
    throw new Error('pdfmake: failed to locate vfs_fonts dictionary');
  }

  type PdfMake = { createPdf: (d: TDocumentDefinitions) => { download: (f: string) => void }; vfs?: Record<string, string> };
  const maybeWrapped = pdfMakeMod as unknown as { default?: PdfMake };
  const pdfMake: PdfMake = maybeWrapped.default ?? (pdfMakeMod as unknown as PdfMake);
  pdfMake.vfs = vfs;

  const docDef = buildResumeDocDef(resume);
  const filename = (resume.basics?.name ?? 'resume').replace(/\s+/g, '-').toLowerCase() + '.pdf';
  pdfMake.createPdf(docDef).download(filename);
}
