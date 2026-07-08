package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNoteCreate_EmitsNote(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"n1","title":"Hi","body":"B","is_todo":0}`))
	}))
	defer srv.Close()

	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"note", "create", "--title", "Hi", "--body", "B",
		"--token", "test-token", "--base-url", srv.URL})
	if err := root.Execute(); err != nil {
		t.Fatalf("create: %v\n%s", err, buf.String())
	}
	var out noteOut
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("bad json: %v\n%s", err, buf.String())
	}
	if out.ID != "n1" || out.IsTodo != false {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestNoteSearch_EmitsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"items":[{"id":"n1","title":"A"}],"has_more":false}`))
	}))
	defer srv.Close()

	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"note", "search", "A", "--token", "test-token", "--base-url", srv.URL})
	if err := root.Execute(); err != nil {
		t.Fatalf("search: %v\n%s", err, buf.String())
	}
	var out noteListOut
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out.Count != 1 || out.Notes[0].ID != "n1" {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestNoteImport_ParsesTitleFromHeading(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"n9","title":"My Title","body":"Line one"}`))
	}))
	defer srv.Close()

	dir := t.TempDir()
	fp := filepath.Join(dir, "note.md")
	_ = os.WriteFile(fp, []byte("# My Title\nLine one"), 0o644)

	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"note", "import", fp, "--token", "test-token", "--base-url", srv.URL})
	if err := root.Execute(); err != nil {
		t.Fatalf("import: %v\n%s", err, buf.String())
	}
	var out struct {
		Note         noteOut `json:"note"`
		ImportedFrom string  `json:"imported_from"`
	}
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out.Note.ID != "n9" || out.ImportedFrom != fp {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestParseMarkdownFile(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "x.md")
	_ = os.WriteFile(fp, []byte("# Heading\nbody text"), 0o644)
	title, body, err := parseMarkdownFile(fp)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if title != "Heading" || body != "body text" {
		t.Fatalf("got title=%q body=%q", title, body)
	}
}
