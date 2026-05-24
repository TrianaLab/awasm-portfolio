import { test, expect } from '@playwright/test';
import { execFileSync } from 'node:child_process';
import { mkdtempSync, writeFileSync, existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

// Asserts the actual rendered PDF — as an ATS reader sees it — keeps
// the resume content in a parseable order. The browser produces the
// PDF, playwright captures the download, then we shell out to pypdf
// for the text extraction (the same approach an ATS pipeline uses).
//
// pypdf is checked at startup; missing tooling skips the spec rather
// than failing the suite, so developers without the Python dep can
// still run the rest of e2e.

function extractPdfText(pdfPath: string): string | null {
  try {
    return execFileSync(
      'python3',
      ['-c', `import sys, pypdf; r = pypdf.PdfReader(sys.argv[1]); print('\\n'.join(p.extract_text() for p in r.pages))`, pdfPath],
      { encoding: 'utf8' },
    );
  } catch {
    return null;
  }
}

test.describe('PDF ATS extraction', () => {
  test('downloaded PDF extracts to clean linear text in expected order', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });

    const [download] = await Promise.all([
      page.waitForEvent('download', { timeout: 30_000 }),
      page.getByRole('button', { name: /download resume as pdf/i }).click(),
    ]);

    const dir = mkdtempSync(join(tmpdir(), 'resume-pdf-'));
    const path = join(dir, 'resume.pdf');
    await download.saveAs(path);
    expect(existsSync(path), 'PDF must be written to disk').toBe(true);

    const text = extractPdfText(path);
    if (text === null) {
      test.skip(true, 'pypdf not installed (pip install pypdf) — skipping ATS extraction check');
      return;
    }

    // Write the extracted text alongside the PDF so a failing run can be
    // inspected from the test output directory.
    writeFileSync(join(dir, 'extracted.txt'), text);

    // ── ATS-critical fields ─────────────────────────────────────────
    // Name and label must appear first (every ATS uses them to seed the
    // candidate record).
    expect(text).toMatch(/Eduardo D[íi]az/);
    expect(text).toContain('Platform Engineer');

    // Contact basics must extract as raw, parseable strings — not
    // glyph-mashed clusters.
    expect(text).toContain('edudiazasencio@gmail.com');
    expect(text).toContain('+34 622287557');
    expect(text).toContain('edudiaz.dev');

    // Section ordering — Experience must come before Education must
    // come before Skills. ATSs that bucket by section rely on this.
    const idxExperience = text.indexOf('EXPERIENCE');
    const idxEducation = text.indexOf('EDUCATION');
    const idxSkills = text.indexOf('SKILLS');
    expect(idxExperience).toBeGreaterThan(-1);
    expect(idxEducation).toBeGreaterThan(idxExperience);
    expect(idxSkills).toBeGreaterThan(idxEducation);

    // Dates must extract with the ASCII hyphen — the "→" arrow gets
    // dropped by some readers and collapses ranges to "Feb 2024Jul 2024".
    expect(text).not.toContain('→');
    expect(text).toMatch(/\b\w{3,9}\s+\d{4}\s+-\s+(Present|\w{3,9}\s+\d{4})/);

    // The skills section must not interleave a skill's name with another
    // skill's keywords — i.e. the keyword line must immediately follow
    // its label. Checking one representative pair: Kubernetes belongs
    // under "Kubernetes and Cloud-Native".
    const labelIdx = text.indexOf('Kubernetes and Cloud-Native');
    expect(labelIdx).toBeGreaterThan(-1);
    const k8sIdx = text.indexOf('Kubernetes', labelIdx + 'Kubernetes and Cloud-Native'.length);
    expect(k8sIdx).toBeGreaterThan(labelIdx);
    // No other "Skill: " header should sneak between the label and its keywords.
    const interloper = text.slice(labelIdx + 'Kubernetes and Cloud-Native'.length, k8sIdx);
    expect(interloper).not.toMatch(/Cloud Platforms|Observability|Platform Engineering/);
  });
});
