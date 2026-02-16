# Contributing to httprobe

Thanks for your interest in contributing! Here's how to get started.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/httprobe.git`
3. Create a branch: `git checkout -b my-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Commit and push
7. Open a pull request

## Development

```bash
go build -o httprobe .
go test -v ./...
```

## Reporting Issues

Please use GitHub Issues for bug reports and feature requests. Include:
- Go version (`go version`)
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Add tests for new features
- Keep commits focused and well-described
