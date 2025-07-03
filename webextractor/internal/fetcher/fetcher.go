package fetcher

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

// Fetcher encapsulates HTTP client logic.
// It exposes a simple API to download an URL and return the root *html.Node.
// The HTTP client is configured with a timeout and a custom User-Agent.
// No other responsibilities live here.

type Fetcher struct {
	client *http.Client
}

// New returns a Fetcher with the given timeout.
func New(timeout time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Fetch retrieves the page located at url and parses the body as HTML.
func (f *Fetcher) Fetch(url string) (*html.Node, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "WebExtractor/0.1")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("html parse error: %w", err)
	}
	if doc == nil {
		return nil, errors.New("empty document")
	}
	return doc, nil
}
