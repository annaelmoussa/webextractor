package fetcher

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"webextractor/internal/htmlparser"
	"webextractor/internal/types"
)

// Fetcher encapsule la logique du client HTTP.
type Fetcher struct {
	client *http.Client
	userAgent types.UserAgent
}

// New retourne un Fetcher avec le timeout donné.
func New(timeout time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		userAgent: types.DefaultUserAgent,
	}
}

// NewWithUserAgent retourne un Fetcher avec un User-Agent personnalisé.
func NewWithUserAgent(timeout time.Duration, userAgent types.UserAgent) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		userAgent: userAgent,
	}
}

// Fetch récupère la page située à l'URL et analyse le corps comme HTML.
func (f *Fetcher) Fetch(url string) (*htmlparser.Node, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", f.userAgent.String())

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	doc, err := htmlparser.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("html parse error: %w", err)
	}
	if doc == nil {
		return nil, errors.New("empty document")
	}
	return doc, nil
}
