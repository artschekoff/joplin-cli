package joplin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListFolders_Paginates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "1":
			_, _ = w.Write([]byte(`{"items":[{"id":"f1","title":"A"}],"has_more":true}`))
		default:
			_, _ = w.Write([]byte(`{"items":[{"id":"f2","title":"B"}],"has_more":false}`))
		}
	}))
	defer srv.Close()

	got, err := testClient(srv.URL).ListFolders()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 || got[0].ID != "f1" || got[1].ID != "f2" {
		t.Fatalf("pagination failed: %+v", got)
	}
}

func TestCreateFolder_WithParent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/folders" {
			t.Errorf("bad request: %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"f3","title":"C","parent_id":"f1"}`))
	}))
	defer srv.Close()

	f, err := testClient(srv.URL).CreateFolder("C", "f1")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if f.ID != "f3" || f.ParentID != "f1" {
		t.Fatalf("bad folder: %+v", f)
	}
}
