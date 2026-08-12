// User-facing branding constants. Override per-deployment via Vite env
// variables (frontend/.env, frontend/.env.local, or shell environment at
// build time). All identifiers must start with VITE_ to be exposed to the
// client bundle.

export const DOMAIN = import.meta.env.VITE_DOMAIN ?? 'edudiaz.dev';
export const GITHUB_REPO = import.meta.env.VITE_GITHUB_REPO ?? 'TrianaLab/awasm-portfolio';

// Fallback <title>, shown until the canonical résumé loads and
// resume-select.pageTitle() rebuilds it from basics.name + the current role.
export const PAGE_TITLE = DOMAIN;
