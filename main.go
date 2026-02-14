package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	concurrency := flag.Int("c", 10, "concurrency level")
	timeout := flag.Int("t", 5, "timeout in seconds")
	showStatus := flag.Bool("s", false, "show status codes")
	flag.Parse()

	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	urls := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urls {
				if !strings.HasPrefix(url, "http") {
					url = "https://" + url
				}
				resp, err := client.Get(url)
				if err != nil {
					continue
				}
				resp.Body.Close()
				if *showStatus {
					fmt.Printf("[%d] %s\n", resp.StatusCode, url)
				} else {
					fmt.Println(url)
				}
			}
		}()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls <- line
		}
	}
	close(urls)
	wg.Wait()
}
