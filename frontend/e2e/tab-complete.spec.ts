import { test, expect } from '@playwright/test';

// Cross-cutting helpers ─ all tab-completion specs hit the same first
// terminal and need a focused helper textarea to send keystrokes into.
async function bootTerminal(page: import('@playwright/test').Page) {
  await page.goto('/');
  await expect(page.locator('.xterm')).toContainText('Welcome', { timeout: 10_000 });
  const helper = page.locator('.xterm-helper-textarea').first();
  await helper.focus();
  return helper;
}

async function readVisibleTerminal(page: import('@playwright/test').Page): Promise<string> {
  return page.evaluate(() => {
    const rows = Array.from(document.querySelectorAll('.xterm-rows > div'));
    return rows.map((r) => (r as HTMLElement).innerText).join('\n');
  });
}

test.describe('terminal tab completion', () => {
  test('single match: completes the subcommand and adds a space', async ({ page }) => {
    await bootTerminal(page);
    await page.keyboard.type('kubectl ge');
    await page.keyboard.press('Tab');
    await page.waitForTimeout(300);

    const text = await readVisibleTerminal(page);
    // Should have extended "ge" → "get" (with trailing space). The
    // unfinished "kubectl ge" must no longer be present at end of line.
    expect(text, `terminal:\n${text}`).toContain('kubectl get');
    expect(text, `terminal:\n${text}`).not.toMatch(/kubectl ge$/m);
  });

  test('multiple matches: first Tab lists candidates and substitutes one', async ({ page }) => {
    await bootTerminal(page);
    // "de" matches both "describe" and "delete" — a single Tab should
    // list candidates AND swap the first match into the buffer (cycle
    // mode is now active).
    await page.keyboard.type('kubectl de');
    await page.keyboard.press('Tab');
    await page.waitForTimeout(200);

    const text = await readVisibleTerminal(page);
    expect(text).toContain('delete');
    expect(text).toContain('describe');
  });

  test('cycle: repeated Tab presses iterate through candidates', async ({ page }) => {
    await bootTerminal(page);
    await page.keyboard.type('kubectl de');
    await page.keyboard.press('Tab');
    await page.waitForTimeout(150);
    await page.keyboard.press('Tab');
    await page.waitForTimeout(150);

    // After cycling, the prompt line should contain the SECOND
    // candidate substituted in. Candidates are sorted alphabetically,
    // so the order is "delete" → "describe".
    const text = await readVisibleTerminal(page);
    expect(text).toMatch(/kubectl describe/);
  });

  test('Enter while cycling commits and adds a space (does not submit)', async ({ page }) => {
    await bootTerminal(page);
    await page.keyboard.type('kubectl de');
    await page.keyboard.press('Tab'); // enter cycle, "delete" substituted
    await page.waitForTimeout(200); // let the async completion settle
    await page.keyboard.press('Enter'); // commit + advance, NO submit
    await page.waitForTimeout(150);

    // Now type a name and submit a real command — proves the buffer
    // wasn't sent on the cycle-Enter and that delete now has 2 args.
    await page.keyboard.type('namespace default');
    await page.keyboard.press('Enter');
    await page.waitForTimeout(600);

    const text = await readVisibleTerminal(page);
    // The composed command ran (output contains "deleted"), and the
    // cycle-Enter did NOT trigger a premature "unknown command" path.
    expect(text).not.toMatch(/unknown command/);
    expect(text).toMatch(/deleted/);
  });

  test('dynamic completion: resource kinds after `get `', async ({ page }) => {
    await bootTerminal(page);
    await page.keyboard.type('kubectl get ');
    await page.keyboard.press('Tab');
    await page.waitForTimeout(200);

    const text = await readVisibleTerminal(page);
    // At least one canonical resource kind should appear in the listing.
    expect(text).toContain('namespace');
  });

  test('dynamic completion: resource names after `get namespace `', async ({ page }) => {
    await bootTerminal(page);
    // The preloaded data ships a "default" namespace among others; we
    // want the second positional to be filled in by Tab.
    await page.keyboard.type('kubectl get namespace ');
    await page.keyboard.press('Tab');
    await page.waitForTimeout(200);

    const text = await readVisibleTerminal(page);
    expect(text).toMatch(/kubectl get namespace \S+/);
  });
});
