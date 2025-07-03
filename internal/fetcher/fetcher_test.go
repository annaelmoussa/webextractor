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
