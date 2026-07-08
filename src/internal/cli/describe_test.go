package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func runDescribe(t *testing.T) map[string]any {
	t.Helper()
	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"describe"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("describe output is not valid JSON: %v\n%s", err, buf.String())
	}
	return got
}

func TestDescribe_EmitsBinaryAndCommands(t *testing.T) {
	got := runDescribe(t)
	if got["binary"] != "joplin-cli" {
		t.Fatalf("binary field wrong: %v", got["binary"])
	}
	cmds, ok := got["commands"].([]any)
	if !ok || len(cmds) == 0 {
		t.Fatalf("commands array missing or empty")
	}
}

func TestDescribe_IncludesLoginAndLogout(t *testing.T) {
	got := runDescribe(t)
	names := map[string]bool{}
	for _, c := range got["commands"].([]any) {
		m := c.(map[string]any)
		names[m["name"].(string)] = true
	}
	for _, want := range []string{"login", "logout", "describe"} {
		if !names[want] {
			t.Errorf("expected %q in commands, got %v", want, names)
		}
	}
}

func TestDescribe_IncludesNestedCommandsWithOutput(t *testing.T) {
	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"describe"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	var doc struct {
		Commands []struct {
			Name   string `json:"name"`
			Output string `json:"output"`
		} `json:"commands"`
	}
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	byName := map[string]string{}
	for _, c := range doc.Commands {
		byName[c.Name] = c.Output
	}
	for _, want := range []string{
		"note search", "note get", "note create", "note update", "note delete", "note import",
		"notebook list", "notebook create", "notebook delete", "notebook notes",
		"tag list", "tag create", "tag delete", "tag add", "tag remove", "tag notes",
		"ping", "login", "logout",
	} {
		if _, ok := byName[want]; !ok {
			t.Errorf("describe missing command %q", want)
		}
	}
	// Every leaf must carry an output schema so LLMs know the response shape.
	for _, name := range []string{"note search", "note get", "tag add"} {
		if strings.TrimSpace(byName[name]) == "" {
			t.Errorf("command %q missing output schema", name)
		}
	}
}
