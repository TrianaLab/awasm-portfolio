/* eslint-env node */
module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  extends: ['eslint:recommended', 'plugin:@typescript-eslint/recommended', 'plugin:svelte/recommended', 'prettier'],
  plugins: ['@typescript-eslint'],
  parserOptions: { ecmaVersion: 2022, sourceType: 'module', extraFileExtensions: ['.svelte'] },
  env: { browser: true, es2022: true, node: true },
  overrides: [
    {
      files: ['*.svelte'],
      parser: 'svelte-eslint-parser',
      parserOptions: { parser: '@typescript-eslint/parser' },
    },
  ],
  ignorePatterns: ['node_modules/', 'dist/', '../web/'],
};
