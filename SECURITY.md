# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |

Only the latest release is actively supported. Please run the most recent version.

## Reporting a Vulnerability

If you discover a security vulnerability in awasm-portfolio, please report it responsibly. **Do not open a public GitHub issue.**

### How to Report

Submit a private report via [GitHub Security Advisories](https://github.com/TrianaLab/awasm-portfolio/security/advisories/new) and include:

- A description of the vulnerability
- Steps to reproduce
- Potential impact
- Any suggested fix (if applicable)

### What to Expect

- **Acknowledgment** within 48 hours.
- **Updates** as we investigate and work on a fix.
- **Disclosure** coordinated with you once a fix is released. We aim to resolve critical issues within 30 days.

## Scope

In scope:

- The awasm-portfolio WebAssembly binary and its Go source
- The static frontend in `web/`
- The container image published to `ghcr.io/trianalab/awasm-portfolio`

Out of scope:

- Vulnerabilities in upstream dependencies — please report those to the upstream project and notify us so we can update.
- Third-party CDN-hosted assets referenced by the frontend.
