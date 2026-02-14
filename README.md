# httprobe

A fast, concurrent HTTP endpoint health checker with colored output.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Install

```bash
go install github.com/clarabennett2626/httprobe@latest
```

Or clone and build:

```bash
git clone https://github.com/clarabennett2626/httprobe.git
cd httprobe && go build -o httprobe .
```

## Usage

Check URLs from arguments:

```bash
httprobe google.com github.com example.com
```

Pipe URLs from a file or stdin:

```bash
cat urls.txt | httprobe
echo -e "google.com\ngithub.com" | httprobe
```

### Example Output

```
URL                                                  STATUS TIME
──────────────────────────────────────────────────────────────────────
✓ https://google.com                                  200      185ms
✓ https://github.com                                  200      243ms
✗ https://doesnotexist.invalid                        ERROR   timeout
──────────────────────────────────────────────────────────────────────
2 up  1 down  3 total
```

## Features

- 🚀 **Concurrent** — all URLs checked in parallel
- 🎨 **Colored output** — green for 2xx, yellow for 3xx, red for 4xx/5xx/errors
- ⏱️ **Response times** — millisecond precision
- 📥 **Stdin support** — pipe URL lists, one per line (# comments ignored)
- 🔗 **Auto-prefix** — bare domains get `https://` automatically

## License

MIT
