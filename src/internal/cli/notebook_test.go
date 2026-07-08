package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotebookList_EmitsNotebooks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"items":[{"id":"f1","title":"Journal"}],"has_more":false}`))
	}))
	defer srv.Close()

	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"notebook", "list", "--token", "t", "--base-url", srv.URL})
	if err := root.Execute(); err != nil {
		t.Fatalf("list: %v\n%s", err, buf.String())
	}
	var out struct {
		Notebooks []folderOut `json:"notebooks"`
		Count     int         `json:"count"`
	}
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out.Count != 1 || out.Notebooks[0].Title != "Journal" {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestNotebookCreate_EmitsNotebook(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"f2","title":"Work"}`))
	}))
	defer srv.Close()

	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"notebook", "create", "--title", "Work", "--token", "t", "--base-url", srv.URL})
	if err := root.Execute(); err != nil {
		t.Fatalf("create: %v\n%s", err, buf.String())
	}
	var out folderOut
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out.ID != "f2" {
		t.Fatalf("unexpected: %+v", out)
	}
}
