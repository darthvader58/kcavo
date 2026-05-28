# Contributing to KCAVO

Thank you for your interest in contributing to KCAVO. This project aims to provide a lightweight, practical `kubectl` plugin for Kubernetes cost visibility, GPU allocation reporting, and basic optimization guidance.

## Project Scope

KCAVO estimates cost from Kubernetes resource specifications. It is not intended to be a billing system, cloud invoice reconciler, or replacement for metrics-backed cost platforms.

Good contributions usually improve one of these areas:

- More accurate request-based cost estimation
- Clearer Kubernetes resource analysis
- Better GPU allocation reporting
- Safer and more useful optimization recommendations
- CLI usability, structured output, tests, and documentation

## Project Layout

```text
.
├── cmd/                    # Cobra commands and CLI wiring
│   ├── analyze.go          # Cost analysis command
│   ├── gpu.go              # GPU analysis command
│   ├── optimize.go         # Optimization recommendation command
│   ├── root.go             # Root command, config, and shared flags
│   └── visualize.go        # Resource visualization command
├── configs/
│   └── pricing.yaml        # Example pricing configuration
├── pkg/
│   ├── cost/               # Pricing and cost calculation logic
│   ├── gpu/                # GPU allocation analysis
│   ├── kubernetes/         # Kubernetes client wrapper
│   ├── optimize/           # Recommendation logic
│   └── visualize/          # Table, JSON, and YAML output helpers
├── main.go                 # CLI entry point
├── Makefile                # Common development commands
├── install.sh              # Local installation helper
├── Readme.md               # User documentation
└── go.mod                  # Go module definition
```

## Local Development

Requirements:

- Go 1.21 or newer
- `kubectl` for testing against a real cluster
- Access to a Kubernetes cluster for integration testing

Common commands:

```bash
make fmt
make test
make build
```

Run the CLI locally:

```bash
go run . analyze
go run . visualize
go run . gpu
go run . optimize
```

## Contribution Workflow

1. Fork the repository.
2. Create a focused branch from `main`.
3. Make a small, reviewable change.
4. Add or update tests for behavior changes.
5. Run `make fmt` and `make test`.
6. Open a pull request with a clear description.

Keep pull requests scoped. Separate unrelated changes into separate PRs.

## Pull Request Format

Use this structure in your pull request description:

```markdown
## Summary

- What changed?
- Why is this change needed?

## Testing

- [ ] `make fmt`
- [ ] `make test`
- [ ] Tested against a Kubernetes cluster, if applicable

## Notes

- Mention limitations, follow-up work, or compatibility concerns.
```

## Code Guidelines

- Keep CLI output clear, predictable, and script-friendly.
- Do not print progress text for JSON or YAML output.
- Prefer structured Kubernetes APIs over string parsing.
- Keep cost logic deterministic and covered by tests.
- Treat estimates as estimates. Avoid presenting inferred savings as exact billing data.
- Keep dependencies minimal and relevant to the CLI.
- Add comments only when they explain non-obvious behavior.

## Testing Guidelines

Add unit tests for:

- Cost calculation changes
- GPU allocation logic
- Optimization recommendation behavior
- Sorting, filtering, or output format behavior

For changes that require a real cluster, document what was tested in the pull request.

## Documentation Guidelines

Update `Readme.md` when a change affects:

- Installation
- Commands or flags
- Configuration
- Output formats
- Cost model behavior

Documentation should be concise, accurate, and focused on user workflows.

## Commit Message Style

Use short, conventional commit messages:

```text
feat: add namespace cost summary
fix: avoid double-counting GPU limits
docs: update configuration examples
test: cover optimizer recommendation sorting
chore: clean install script output
```

## Review Expectations

Maintainers will look for:

- Correct behavior
- Clear user impact
- Tests for meaningful logic changes
- Minimal scope
- Professional CLI and documentation quality

Security-sensitive changes, dependency additions, and Kubernetes permission changes may require extra review.
