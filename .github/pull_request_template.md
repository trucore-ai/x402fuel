## Description

<!-- What does this PR do? -->

## Type of change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Checklist

- [ ] Tests added and passing (`go test ./...`)
- [ ] Code follows TDD: tests written before implementation
- [ ] Logging added at key decision points
- [ ] No key material in logs (verified)
- [ ] `go vet ./...` passes
- [ ] Errors wrapped with context (`fmt.Errorf("context: %w", err)`)
- [ ] Dependencies injected via constructors (no global state)
- [ ] README updated if needed

## Related issues

<!-- Link to issues this PR addresses -->

Closes #