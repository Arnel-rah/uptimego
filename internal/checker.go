package checker

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CheckResult struct {
	Up      bool
	Latency time.Duration
	Error   error
}

func CheckEndpoint(url string, timeout time.Duration) CheckResult {
	return CheckWithRetry(url, timeout)
}

func CheckWithRetry(url string, timeout time.Duration) CheckResult {
	const maxRetries = 3
	const baseRetryDelay = 2 * time.Second

	var lastErr error
	var lastLatency time.Duration

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("Tentative %d/%d pour %s...\n", attempt, maxRetries, url)

		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("Échec tentative %d/%d pour %s : %v\n", attempt, maxRetries, url, err)
			lastErr = err
			if attempt < maxRetries {
				time.Sleep(baseRetryDelay * time.Duration(attempt))
				continue
			}
			return CheckResult{Up: false, Latency: lastLatency, Error: lastErr}
		}

		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Do(req)
		latency := time.Since(start)

		if err != nil {
			fmt.Printf("Échec tentative %d/%d pour %s : %v\n", attempt, maxRetries, url, err)
			lastErr = err
			lastLatency = latency
			if strings.Contains(err.Error(), "EOF") && attempt == 1 {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			if attempt < maxRetries {
				time.Sleep(baseRetryDelay * time.Duration(1<<uint(attempt-1)))
				continue
			}
			return CheckResult{Up: false, Latency: lastLatency, Error: lastErr}
		}

		if resp.StatusCode >= 400 && resp.StatusCode != 500 {
			return CheckResult{Up: false, Latency: latency, Error: fmt.Errorf("status %d - no retry", resp.StatusCode)}
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Printf("Error closing response body: %v\n", err)
			}
		}()

		up := resp.StatusCode >= 200 && resp.StatusCode < 400
		fmt.Printf("Succès à la tentative %d/%d pour %s en %d ms (status %d)\n", attempt, maxRetries, url, latency.Milliseconds(), resp.StatusCode)

		return CheckResult{
			Up:      up,
			Latency: latency,
			Error:   nil,
		}
	}

	return CheckResult{
		Up:      false,
		Latency: lastLatency,
		Error:   fmt.Errorf("abandon après %d tentatives : %v", maxRetries, lastErr),
	}
}

func FormatResult(name, url string, result CheckResult) string {
	if result.Error != nil {
		return fmt.Sprintf("%s (%s) → DOWN (%v)", name, url, result.Error)
	}
	status := "UP"
	if !result.Up {
		status = "DOWN"
	}
	return fmt.Sprintf("%s (%s) → %s (%d ms)", name, url, status, result.Latency.Milliseconds())
}
