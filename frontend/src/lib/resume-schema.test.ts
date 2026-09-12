// JSON Resume compliance gate.
//
// The canonical document lives in the Go backend and is produced by the same
// kubectl command the browser runs. This suite regenerates it and validates it
// against the official @jsonresume/schema package, so a change that would
// break standard JSON Resume tooling or third-party themes fails CI.

import { describe, it, expect, beforeAll } from 'vitest';
import { execFileSync } from 'node:child_process';
import { validate, schema } from '@jsonresume/schema';
import { FEATURED, SELECTED_SYSTEM_URLS, TERMINAL_SUGGESTIONS } from './portfolio';
import { RESUME_QUERY } from './wasm';
import { buildResumeDocDef } from './pdf';
import {
  experienceGroups,
  featuredProjects,
  headlineTitle,
  selectedSystems,
} from './resume-select';
import type { Resume } from './schema';

const REPO_ROOT = new URL('../../..', import.meta.url).pathname;

function loadCanonicalResume(): Resume & Record<string, unknown> {
  const out = execFileSync('go', ['run', 'cli.go', ...RESUME_QUERY.split(' ').slice(1)], {
    cwd: REPO_ROOT,
    encoding: 'utf8',
    maxBuffer: 8 * 1024 * 1024,
  });
  const parsed = JSON.parse(out) as unknown;
  return (Array.isArray(parsed) ? parsed[0] : parsed) as Resume & Record<string, unknown>;
}

function runCli(command: string): string {
  return execFileSync('go', ['run', 'cli.go', ...command.split(' ').slice(1)], {
    cwd: REPO_ROOT,
    encoding: 'utf8',
    maxBuffer: 8 * 1024 * 1024,
  });
}

function validationErrors(doc: unknown): unknown[] {
  let result: unknown[] = [];
  validate(doc, (errors) => {
    result = errors ?? [];
  });
  return result;
}

describe('canonical JSON Resume document', () => {
  let doc: Resume & Record<string, unknown>;

  beforeAll(() => {
    doc = loadCanonicalResume();
  }, 300_000);

  it('validates against the official JSON Resume schema with no preprocessing', () => {
    expect(validationErrors(doc)).toEqual([]);
  });

  it('carries no proprietary root-level properties', () => {
    const official = new Set(Object.keys(schema.properties));
    expect(Object.keys(doc).filter((key) => !official.has(key))).toEqual([]);
  });

  it('fails loudly when the document stops matching the schema', () => {
    const broken = { ...doc, work: { position: 'not an array' } };
    expect(validationErrors(broken).length).toBeGreaterThan(0);
  });

  it('feeds the PDF and ATS pipeline from the same document', () => {
    const def = buildResumeDocDef(doc);
    const text = JSON.stringify(def.content);
    expect(text).toContain(doc.basics?.name);
    expect(text).toContain('EXPERIENCE');
    expect(text).toContain(doc.work?.[0]?.position);
  });

  it('resolves every résumé reference held by the presentation config', () => {
    expect(featuredProjects(doc).map((p) => p.config.url)).toEqual(FEATURED.map((f) => f.url));
    expect(selectedSystems(doc)).toHaveLength(SELECTED_SYSTEM_URLS.length);
  });

  // The suggestions are printed in the terminal welcome banner and offered as
  // copy-to-clipboard chips, so a stale resource name hands a visitor an error
  // as their first interaction.
  it.each(TERMINAL_SUGGESTIONS.map((s) => s.command))('suggested command runs clean: %s', (command) => {
    expect(runCli(command).trim()).not.toMatch(/^Error:/);
  }, 300_000);

  it('supplies everything the home-page view models need', () => {
    // The hero headline reads basics.label, so the document must carry one.
    expect(headlineTitle(doc)).toBeTruthy();
    expect(experienceGroups(doc).length).toBeGreaterThan(0);
    for (const { config, entry } of featuredProjects(doc)) {
      expect(entry.summary, `${config.url} has no summary`).toBeTruthy();
    }
  });
});
