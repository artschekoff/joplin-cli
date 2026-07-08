package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRoot_HelpListsBinary(t *testing.T) {
	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"--help"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute --help: %v", err)
	}
	if !strings.Contains(buf.String(), "joplin-cli") {
		t.Fatalf("help missing binary name:\n%s", buf.String())
	}
}

func TestRoot_HasPersistentFlags(t *testing.T) {
	root := NewRootCmd()
	for _, name := range []string{"format", "token", "base-url"} {
		if root.PersistentFlags().Lookup(name) == nil {
			t.Errorf("missing persistent flag --%s", name)
		}
	}
}
