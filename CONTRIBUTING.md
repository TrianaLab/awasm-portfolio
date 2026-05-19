# Contributing to awasm-portfolio

Thanks for your interest in contributing! This guide covers prereqs, the dev workflow, and CI gates.

## Code of Conduct

By participating in this project, you agree to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).

## Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Node 22+](https://nodejs.org/)
- [Git](https://git-scm.com/)
- `make`
- [golangci-lint](https://golangci-lint.run/welcome/install/)

## Setting up your environment

```bash
git clone https://github.com/<your-username>/awasm-portfolio.git
cd awasm-portfolio
go mod download
cd frontend && npm install
```

## Dev workflow

```bash
make wasm        # build the Go WASM artifact once
cd frontend && npm run dev    # then HMR for the Svelte side
```

Or for a full production build:

```bash
make build && make run        # http://127.0.0.1:8000
```

Before pushing:

```bash
make ci          # runs every quality gate
```

`make ci` is the same set of checks the GitHub workflows run. If it passes locally it'll pass in CI.

## Quality gates (`make ci`)

| Gate          | What it checks                                           |
| ------------- | -------------------------------------------------------- |
| `ci-fmt`      | gofmt -s on every non-vendored Go file                   |
| `ci-vet`      | go vet ./...                                             |
| `ci-cyclo`    | gocyclo ≤ 15 across all production + test Go code        |
| `ci-lint`     | golangci-lint (with misspell)                            |
| `ci-test`     | 100% statement coverage on every Go package              |
| `ci-ui`       | svelte-check, vitest, eslint, vite build                 |

## Pull-request convention

PR titles must follow the conventional-commit format — the auto-release workflow uses them to compute the next semver tag.

```
feat: add tab completion to the terminal              # → minor bump
fix: handle empty highlights in PDF section           # → patch bump
feat!: switch resume schema to JSON Resume v1.1       # → major bump
docs: clarify the dev workflow                        # → no release
```

Allowed types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `ci`, `perf`, `build`, `style`. Append `!` before `:` for breaking changes.

## Project layout

```
cmd/                      Cobra command handlers (kubectl, get/create/delete/describe/version)
internal/
  models/                 Resource interface + shared Meta struct
    types/                Concrete resource types (Award, Basics, Education, ...)
  preload/                Static portfolio content
  repository/             In-memory key-value store
  service/                Flat business-logic helpers (Create/Delete/Get/Describe)
  ui/                     Output formatters: table, details (YAML), JSON
  util/                   Kind normalisation + aliases
cli.go                    Native CLI entrypoint  (excluded from WASM build)
main.go                   WASM entrypoint        (//go:build js && wasm)
frontend/                 Svelte 5 + Vite SPA
  src/lib/                wasm bridge, theme store, pdf builder, schema types
  src/components/         App, Terminal, ResumeView, Section, TimelineEntry, ...
web/                      Build output (Vite + Go WASM artefacts) — gitignored
```

## Releasing

Releases are automated. When a PR is merged into `main`:

1. `auto-release.yml` inspects the PR title.
2. It computes the next semver tag (major / minor / patch / skip).
3. A GitHub Release is created with auto-generated notes.
4. `cd.yml` and `release.yml` then build the WASM + Svelte bundle, push the container image to GHCR, and deploy to GitHub Pages.

## License

By contributing you agree your changes will be licensed under the [MIT License](LICENSE).
