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

func TestWriteStructured(t *testing.T) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	doc := StructuredResult{
		URL:        "http://example.com",
		Title:      "Example Title",
		H1:         []string{"Main Header"},
		H2:         []string{"Sub Header"},
		Paragraphs: []string{"First paragraph", "Second paragraph"},
		Links:      []string{"https://link1.com", "https://link2.com"},
		Images:     []string{"/image1.jpg", "/image2.png"},
		Lists:      []string{"Item 1 | Item 2"},
	}

	if err := WriteStructured("-", doc); err != nil {
		t.Fatalf("write structured: %v", err)
	}
	w.Close()
	os.Stdout = orig

	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)
	if !strings.Contains(out, "example.com") {
		t.Fatalf("unexpected output %s", out)
	}
	if !strings.Contains(out, "Example Title") {
		t.Fatalf("title missing in output")
	}
	if !strings.Contains(out, "Main Header") {
		t.Fatalf("h1 missing in output")
	}

	tmp, err := os.CreateTemp("", "we_structured.json")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	tmp.Close()

	if err := WriteStructured(tmp.Name(), doc); err != nil {
		t.Fatalf("write structured file: %v", err)
	}

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("read temp: %v", err)
	}

	if !strings.Contains(string(data), "example.com") {
		t.Fatalf("file content wrong")
	}
	if !strings.Contains(string(data), "Example Title") {
		t.Fatalf("title missing in file")
	}

	os.Remove(tmp.Name())
}
