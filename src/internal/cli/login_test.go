package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/artschekoff/joplin-cli/src/internal/auth"
)

func TestLogin_FromStdin_Stores(t *testing.T) {
	withFakeAuth(t)
	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetIn(strings.NewReader("tok-abcdef0123456789\n"))
	root.SetArgs([]string{"login"})
	if err := root.Execute(); err != nil {
		t.Fatalf("login: %v\n%s", err, buf.String())
	}
	got, err := auth.LoadToken()
	if err != nil || got != "tok-abcdef0123456789" {
		t.Fatalf("token not stored: %q err=%v", got, err)
	}
}

func TestLogout_RemovesToken(t *testing.T) {
	withFakeAuth(t)
	_ = auth.SaveToken("tok-abcdef0123456789")
	root := NewRootCmd()
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"logout"})
	if err := root.Execute(); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if _, err := auth.LoadToken(); err == nil {
		t.Fatal("token still present after logout")
	}
}
