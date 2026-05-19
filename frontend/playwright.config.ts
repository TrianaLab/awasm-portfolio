import { defineConfig } from '@playwright/test';

// E2E config: serves the production bundle from ../web on port 8765
// before the tests run. The production bundle (Go WASM + Svelte +
// pdfmake) must already be built — the global setup step builds it.
export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  workers: 1,
  reporter: process.env.CI ? 'list' : 'line',
  use: {
    baseURL: 'http://127.0.0.1:8765',
    trace: 'retain-on-failure',
    video: 'retain-on-failure',
    headless: true,
  },
  webServer: {
    command: 'python3 -m http.server 8765 --bind 127.0.0.1 --directory ../web',
    url: 'http://127.0.0.1:8765',
    reuseExistingServer: !process.env.CI,
    stdout: 'ignore',
    stderr: 'pipe',
  },
});
