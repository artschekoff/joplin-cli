package joplin

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListTags_Paginates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") == "1" {
			_, _ = w.Write([]byte(`{"items":[{"id":"t1","title":"work"}],"has_more":false}`))
			return
		}
		_, _ = w.Write([]byte(`{"items":[],"has_more":false}`))
	}))
	defer srv.Close()

	got, err := testClient(srv.URL).ListTags()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 1 || got[0].Title != "work" {
		t.Fatalf("bad tags: %+v", got)
	}
}

func TestAddTagToNote_PostsNoteID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/tags/t1/notes" {
			t.Errorf("bad request: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(raw, &got)
		if got["id"] != "n1" {
			t.Errorf("body wrong: %v", got)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	if err := testClient(srv.URL).AddTagToNote("t1", "n1"); err != nil {
		t.Fatalf("add: %v", err)
	}
}

func TestRemoveTagFromNote_Path(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/tags/t1/notes/n1" {
			t.Errorf("bad request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	if err := testClient(srv.URL).RemoveTagFromNote("t1", "n1"); err != nil {
		t.Fatalf("remove: %v", err)
	}
}
