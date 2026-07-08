package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing_ReportsOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("JoplinClipperServer"))
	}))
	defer srv.Close()

	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"ping", "--base-url", srv.URL})
	if err := root.Execute(); err != nil {
		t.Fatalf("ping: %v\n%s", err, buf.String())
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("bad json: %v\n%s", err, buf.String())
	}
	if out["ok"] != true || out["message"] != "JoplinClipperServer" {
		t.Fatalf("unexpected: %v", out)
	}
}
