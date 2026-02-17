package main

import (
	"bufio"
	"encoding/json"
	"flag"
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
	URL      string  `json:"url"`
	Status   int     `json:"status"`
	Duration float64 `json:"duration_ms"`
	Error    string  `json:"error,omitempty"`
}

func probe(url string, timeout time.Duration) result {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	client := &http.Client{Timeout: timeout}
	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)
	ms := float64(elapsed.Microseconds()) / 1000.0

	if err != nil {
		return result{URL: url, Error: err.Error(), Duration: ms}
	}
	defer resp.Body.Close()

	return result{URL: url, Status: resp.StatusCode, Duration: ms}
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
	if r.Error != "" {
		fmt.Printf("%s✗ %-50s %sERROR%s  %s%v%s\n", red, r.URL, red, reset, dim, r.Error, reset)
		return
	}
	c := colorForStatus(r.Status)
	fmt.Printf("%s✓ %-50s %s%d%s    %s%6.0fms%s\n", c, r.URL, c, r.Status, reset, dim, r.Duration, reset)
}

func main() {
	timeoutSec := flag.Int("t", 10, "timeout per request in seconds")
	concurrency := flag.Int("c", 10, "max concurrent requests")
	jsonOutput := flag.Bool("j", false, "output results as JSON")
	failOnly := flag.Bool("f", false, "show only failed requests (status >= 400 or errors)")
	quiet := flag.Bool("q", false, "suppress output, exit with code 1 if any probe fails")
	flag.Parse()

	timeout := time.Duration(*timeoutSec) * time.Second
	var urls []string

	// Collect URLs from positional args
	urls = append(urls, flag.Args()...)

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
		fmt.Fprintf(os.Stderr, "Usage: httprobe [flags] <url> [url...]\n")
		fmt.Fprintf(os.Stderr, "       echo 'example.com' | httprobe\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	results := make([]result, len(urls))
	sem := make(chan struct{}, *concurrency)
	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			sem <- struct{}{}
			results[i] = probe(url, timeout)
			<-sem
		}(i, url)
	}
	wg.Wait()

	up, down := 0, 0
	for _, r := range results {
		if r.Error != "" || r.Status >= 400 {
			down++
		} else {
			up++
		}
	}

	if *quiet {
		if down > 0 {
			os.Exit(1)
		}
		return
	}

	if *jsonOutput {
		filtered := results
		if *failOnly {
			filtered = nil
			for _, r := range results {
				if r.Error != "" || r.Status >= 400 {
					filtered = append(filtered, r)
				}
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(filtered)
		return
	}

	fmt.Printf("%s%-52s %-6s %s%s\n", cyan, "URL", "STATUS", "TIME", reset)
	fmt.Println(strings.Repeat("─", 70))

	for _, r := range results {
		isFail := r.Error != "" || r.Status >= 400
		if *failOnly && !isFail {
			continue
		}
		printResult(r)
	}

	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("%s%d up%s  %s%d down%s  %s%d total%s\n",
		green, up, reset, red, down, reset, dim, len(results), reset)
}
