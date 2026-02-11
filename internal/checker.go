package checker

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	const retryDelay = 2 * time.Second

	var lastErr error
	var lastLatency time.Duration

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("Tentative %d/%d pour %s...\n", attempt, maxRetries, url)

		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return CheckResult{Up: false, Latency: 0, Error: lastErr}
		}

		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Do(req)
		latency := time.Since(start)

		if err != nil {
			lastErr = err
			lastLatency = latency
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return CheckResult{Up: false, Latency: lastLatency, Error: lastErr}
		}

		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				fmt.Printf("Error closing response body: %v\n", err)
			}
		}(resp.Body)

		up := resp.StatusCode >= 200 && resp.StatusCode < 400

		fmt.Printf("Succès à la tentative %d/%d pour %s en %d ms\n", attempt, maxRetries, url, latency.Milliseconds())

		return CheckResult{
			Up:      up,
			Latency: latency,
			Error:   nil,
		}
	}

	return CheckResult{
		Up:      false,
		Latency: 0,
		Error:   fmt.Errorf("all %d retries failed", maxRetries),
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
