package io

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	doc := DocumentResult{
		URL:     "http://example.com",
		Results: []Result{{Selector: "p", Matches: []string{"hello"}}},
	}
	if err := Write("-", doc); err != nil {
		t.Fatalf("write: %v", err)
	}
	w.Close()
	os.Stdout = orig

	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)
	if !strings.Contains(out, "example.com") {
		t.Fatalf("unexpected output %s", out)
	}

	tmp, err := os.CreateTemp("", "we.json")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	tmp.Close()
	if err := Write(tmp.Name(), doc); err != nil {
		t.Fatalf("write file: %v", err)
	}
	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("read temp: %v", err)
	}
	if !strings.Contains(string(data), "example.com") {
		t.Fatalf("file content wrong")
	}
}
