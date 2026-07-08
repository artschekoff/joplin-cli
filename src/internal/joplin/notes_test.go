package joplin

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchNotes_QueryAndFields(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("query") != "hello" || q.Get("type") != "note" {
			t.Errorf("bad query: %s", r.URL.RawQuery)
		}
		if q.Get("fields") != noteFields {
			t.Errorf("fields: %s", q.Get("fields"))
		}
		_, _ = w.Write([]byte(`{"items":[{"id":"n1","title":"Hello","is_todo":0}],"has_more":false}`))
	}))
	defer srv.Close()

	page, err := testClient(srv.URL).SearchNotes("hello", 20)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(page.Items) != 1 || page.Items[0].Title != "Hello" {
		t.Fatalf("bad result: %+v", page)
	}
}

func TestCreateNote_PostsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/notes" {
			t.Errorf("bad request line: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(raw, &got)
		if got["title"] != "T" || got["body"] != "B" || got["is_todo"] != float64(1) {
			t.Errorf("body wrong: %v", got)
		}
		_, _ = w.Write([]byte(`{"id":"n2","title":"T","body":"B","is_todo":1}`))
	}))
	defer srv.Close()

	body := "B"
	n, err := testClient(srv.URL).CreateNote(NoteCreate{Title: "T", Body: &body, IsTodo: true})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if n.ID != "n2" || n.IsTodo != 1 {
		t.Fatalf("bad note: %+v", n)
	}
}

func TestDeleteNote_Permanent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/notes/n3" {
			t.Errorf("bad delete: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("permanent") != "1" {
			t.Errorf("permanent not set: %s", r.URL.RawQuery)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	if err := testClient(srv.URL).DeleteNote("n3", true); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
