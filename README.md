# httprobe

Fast HTTP probe — check which URLs are alive from stdin. Reads URLs line by line and outputs the ones that respond.

## Install

```bash
go install github.com/clarabennett2626/httprobe@latest
```

## Usage

```bash
# Basic usage — pipe URLs in
cat urls.txt | httprobe

# Show status codes
echo -e "google.com\ngithub.com" | httprobe -s
# [301] https://google.com
# [200] https://github.com

# Adjust concurrency and timeout
cat large-list.txt | httprobe -c 20 -t 10
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c` | 10 | Concurrency level |
| `-t` | 5 | Timeout in seconds |
| `-s` | false | Show HTTP status codes |

## How it works

- Reads URLs from stdin (one per line)
- Automatically prepends `https://` if no scheme provided
- Probes each URL concurrently
- Outputs only responsive URLs

Useful for bug bounty recon, monitoring, or quickly checking which hosts from a list are up.

## License

MIT
