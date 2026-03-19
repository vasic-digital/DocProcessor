# Contributing

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch
4. Make your changes
5. Run tests: `go test ./... -race -count=1`
6. Submit a pull request

## Development Setup

```bash
git clone git@github.com:YOUR_USERNAME/DocProcessor.git
cd DocProcessor
go mod tidy
go test ./... -race -count=1
```

## Code Standards

- All code must pass `go vet`
- All code must be formatted with `gofmt`
- All public functions must have doc comments
- All source files must have SPDX license headers
- Tests must pass with `-race -count=1`

## Testing Requirements

Every package must have:
- Unit tests (`*_test.go`)
- Integration tests (`*_integration_test.go`) where applicable
- Stress tests (`*_stress_test.go`) for concurrent code
- Security tests (`*_security_test.go`) for input validation

**NO test may ever be removed, disabled, skipped, or left broken.**

## Commit Messages

Follow conventional commit format:
```
type: short description

Longer description if needed.

Co-Authored-By: Your Name <your@email.com>
```

Types: `feat`, `fix`, `test`, `docs`, `refactor`, `chore`

## Pull Request Process

1. Ensure all tests pass
2. Update documentation if needed
3. Add tests for new functionality
4. Request review from maintainers

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
