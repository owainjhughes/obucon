package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	defaultHost    = "127.0.0.1"
	defaultPort    = "8080"
	defaultTimeout = 3 * time.Second
	endpointPath   = "/health"
)

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func main() {
	host := envOr("HEALTHCHECK_HOST", defaultHost)
	port := envOr("PORT", defaultPort)
	url := fmt.Sprintf("http://%s:%s%s", host, port, endpointPath)

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: request to %s failed: %v\n", url, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "healthcheck: %s returned status %d\n", url, resp.StatusCode)
		os.Exit(1)
	}
}
