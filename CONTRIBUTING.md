# Contributing to DeadEndProxy

🧠 **Thanks for your interest in contributing to DeadEndProxy!**  
We're building a fast, lightweight reverse proxy written in Go with system-level integrations for Linux and BSD.

## Getting Started

1. Fork the repository.
2. Clone your fork:

```sh
git clone https://github.com/yourusername/DeadEndProxy.git
cd DeadEndProxy
```

3. Make sure you have Go 1.22+ installed.
4. Build the binary:

```sh
go build -o deadendproxy-bin ./cmd
```

## Running Tests

Run all tests using:

```sh
go test ./...
```

## Code Style

- Follow Go conventions (`go fmt` is your friend).
- Use clear commit messages. Recommended format:

```
feat(config): add support for dynamic port binding
fix(proxy): resolve nil pointer on empty route
```

## Pull Request Process

1. Create a new branch from `main`.
2. Commit your changes clearly.
3. Push to your fork and open a PR.
4. Wait for code review.

## Feature Proposals

We welcome feature requests! Open an issue with the `[Feature]` tag and explain your idea. Please include:

- Why the feature is useful
- Example usage
- Any breaking changes

## Code of Conduct

All contributors are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md).

---

✉️ Questions or ideas?  
Email us at: **admin@devinsidercode.com**