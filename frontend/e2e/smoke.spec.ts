import { test, expect, type Page } from '@playwright/test';

// The terminal is now an alternate view behind #/terminal, so anything that
// exercises it starts there instead of on the landing page.
const TERMINAL = '/#/terminal';

/** The primary nav link with this exact label. */
function nav(page: Page, label: string) {
  return page.getByRole('navigation', { name: 'Primary' }).getByRole('link', { name: label, exact: true });
}

async function terminalReady(page: Page) {
  await expect(page.locator('.xterm')).toBeVisible({ timeout: 10_000 });
  await expect(page.locator('.xterm')).toContainText('Welcome', { timeout: 10_000 });
}

/**
 * Reads a window's box once it stops moving. The pane mounts from a lazy
 * chunk and the window animates in, so a box read straight after boot is a
 * mid-transition value that later geometry assertions would compare against.
 */
async function settledBox(locator: ReturnType<Page['locator']>) {
  let last = await locator.boundingBox();
  for (let i = 0; i < 20; i++) {
    await locator.page().waitForTimeout(100);
    const now = await locator.boundingBox();
    if (now && last && now.width === last.width && now.height === last.height) return now;
    last = now;
  }
  return last;
}

test.describe('portfolio', () => {
  test('home page presents the portfolio without loading the terminal', async ({ page }) => {
    const consoleErrors: string[] = [];
    const pageErrors: string[] = [];
    page.on('console', (msg) => {
      if (msg.type() !== 'error') return;
      const text = msg.text();
      // Skip GitHub API rate-limit / not-found noise — the repo card
      // retries from localStorage and degrades gracefully when unauth
      // calls hit the 60/hr quota (very common in CI).
      if (/Failed to load resource/i.test(text)) return;
      consoleErrors.push(text);
    });
    page.on('pageerror', (err) => pageErrors.push(err.message));

    await page.goto('/');
    await expect(page.getByRole('heading', { level: 1 })).toHaveText('Eduardo Díaz');
    await expect(page.getByRole('heading', { name: 'pacto' })).toBeVisible({ timeout: 10_000 });
    // The heavy chunk stays unloaded until the visitor asks for it.
    await expect(page.locator('.xterm')).toHaveCount(0);

    expect(pageErrors, `page errors:\n${pageErrors.join('\n')}`).toEqual([]);
    expect(consoleErrors, `console errors:\n${consoleErrors.join('\n')}`).toEqual([]);
  });

  // The prose was rewritten to drop the em dash; it is the tell the author
  // asked to keep out of the copy, so both the editorial strings and the
  // canonical résumé document are checked against it.
  test('rendered copy carries no em dashes', async ({ page }) => {
    for (const route of ['/', '/#/resume']) {
      await page.goto(route);
      await expect(page.getByRole('heading', { level: 1 })).toBeVisible({ timeout: 10_000 });
      const text = (await page.locator('body').innerText()) ?? '';
      expect(text, `em dash in rendered copy at ${route}`).not.toContain('—');
    }
  });

  test('hash routes are linkable, reloadable and survive back/forward', async ({ page }) => {
    await page.goto('/#/resume');
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible({ timeout: 10_000 });
    await expect(page.locator('body')).toContainText('Experience');

    await page.reload();
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible({ timeout: 10_000 });

    await nav(page, 'Work').click();
    await expect(nav(page, 'Work')).toHaveAttribute('aria-current', 'page');
    await expect(page.locator('#work')).toBeInViewport();

    await page.goBack();
    await expect(page.getByRole('heading', { level: 1 })).toHaveText(/Eduardo Díaz|Eduardo/);
  });

  test('home page scrolls normally instead of trapping scroll in a pane', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'pacto' })).toBeVisible({ timeout: 10_000 });

    const scrollable = await page.evaluate(
      () => document.documentElement.scrollHeight > window.innerHeight + 100,
    );
    expect(scrollable, 'the home page must be a normal scrolling document').toBe(true);
  });

  test('mobile: the Menu button expands the nav and navigating closes it', async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 780 });
    await page.goto('/');

    // Located by the control it owns — its label flips between Menu and Close.
    const toggle = page.locator('button[aria-controls="primary-nav"]');
    await expect(toggle).toHaveText('Menu');
    const box = await toggle.boundingBox();
    expect(box!.height, 'menu button must be a comfortable touch target').toBeGreaterThanOrEqual(44);

    await expect(nav(page, 'Résumé')).toBeHidden();
    await expect(toggle).toHaveAttribute('aria-expanded', 'false');
    await toggle.click();
    await expect(toggle).toHaveAttribute('aria-expanded', 'true');
    await expect(nav(page, 'Résumé')).toBeVisible();

    await nav(page, 'Résumé').click();
    await expect(nav(page, 'Résumé')).toBeHidden();
    await expect(page.locator('body')).toContainText('Experience', { timeout: 10_000 });
  });

  test('no horizontal overflow on a phone viewport', async ({ page }) => {
    await page.setViewportSize({ width: 360, height: 780 });
    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'pacto' })).toBeVisible({ timeout: 10_000 });
    const overflow = await page.evaluate(
      () => document.documentElement.scrollWidth - document.documentElement.clientWidth,
    );
    expect(overflow).toBeLessThanOrEqual(1);
  });

  test('education section renders course names as text, not [object Object]', async ({ page }) => {
    await page.goto('/#/resume');
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible({ timeout: 10_000 });

    // Secondary detail is collapsed on the page; open it before reading text.
    for (const d of await page.locator('details').all()) {
      await d.evaluate((el) => ((el as HTMLDetailsElement).open = true));
    }

    const body = await page.locator('body').innerText();
    expect(body, 'rendered body must include a real course name').toContain('Python Fundamentals');
    expect(body, 'should not contain [object Object] anywhere').not.toContain('[object Object]');
  });

  test('GitHub repo card renders the repo name (and meta when API resolves)', async ({ page }) => {
    await page.goto('/');
    const card = page.getByRole('link', { name: /github repository trianalab\/awasm-portfolio/i });
    await expect(card).toBeVisible({ timeout: 10_000 });
    await expect(card).toContainText('TrianaLab/awasm-portfolio');
  });

  test('download button triggers a PDF download', async ({ page }) => {
    const errs: string[] = [];
    page.on('console', (msg) => {
      if (msg.type() === 'error') errs.push(`[console] ${msg.text()}`);
    });
    page.on('pageerror', (err) => errs.push(`[pageerror] ${err.message}`));

    await page.goto('/');
    await expect(page.getByRole('heading', { name: 'pacto' })).toBeVisible({ timeout: 10_000 });

    try {
      const [download] = await Promise.all([
        page.waitForEvent('download', { timeout: 30_000 }),
        page.getByRole('button', { name: /download résumé/i }).first().click(),
      ]);
      const filename = download.suggestedFilename();
      expect(filename).toMatch(/\.pdf$/);
    } catch (e) {
      throw new Error(`download didn't fire.\nErrors:\n${errs.join('\n')}`);
    }
  });
});

