# CI-specific targets. Included by the main Makefile.
# This file is the single source of truth for all CI quality gates.
# Do not edit without verifying that the GitHub Actions pipelines still match.

.PHONY: ci ci-static ci-test ci-fmt ci-vet ci-cyclo ci-lint

ci: ci-static ci-test

ci-static: ci-fmt ci-vet ci-cyclo ci-lint

ci-fmt:
	@echo "==> Checking formatting..."
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './web/*'))" || \
		(echo "gofmt found unformatted files:" && gofmt -l $$(find . -name '*.go' -not -path './web/*') && exit 1)

ci-vet:
	@echo "==> Running go vet..."
	@go vet ./...

ci-cyclo:
	@echo "==> Checking cyclomatic complexity..."
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@gocyclo -over 15 $$(find . -name '*.go' -not -path './web/*')

ci-lint:
	@echo "==> Running golangci-lint..."
	@golangci-lint run

ci-test:
	@echo "==> Running unit tests with coverage..."
	@# The root package (cli.go, main.go) is excluded — it is the entrypoint
	@# (WASM and native CLI bootstrap) and is also excluded by codecov.yml.
	@go test $$(go list ./... | grep -v /tests/ | grep -v '^github.com/TrianaLab/awasm-portfolio$$') -coverprofile=coverage.out -covermode=atomic
	@total=$$(go tool cover -func=coverage.out | grep '^total:' | awk '{print $$NF}'); \
	if [ "$$total" != "100.0%" ]; then \
		echo "FAIL: total coverage is $$total, expected 100.0%"; \
		go tool cover -func=coverage.out | grep -v '100.0%'; \
		exit 1; \
	fi
	@echo "    total coverage: 100.0%"
