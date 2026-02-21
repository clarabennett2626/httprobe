# httprobe 🔍

A fast, concurrent HTTP endpoint health checker with colored output.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/clarabennettdev/httprobe)](https://goreportcard.com/report/github.com/clarabennettdev/httprobe)

Pipe a list of URLs and get instant feedback on which ones are alive, with response times, status codes, and color-coded output.

## Install

```bash
go install github.com/clarabennettdev/httprobe@latest
```

## Usage

```bash
# Check URLs from a file
cat urls.txt | httprobe

# With concurrency and timeout
cat urls.txt | httprobe -c 20 -t 5000

# JSON output for scripting
cat urls.txt | httprobe -j
```

## Options

| Flag | Description | Default |
|------|-------------|---------|
| `-c` | Concurrency (parallel requests) | `10` |
| `-t` | Timeout per request (ms) | `3000` |
| `-j` | JSON output | `false` |

## Example Output

```
✓ 200  https://example.com         123ms
✓ 301  https://google.com          89ms
✗ ERR  https://doesnotexist.xyz    timeout
```

## License

MIT
