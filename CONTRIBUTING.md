# Contributing to x402Fuel

Thanks for your interest in x402Fuel! This document outlines how to contribute.

## Getting Started

1. Fork the repo and clone it
2. `cd x402fuel && go mod download`
3. Make your changes on a feature branch
4. Run `go test ./...` to ensure all tests pass
5. Run `go vet ./...` to check for issues
6. Submit a PR against `main`

## TDD Requirement

**All code changes must follow test-driven development:**

1. Write a failing test that defines the expected behavior
2. Write the minimal code to make the test pass
3. Refactor while keeping tests green

PRs without tests will not be merged.

## Logging

Use the standard logging format for all new code:

```go
log.Printf("[INFO] [component] message key=value")
```

See the [Engineering Standards](https://github.com/trucore-ai/x402fuel/blob/main/CONTRIBUTING.md#engineering-standards) section for details.

## Engineering Standards

- **Single Responsibility:** Each type and function does one thing
- **Dependency Injection:** Dependencies passed via constructors, no global state
- **Error Wrapping:** Use `fmt.Errorf("context: %w", err)` — never swallow errors
- **Fail Fast:** Validate inputs at function boundaries
- **No key material in logs:** Enforced by tests, not convention

## Commit Convention

We use conventional commits:
- `feat:` new feature
- `fix:` bug fix
- `docs:` documentation
- `test:` adding tests
- `refactor:` code change that neither fixes a bug nor adds a feature
- `chore:` build process or tooling changes
- `ci:` CI configuration changes

## Questions?

Open an issue or reach out on [x402 Slack](http://slack.x402.org/).