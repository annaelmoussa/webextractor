package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><body><div>Hello</div></body></html>`))
	}))
	defer srv.Close()

	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	origArgs := os.Args
	os.Args = []string{"cmd", "-url", srv.URL, "-sel", "div"}

	main()

	w.Close()
	os.Stdout = origStdout
	os.Args = origArgs

	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)
	if !strings.Contains(out, "Hello") {
		t.Fatalf("output not correct: %s", out)
	}
}
