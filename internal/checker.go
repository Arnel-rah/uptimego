package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type CheckResult struct {
	Up      bool
	Latency time.Duration
	Error   error
}

func CheckEndpoint(url string, timeout time.Duration) CheckResult {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return CheckResult{Up: false, Latency: 0, Error: err}
	}

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return CheckResult{Up: false, Latency: latency, Error: err}
	}
	defer resp.Body.Close()
	up := resp.StatusCode >= 200 && resp.StatusCode < 400

	return CheckResult{
		Up:      up,
		Latency: latency,
		Error:   nil,
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