test.describe('terminal', () => {
  test('boots the WASM worker and runs a kubectl command', async ({ page }) => {
    await page.goto(TERMINAL);
    await terminalReady(page);

    // xterm captures keystrokes via its hidden textarea.
    await page.locator('.xterm-helper-textarea').focus();
    await page.keyboard.type('kubectl get namespace');
    await page.keyboard.press('Enter');

    // The WASM service should answer with at least the default namespace.
    await expect(page.locator('.xterm')).toContainText('default', { timeout: 10_000 });
  });

  test('nothing runs on its own — the welcome screen only suggests commands', async ({ page }) => {
    await page.goto(TERMINAL);
    await terminalReady(page);
    await page.waitForTimeout(1500);
    await expect(page.locator('.xterm')).not.toContainText('EXPERIENCE');
  });

  test('state survives a round trip to the résumé view', async ({ page }) => {
    await page.goto(TERMINAL);
    await terminalReady(page);

    await page.getByRole('button', { name: 'New terminal' }).click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(2);

    await page.locator('.xterm-helper-textarea').first().focus();
    await page.keyboard.type('kubectl get namespace');
    await page.keyboard.press('Enter');
    await page.waitForTimeout(400);
    await expect(page.locator('.xterm').first()).toContainText('default');

    await nav(page, 'Résumé').click();
    await page.waitForTimeout(300);
    await nav(page, 'Terminal').click();
    await page.waitForTimeout(300);

    await expect(page.locator('[role="dialog"]')).toHaveCount(2);
    await expect(page.locator('.xterm').first()).toContainText('default');
  });

  test('is fully inert while another view is active', async ({ page }) => {
    await page.goto(TERMINAL);
    await terminalReady(page);

    await nav(page, 'Résumé').click();
    await expect(page.locator('.terminal-layer')).toBeHidden();
    // display:none takes it out of the tab order and the a11y tree, so the
    // résumé and the terminal are never simultaneously interactive.
    await expect(page.locator('.xterm-helper-textarea')).toBeHidden();
    await expect(page.locator('[role="dialog"]').first()).toBeHidden();
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
      await page.goto(TERMINAL);
      await terminalReady(page);

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

  test('terminal view does not grow the page height indefinitely', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto(TERMINAL);
    await terminalReady(page);
    // Let xterm + ResizeObserver settle.
    await page.waitForTimeout(500);
    const initial = await page.evaluate(() => document.documentElement.scrollHeight);
    await page.waitForTimeout(1000);
    const after = await page.evaluate(() => document.documentElement.scrollHeight);
    // Allow ±10px for the cursor blink and other normal variations.
    expect(after, `scrollHeight grew from ${initial}px to ${after}px (feedback loop?)`).toBeLessThanOrEqual(
      initial + 10,
    );
    // The terminal view is the one view that locks to the viewport.
    expect(after).toBeLessThanOrEqual(810);
  });

  test('window manager: opens a second window via the + button', async ({ page }) => {
    await page.goto(TERMINAL);
    await terminalReady(page);
    await expect(page.locator('[role="dialog"]')).toHaveCount(1);

    await page.getByRole('button', { name: 'New terminal' }).click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(2);
  });

  test('window manager: close removes the window', async ({ page }) => {
    await page.goto(TERMINAL);
    await terminalReady(page);
    await page.getByRole('button', { name: 'New terminal' }).click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(2);

    await page.locator('[role="dialog"]').last().getByRole('button', { name: /close window/i }).click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(1);
  });

  test('window manager: minimize hides the window and dock entry restores it', async ({ page }) => {
    await page.goto(TERMINAL);
    await terminalReady(page);

    await page.getByRole('button', { name: /minimize window/i }).first().click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(0);

    await page.getByRole('toolbar', { name: /minimized windows/i }).getByRole('button').first().click();
    await expect(page.locator('[role="dialog"]')).toHaveCount(1);
  });

  test('window manager: double-clicking the chrome toggles maximize', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto(TERMINAL);
    await terminalReady(page);

    const win = page.locator('[role="dialog"]').first();
    const before = await settledBox(win);
    const desktop = await page.locator('.desktop').boundingBox();
    expect(before!.width).toBeLessThan(desktop!.width);

    // Window uses manual double-click detection (350 ms window). Two
    // quick clicks on the chrome bar must trigger maximize.
    const chrome = win.locator('.chrome');
    await chrome.click({ position: { x: 300, y: 10 } });
    await chrome.click({ position: { x: 300, y: 10 } });
    // 250 ms covers the 180 ms maximize/snap CSS transition plus a small
    // jitter buffer; boundingBox() reads the in-flight value mid-anim.
    await page.waitForTimeout(250);
    let now = await win.boundingBox();
    expect(Math.abs(now!.width - desktop!.width)).toBeLessThan(3);

    await chrome.click({ position: { x: 300, y: 10 } });
    await chrome.click({ position: { x: 300, y: 10 } });
    await page.waitForTimeout(250);
    now = await win.boundingBox();
    expect(Math.abs(now!.width - before!.width)).toBeLessThan(3);
  });

  test('window manager: resize from the left edge changes width and x', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto(TERMINAL);
    await terminalReady(page);

    const win = page.locator('[role="dialog"]').first();
    const before = await settledBox(win);

    const handle = await win.locator('.rz-w').boundingBox();
    const startX = handle!.x + handle!.width / 2;
    const y = handle!.y + handle!.height / 2;
    await page.mouse.move(startX, y);
    await page.mouse.down();
    await page.mouse.move(startX + 60, y, { steps: 6 });
    await page.mouse.up();
    await page.waitForTimeout(150);

    const after = await win.boundingBox();
    expect(after!.x - before!.x).toBeGreaterThan(40);
    expect(before!.width - after!.width).toBeGreaterThan(40);
  });

  test('window manager: resize from the top edge changes height and y', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto(TERMINAL);
    await terminalReady(page);

    const win = page.locator('[role="dialog"]').first();
    const before = await settledBox(win);

    const handle = await win.locator('.rz-n').boundingBox();
    const x = handle!.x + handle!.width / 2;
    const startY = handle!.y + handle!.height / 2;
    await page.mouse.move(x, startY);
    await page.mouse.down();
    await page.mouse.move(x, startY + 50, { steps: 6 });
    await page.mouse.up();
    await page.waitForTimeout(150);

    const after = await win.boundingBox();
    expect(after!.y - before!.y).toBeGreaterThan(30);
    expect(before!.height - after!.height).toBeGreaterThan(30);
  });

  test('window manager: maximize fills the desktop and toggles back', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto(TERMINAL);
    await terminalReady(page);

    const win = page.locator('[role="dialog"]').first();
    const before = await settledBox(win);
    const desktop = await page.locator('.desktop').boundingBox();
    expect(before).not.toBeNull();
    expect(desktop).not.toBeNull();
    expect(before!.width).toBeLessThan(desktop!.width);

    await win.getByRole('button', { name: /maximize window/i }).click();
    await page.waitForTimeout(250);
    const max = await win.boundingBox();
    expect(Math.abs(max!.width - desktop!.width)).toBeLessThan(3);
    expect(Math.abs(max!.height - desktop!.height)).toBeLessThan(3);

    await win.getByRole('button', { name: /maximize window/i }).click();
    await page.waitForTimeout(250);
    const restored = await win.boundingBox();
    expect(Math.abs(restored!.width - before!.width)).toBeLessThan(3);
    expect(Math.abs(restored!.height - before!.height)).toBeLessThan(3);
  });

  test('window manager: prompt stays visible after maximizing during long output', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto(TERMINAL);
    await terminalReady(page);

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
    await page.goto(TERMINAL);
    await terminalReady(page);

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

  test('window manager: dragging a snapped window restores its previous size', async ({ page }) => {
    await page.setViewportSize({ width: 1280, height: 800 });
    await page.goto(TERMINAL);
    await terminalReady(page);

    const win = page.locator('[role="dialog"]').first();
    const before = await settledBox(win);

    // Maximize first (same code path as snap — both set previousGeometry).
    await win.getByRole('button', { name: /maximize window/i }).click();
    await page.waitForTimeout(200);
    const maxed = await win.boundingBox();
    expect(Math.abs(maxed!.width - 1280)).toBeLessThan(3);

    // Now drag the chrome to "pull off" — window should restore to the
    // original size and stay under the cursor.
    const chrome = win.locator('.chrome');
    const chromeBox = await chrome.boundingBox();
    const startX = chromeBox!.x + 300;
    const startY = chromeBox!.y + 12;
    await page.mouse.move(startX, startY);
    await page.mouse.down();
    await page.mouse.move(startX + 60, startY + 60, { steps: 8 });
    await page.mouse.up();
    await page.waitForTimeout(200);

    const after = await win.boundingBox();
    // Width must shrink back from maximized → roughly the pre-maximize size.
    expect(Math.abs(after!.width - before!.width)).toBeLessThan(20);
    expect(Math.abs(after!.height - before!.height)).toBeLessThan(20);
  });

  test('phone viewport: terminal fills the desktop', async ({ page }) => {
    // iPhone-ish viewport. On phones the window forcibly fills the
    // desktop area so the terminal is actually usable.
    await page.setViewportSize({ width: 390, height: 780 });
    await page.goto(TERMINAL);
    await terminalReady(page);

    const win = page.locator('[role="dialog"]').first();
    const desktop = await page.locator('.desktop').boundingBox();
    const winBox = await win.boundingBox();
    expect(desktop).not.toBeNull();
    expect(winBox).not.toBeNull();
    // Window must span the whole desktop area on phones.
    expect(Math.abs(winBox!.width - desktop!.width)).toBeLessThan(3);
    expect(Math.abs(winBox!.height - desktop!.height)).toBeLessThan(3);
    // Resize handles must be hidden so they don't trap taps.
    await expect(win.locator('.rz-se')).toBeHidden();
  });

  test('window manager: focus brings background window to front', async ({ page }) => {
    await page.goto(TERMINAL);
    await terminalReady(page);
    await page.getByRole('button', { name: 'New terminal' }).click();
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
    await page.goto(TERMINAL);
    await terminalReady(page);

    // Snapshot the preloaded resume content first.
    await nav(page, 'Résumé').click();
    await expect(page.locator('body')).toContainText('Eduardo', { timeout: 10_000 });

    // Mutate via the terminal: delete + recreate the resume.
    await nav(page, 'Terminal').click();
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
    await nav(page, 'Résumé').click();
    await page.waitForTimeout(500);
    const body = await page.locator('main').innerText();
    expect(body, 'after recreate the resume should be empty — no Experience section').not.toContain(
      'Experience',
    );
  });
});
