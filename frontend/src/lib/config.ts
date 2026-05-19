// User-facing branding constants. Override per-deployment via Vite env
// variables (frontend/.env, frontend/.env.local, or shell environment at
// build time). All identifiers must start with VITE_ to be exposed to the
// client bundle.

export const BRAND = import.meta.env.VITE_BRAND ?? 'edudiaz';
export const DOMAIN = import.meta.env.VITE_DOMAIN ?? 'edudiaz.dev';
export const GITHUB_REPO = import.meta.env.VITE_GITHUB_REPO ?? 'TrianaLab/awasm-portfolio';

// The page <title> shown in the browser tab.
export const PAGE_TITLE = `${BRAND} - awasm portfolio`;
