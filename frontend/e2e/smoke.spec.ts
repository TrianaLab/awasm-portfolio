import { test, expect } from '@playwright/test';

test.describe('awasm-portfolio smoke', () => {
  test('loads without console errors and boots the WASM worker', async ({ page }) => {
    const consoleErrors: string[] = [];
    const pageErrors: string[] = [];
    page.on('console', (msg) => {
      if (msg.type() === 'error') consoleErrors.push(msg.text());
    });
    page.on('pageerror', (err) => pageErrors.push(err.message));

    await page.goto('/');
    // Wait for the xterm welcome banner to render — it only shows after
    // the Worker + WASM bootstrap chain has set things up.
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });
    await expect(page.locator('.xterm')).toContainText('Welcome', { timeout: 10_000 });

    expect(pageErrors, `page errors:\n${pageErrors.join('\n')}`).toEqual([]);
    expect(consoleErrors, `console errors:\n${consoleErrors.join('\n')}`).toEqual([]);
  });

  test('runs a kubectl command in the terminal', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.xterm')).toContainText('Welcome', { timeout: 10_000 });

    // xterm captures keystrokes via its hidden textarea.
    const helper = page.locator('.xterm-helper-textarea');
    await helper.focus();
    await page.keyboard.type('kubectl get namespace');
    await page.keyboard.press('Enter');

    // The WASM service should answer with at least the default namespace.
    await expect(page.locator('.xterm')).toContainText('default', { timeout: 10_000 });
  });

  test('switches to resume view and renders sections', async ({ page }) => {
    const errs: string[] = [];
    page.on('console', (msg) => {
      if (msg.type() === 'error') errs.push(`[console] ${msg.text()}`);
    });
    page.on('pageerror', (err) => errs.push(`[pageerror] ${err.message}`));

    await page.goto('/');
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });

    await page.getByRole('tab', { name: /resume/i }).click();
    try {
      await expect(page.getByRole('heading', { level: 1 })).toBeVisible({ timeout: 10_000 });
    } catch (e) {
      const bodyText = await page.locator('body').innerText();
      throw new Error(`heading not found.\nBody text:\n${bodyText}\nErrors:\n${errs.join('\n')}`);
    }
    await expect(page.locator('body')).toContainText('Experience');
  });

  test('download button triggers a PDF download', async ({ page }) => {
    const errs: string[] = [];
    page.on('console', (msg) => {
      if (msg.type() === 'error') errs.push(`[console] ${msg.text()}`);
    });
    page.on('pageerror', (err) => errs.push(`[pageerror] ${err.message}`));

    await page.goto('/');
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });

    try {
      const [download] = await Promise.all([
        page.waitForEvent('download', { timeout: 30_000 }),
        page.getByRole('button', { name: /download resume as pdf/i }).click(),
      ]);
      const filename = download.suggestedFilename();
      expect(filename).toMatch(/\.pdf$/);
    } catch (e) {
      throw new Error(`download didn't fire.\nErrors:\n${errs.join('\n')}`);
    }
  });
});
