package httputil

import (
	"net/http"
	"net/url"
	"os"
	"time"
)

func NewHTTPClient(timeout time.Duration) *http.Client {
	transport := &http.Transport{
		Proxy: proxyFromEnv(),
	}
	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

func proxyFromEnv() func(*http.Request) (*url.URL, error) {
	proxyURL := os.Getenv("HTTP_PROXY")
	if proxyURL == "" {
		proxyURL = os.Getenv("HTTPS_PROXY")
	}
	if proxyURL == "" {
		return http.ProxyFromEnvironment
	}
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return nil
	}
	return http.ProxyURL(parsed)
}
