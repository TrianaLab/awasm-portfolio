import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { resolve } from 'node:path';

// Vite emits the production bundle directly into the existing web/
// directory so the Go-side `make wasm` target can drop app.wasm and
// wasm_exec.js alongside it without coordination.
export default defineConfig({
  plugins: [svelte()],
  build: {
    outDir: resolve(__dirname, '../web'),
    emptyOutDir: false, // preserve assets/app.wasm placed by Go build
    target: 'esnext',
    sourcemap: true,
  },
  server: {
    port: 8000,
    fs: {
      allow: ['..'], // permit reading wasm + wasm_exec.js from sibling dirs
    },
  },
  worker: {
    // Classic IIFE — the WASM worker uses importScripts() which is only
    // available in non-module workers.
    format: 'iife',
  },
});
