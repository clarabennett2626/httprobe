# Changelog

All notable changes to httprobe will be documented in this file.

## [Unreleased]

### Added
- `-f` flag for fail-only mode (only show non-responding hosts)
- `-q` flag for quiet mode (exit code only, no output)
- `-c` flag for concurrency control
- `-t` flag for custom timeout
- `-j` flag for JSON output format
- `CONTRIBUTING.md` with contribution guidelines

### Fixed
- Improved error handling for malformed URLs

## [1.0.0] - 2026-02-14

### Added
- Initial release
- HTTP/HTTPS probe support
- Configurable ports via `-p` flag
- Stdin/file input support
