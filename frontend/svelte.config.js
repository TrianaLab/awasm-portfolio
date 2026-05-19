import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

export default {
  preprocess: vitePreprocess(),
  compilerOptions: {
    // Svelte 5 runes ($state, $derived, $effect) are the canonical
    // reactivity primitives. Disable legacy reactive statements.
    runes: true,
  },
};
