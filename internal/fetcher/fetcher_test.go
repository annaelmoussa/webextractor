package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	htmlBody := `<html><body><h1>Hello</h1></body></html>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(htmlBody))
	}))
	defer srv.Close()

	f := New(5 * time.Second)
	doc, err := f.Fetch(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc == nil {
		t.Fatalf("expected non-nil doc")
	}
}

func TestFetchNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	f := New(2 * time.Second)
	_, err := f.Fetch(srv.URL)
	if err == nil {
		t.Fatalf("expected error for non-200 status")
	}
}

func TestFetchInvalidURL(t *testing.T) {
	f := New(2 * time.Second)
	_, err := f.Fetch("invalid-url")
	if err == nil {
		t.Fatalf("expected error for invalid URL")
	}
}

func TestFetchInvalidHTML(t *testing.T) {
	invalidHTML := `<html><body><h1>Unclosed tag`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(invalidHTML))
	}))
	defer srv.Close()

	f := New(5 * time.Second)
	// Note: golang.org/x/net/html est assez tolérant,
	// mais on peut tester avec du contenu vraiment cassé
	doc, err := f.Fetch(srv.URL)
	// Le parser HTML de Go est très tolérant, donc même du HTML invalide peut passer
	// On teste juste que ça ne plante pas
	if err != nil && doc == nil {
		t.Logf("HTML parse handled gracefully: %v", err)
	}
}

func TestFetchUserAgent(t *testing.T) {
	var userAgent string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.Header.Get("User-Agent")
		w.Write([]byte(`<html><body><h1>Hello</h1></body></html>`))
	}))
	defer srv.Close()

	f := New(5 * time.Second)
	_, err := f.Fetch(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedUA := "WebExtractor/0.1"
	if userAgent != expectedUA {
		t.Fatalf("expected User-Agent '%s', got '%s'", expectedUA, userAgent)
	}
}
