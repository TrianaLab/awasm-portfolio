# Contributing to awasm-portfolio

Thank you for your interest in contributing to awasm-portfolio! This guide will help you get started.

## Code of Conduct

By participating in this project, you agree to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

### Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Git](https://git-scm.com/)
- A terminal with `make` available
- [golangci-lint](https://golangci-lint.run/welcome/install/) (for linting)
- Python 3 (for the local dev server)

### Setting Up Your Development Environment

1. **Fork and clone the repository:**

   ```bash
   git clone https://github.com/<your-username>/awasm-portfolio.git
   cd awasm-portfolio
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Build the WebAssembly binary:**

   ```bash
   make build
   ```

4. **Run the local dev server:**

   ```bash
   make run
   ```

   The portfolio is served at `http://127.0.0.1:8000`.

5. **Run the full CI pipeline locally:**

   ```bash
   make ci
   ```

   This runs everything the CI pipeline checks — formatting, vetting, cyclomatic complexity, linting, and 100% unit test coverage. **Always run `make ci` before pushing** to catch issues early.

   You can also run individual targets:

   ```bash
   make test         # unit tests
   make lint         # gofmt + go vet
   make coverage     # coverage report with HTML output
   ```

## How to Contribute

### Reporting Bugs

If you find a bug, please [open an issue](https://github.com/TrianaLab/awasm-portfolio/issues/new?template=bug_report.yml) using the bug report template. Include:

- Steps to reproduce
- Expected vs. actual behavior
- Your environment (OS, browser, Go version)
- Relevant logs or error messages

### Suggesting Features

Have an idea? [Open a feature request](https://github.com/TrianaLab/awasm-portfolio/issues/new?template=feature_request.yml). Describe the problem you're trying to solve and the solution you'd like to see.

### Submitting Changes

1. **Create a branch** from `main`:

   ```bash
   git checkout -b feat/my-feature
   ```

   Use a descriptive branch name with a prefix: `feat/`, `fix/`, `docs/`, `refactor/`, `test/`.

2. **Make your changes.** Keep commits focused and atomic.

3. **Write or update tests.** All new functionality must include tests. All bug fixes must include a regression test. The project enforces **100% statement coverage**.

4. **Run the CI pipeline locally before pushing:**

   ```bash
   make ci
   ```

5. **Use a conventional commit-style PR title.** The PR title drives semantic version bumps via the auto-release workflow:

   ```
   feat: add reverse-DNS resource shortcut       # minor bump
   fix: render unicode glyphs in skill panel     # patch bump
   feat!: drop legacy kubectl alias              # major bump
   docs: improve customization guide             # no release
   ```

   Allowed types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `ci`, `perf`, `build`, `style`. Add `!` before `:` for breaking changes.

6. **Open a pull request** against `main`. Fill in the PR template and link any related issues.

## Development Guidelines

### Project Structure

```
awasm-portfolio/
  cli.go              # Native CLI entrypoint (excluded from coverage)
  main.go             # WebAssembly entrypoint (js+wasm build tag)
  cmd/                # Cobra command handlers (kubectl-style)
  internal/
    config/           # Internal config
    factory/          # Mock-data factory (gofakeit)
    logger/           # Logging
    middleware/       # Command middleware
    models/           # Domain model types (resume sections)
    preload/          # Static portfolio data (edit here to customize)
    repository/       # In-memory key-value store
    service/          # Business logic
    ui/               # Output formatters (YAML, JSON, table, details)
    util/             # Shared utilities
  web/                # Static frontend (HTML/JS/CSS + xterm.js)
```

To customize the portfolio content, edit `internal/preload/preload.go`.

### Code Style

- Follow standard Go conventions and idioms.
- Code must pass `golangci-lint` (run via `make ci`).
- Keep functions small and focused. Cyclomatic complexity must stay at 15 or below.

### Testing

- **Unit tests** live alongside the code they test (`_test.go` files).
- The project enforces **100% statement coverage**. `make ci` will fail if total coverage drops below 100%.
- Run `make coverage` to generate a coverage report and identify uncovered lines.

### CI Quality Gates

The `make ci` target runs all quality gates in order:

| Gate | What it checks |
|------|---------------|
| `ci-fmt` | All files are `gofmt`-formatted |
| `ci-vet` | `go vet` passes on all packages |
| `ci-cyclo` | No function exceeds cyclomatic complexity 15 |
| `ci-lint` | `golangci-lint` reports zero issues |
| `ci-test` | Unit tests pass with 100% coverage |

## Releasing

Releases are automated. When a PR is merged into `main`:

1. The `auto-release` workflow inspects the PR title.
2. If it's a `feat:`, `fix:`, `refactor:`, `perf:`, or breaking change (`!:`), a new semver tag is computed (minor, patch, or major).
3. A GitHub Release is created with auto-generated notes.
4. The `cd` and `release` workflows then deploy to GitHub Pages and publish a container image to GHCR.

## Questions?

If you're unsure about anything, feel free to [open an issue](https://github.com/TrianaLab/awasm-portfolio/issues) or ask in your pull request.

## License

By contributing to awasm-portfolio, you agree that your contributions will be licensed under the [MIT License](LICENSE).
