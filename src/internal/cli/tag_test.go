package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTagAdd_EmitsTagged(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/tags/t1/notes" {
			t.Errorf("bad request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"tag", "add", "t1", "n1", "--token", "t", "--base-url", srv.URL})
	if err := root.Execute(); err != nil {
		t.Fatalf("add: %v\n%s", err, buf.String())
	}
	var out map[string]any
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out["tagged"] != true || out["tag_id"] != "t1" || out["note_id"] != "n1" {
		t.Fatalf("unexpected: %v", out)
	}
}

func TestTagList_EmitsTags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"items":[{"id":"t1","title":"work"}],"has_more":false}`))
	}))
	defer srv.Close()

	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"tag", "list", "--token", "t", "--base-url", srv.URL})
	if err := root.Execute(); err != nil {
		t.Fatalf("list: %v\n%s", err, buf.String())
	}
	var out struct {
		Tags  []tagOut `json:"tags"`
		Count int      `json:"count"`
	}
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out.Count != 1 || out.Tags[0].Title != "work" {
		t.Fatalf("unexpected: %+v", out)
	}
}
