package joplin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testClient(url string) *Client {
	return NewClient(Config{Token: "test-token", BaseURL: url, TimeoutSeconds: 5, HTTPRetries: 2, HTTPRetryBackoff: 0.01})
}

func TestPing_ReturnsServerString(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ping" {
			t.Errorf("path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte("JoplinClipperServer"))
	}))
	defer srv.Close()

	got, err := testClient(srv.URL).Ping()
	if err != nil {
		t.Fatalf("ping: %v", err)
	}
	if got != "JoplinClipperServer" {
		t.Fatalf("got %q", got)
	}
}

func TestRequest_SendsTokenAndDecodes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("token") != "test-token" {
			t.Errorf("missing token query param: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[{"id":"a"}],"has_more":false}`))
	}))
	defer srv.Close()

	var page Page[map[string]any]
	if err := testClient(srv.URL).request(http.MethodGet, "/notes", nil, nil, &page); err != nil {
		t.Fatalf("request: %v", err)
	}
	if len(page.Items) != 1 || page.Items[0]["id"] != "a" {
		t.Fatalf("decode failed: %+v", page)
	}
}

func TestRequest_APIError_Surfaced(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"error":"note not found"}`))
	}))
	defer srv.Close()

	err := testClient(srv.URL).request(http.MethodGet, "/notes/x", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !containsAll(err.Error(), "404", "note not found") {
		t.Fatalf("error text missing detail: %v", err)
	}
}

func TestRequest_NetworkError_DoesNotLeakToken(t *testing.T) {
	c := NewClient(Config{Token: "SECRET-TOKEN-123", BaseURL: "http://127.0.0.1:1", TimeoutSeconds: 1, HTTPRetries: 1, HTTPRetryBackoff: 0.01})
	err := c.request(http.MethodGet, "/notes", nil, nil, nil)
	if err == nil {
		t.Fatal("expected a network error")
	}
	if strings.Contains(err.Error(), "SECRET-TOKEN-123") {
		t.Fatalf("token leaked in error message: %v", err)
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		found := false
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
