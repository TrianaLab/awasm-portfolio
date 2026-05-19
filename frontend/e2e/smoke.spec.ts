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

  // Run at a few viewport sizes so we catch regressions on smaller windows.
  for (const viewport of [
    { width: 1280, height: 720 },
    { width: 1024, height: 600 },
    { width: 800, height: 500 },
  ]) {
    test(`prompt stays visible after a long-output command @ ${viewport.width}x${viewport.height}`, async ({
      page,
    }) => {
      await page.setViewportSize(viewport);
      await page.goto('/');
      await expect(page.locator('.xterm')).toContainText('Welcome', { timeout: 10_000 });

      const helper = page.locator('.xterm-helper-textarea');
      await helper.focus();
      await page.keyboard.type('kubectl get all -A');
      await page.keyboard.press('Enter');
      await page.waitForTimeout(800);

      const result = await page.evaluate(() => {
        const v = document.querySelector('.xterm-viewport') as HTMLElement | null;
        if (!v) return { ok: false, reason: 'no viewport' };
        const slack = v.scrollHeight - v.scrollTop - v.clientHeight;
        return { ok: slack < 24, slack, scrollHeight: v.scrollHeight, clientHeight: v.clientHeight };
      });
      expect(
        result.ok,
        `terminal must auto-scroll to the prompt after output (slack=${JSON.stringify(result)})`,
      ).toBe(true);
    });
  }

  test('terminal mode does not grow the page height indefinitely', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto('/');
    await expect(page.locator('.xterm')).toContainText('Welcome', { timeout: 10_000 });
    // Let xterm + ResizeObserver settle.
    await page.waitForTimeout(500);
    const initial = await page.evaluate(() => document.documentElement.scrollHeight);
    await page.waitForTimeout(1000);
    const after = await page.evaluate(() => document.documentElement.scrollHeight);
    // Allow ±10px for the cursor blink and other normal variations.
    expect(after, `scrollHeight grew from ${initial}px to ${after}px (feedback loop?)`).toBeLessThanOrEqual(
      initial + 10,
    );
    // The page should fit in the viewport — no body scrolling.
    expect(after).toBeLessThanOrEqual(810);
  });

  test('education section renders course names as text, not [object Object]', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });

    await page.getByRole('tab', { name: /resume/i }).click();
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible({ timeout: 10_000 });

    const body = await page.locator('body').innerText();
    expect(body, 'rendered body must include a real course name').toContain('Python Fundamentals');
    expect(body, 'should not contain [object Object] anywhere').not.toContain('[object Object]');
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

  test('window manager: opens a second window via the + button', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });
    await expect(page.locator('[role="dialog"]')).toHaveCount(1);

    await page.getByRole('button', { name: /open a new terminal/i }).click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(2);
  });

  test('window manager: close removes the window', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });
    await page.getByRole('button', { name: /open a new terminal/i }).click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(2);

    await page.locator('[role="dialog"]').last().getByRole('button', { name: /close window/i }).click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(1);
  });

  test('window manager: minimize hides the window and dock entry restores it', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });

    await page.getByRole('button', { name: /minimize window/i }).first().click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(0);

    await page.getByRole('toolbar', { name: /minimized windows/i }).getByRole('button').first().click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(1);
  });

  test('window manager: maximize fills the desktop and toggles back', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto('/');
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });

    const win = page.locator('[role="dialog"]').first();
    const before = await win.boundingBox();
    const desktop = await page.locator('.desktop').boundingBox();
    expect(before).not.toBeNull();
    expect(desktop).not.toBeNull();
    expect(before!.width).toBeLessThan(desktop!.width);

    await win.getByRole('button', { name: /maximize window/i }).click();
    await page.waitForTimeout(150);
    const max = await win.boundingBox();
    expect(Math.abs(max!.width - desktop!.width)).toBeLessThan(3);
    expect(Math.abs(max!.height - desktop!.height)).toBeLessThan(3);

    await win.getByRole('button', { name: /maximize window/i }).click();
    await page.waitForTimeout(150);
    const restored = await win.boundingBox();
    expect(Math.abs(restored!.width - before!.width)).toBeLessThan(3);
    expect(Math.abs(restored!.height - before!.height)).toBeLessThan(3);
  });

  test('window manager: prompt stays visible after maximizing during long output', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto('/');
    await expect(page.locator('.xterm')).toContainText('Welcome', { timeout: 10_000 });

    const helper = page.locator('.xterm-helper-textarea');
    await helper.focus();
    await page.keyboard.type('kubectl get all -A');
    await page.keyboard.press('Enter');
    await page.waitForTimeout(600);

    await page
      .locator('[role="dialog"]')
      .first()
      .getByRole('button', { name: /maximize window/i })
      .click();
    await page.waitForTimeout(400);

    const ok = await page.evaluate(() => {
      const v = document.querySelector('.xterm-viewport') as HTMLElement | null;
      if (!v) return false;
      return v.scrollHeight - v.scrollTop - v.clientHeight < 24;
    });
    expect(ok, 'after maximize, viewport must be scrolled to the prompt').toBe(true);
  });

  test('window manager: prompt is not clipped after long-output command', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto('/');
    await expect(page.locator('.xterm')).toContainText('Welcome', { timeout: 10_000 });

    const helper = page.locator('.xterm-helper-textarea');
    await helper.focus();
    await page.keyboard.type('kubectl get all -A');
    await page.keyboard.press('Enter');
    await page.waitForTimeout(800);

    // The bottom of the xterm canvas must sit above the window border so
    // the prompt's last row is not partially clipped.
    const clearance = await page.evaluate(() => {
      const dialog = document.querySelector('[role="dialog"]') as HTMLElement | null;
      const xterm = dialog?.querySelector('.xterm') as HTMLElement | null;
      if (!dialog || !xterm) return -1;
      return dialog.getBoundingClientRect().bottom - xterm.getBoundingClientRect().bottom;
    });
    expect(clearance, 'xterm canvas must finish inside the window').toBeGreaterThan(0);
  });

  test('window manager: focus brings background window to front', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });
    await page.getByRole('button', { name: /open a new terminal/i }).click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(2);

    const zBefore = await page.evaluate(() => {
      const ds = Array.from(document.querySelectorAll('[role="dialog"]')) as HTMLElement[];
      return ds.map((d) => parseInt(d.style.zIndex, 10));
    });
    expect(zBefore[1]).toBeGreaterThan(zBefore[0]);

    // Click the first window's chrome — should raise it above the second.
    await page.locator('[role="dialog"]').first().click({ position: { x: 100, y: 5 } });
    const zAfter = await page.evaluate(() => {
      const ds = Array.from(document.querySelectorAll('[role="dialog"]')) as HTMLElement[];
      return ds.map((d) => parseInt(d.style.zIndex, 10));
    });
    expect(zAfter[0]).toBeGreaterThan(zAfter[1]);
  });

  test('resume view re-fetches after a delete + create in the terminal', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('.xterm')).toContainText('Welcome', { timeout: 10_000 });

    // Snapshot the preloaded resume content first.
    await page.getByRole('tab', { name: /resume/i }).click();
    await expect(page.locator('body')).toContainText('Eduardo', { timeout: 10_000 });

    // Mutate via the terminal: delete + recreate the resume.
    await page.getByRole('tab', { name: /terminal/i }).click();
    const helper = page.locator('.xterm-helper-textarea');
    await helper.focus();
    await page.keyboard.type('kubectl delete resume main-resume');
    await page.keyboard.press('Enter');
    await page.waitForTimeout(400);
    await page.keyboard.type('kubectl create resume main-resume');
    await page.keyboard.press('Enter');
    await page.waitForTimeout(400);

    // Switch back to the resume view; refreshResume must refetch and the
    // recreated (empty) resume should render — no Experience section.
    await page.getByRole('tab', { name: /resume/i }).click();
    await page.waitForTimeout(500);
    const body = await page.locator('body').innerText();
    expect(body, 'after recreate the resume should be empty — no Experience section').not.toContain(
      'Experience',
    );
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
