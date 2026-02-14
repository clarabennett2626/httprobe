package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	dim    = "\033[2m"
)

type result struct {
	url      string
	status   int
	duration time.Duration
	err      error
}

func probe(url string, timeout time.Duration) result {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	client := &http.Client{Timeout: timeout}
	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		return result{url: url, err: err, duration: elapsed}
	}
	defer resp.Body.Close()

	return result{url: url, status: resp.StatusCode, duration: elapsed}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return yellow
	case code >= 400:
		return red
	default:
		return dim
	}
}

func printResult(r result) {
	if r.err != nil {
		fmt.Printf("%s✗ %-50s %sERROR%s  %s%v%s\n", red, r.url, red, reset, dim, r.err, reset)
		return
	}
	c := colorForStatus(r.status)
	ms := float64(r.duration.Microseconds()) / 1000.0
	fmt.Printf("%s✓ %-50s %s%d%s    %s%6.0fms%s\n", c, r.url, c, r.status, reset, dim, ms, reset)
}

func main() {
	timeout := 10 * time.Second
	var urls []string

	// Collect URLs from args
	if len(os.Args) > 1 {
		urls = append(urls, os.Args[1:]...)
	}

	// Collect URLs from stdin if piped
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				urls = append(urls, line)
			}
		}
	}

	if len(urls) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: httprobe <url> [url...]\n")
		fmt.Fprintf(os.Stderr, "       echo 'example.com' | httprobe\n")
		os.Exit(1)
	}

	fmt.Printf("%s%-52s %-6s %s%s\n", cyan, "URL", "STATUS", "TIME", reset)
	fmt.Println(strings.Repeat("─", 70))

	var wg sync.WaitGroup
	results := make([]result, len(urls))

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			results[i] = probe(url, timeout)
		}(i, url)
	}
	wg.Wait()

	up, down := 0, 0
	for _, r := range results {
		printResult(r)
		if r.err != nil || r.status >= 400 {
			down++
		} else {
			up++
		}
	}

	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("%s%d up%s  %s%d down%s  %s%d total%s\n",
		green, up, reset, red, down, reset, dim, len(results), reset)
}
