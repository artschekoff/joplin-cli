package cli

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/artschekoff/joplin-cli/src/internal/auth"
)

func TestResolveToken_EnvWins(t *testing.T) {
	withFakeAuth(t)
	_ = auth.SaveToken("stored-token")
	t.Setenv("JOPLIN_TOKEN", "env-token")
	root := NewRootCmd()
	root.SetArgs([]string{"ping"}) // any command; we only need parsed flags
	tok, err := resolveToken(rootWithFlags(t))
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if tok != "env-token" {
		t.Fatalf("env should win, got %q", tok)
	}
}

func TestResolveToken_NoneErrors(t *testing.T) {
	withFakeAuth(t)
	t.Setenv("JOPLIN_TOKEN", "")
	_, err := resolveToken(rootWithFlags(t))
	if err == nil {
		t.Fatal("expected error when no token available")
	}
}

// rootWithFlags returns a root command with persistent flags parsed (empty),
// suitable for exercising resolveToken directly.
func rootWithFlags(t *testing.T) *cobra.Command {
	t.Helper()
	r := NewRootCmd()
	r.SetArgs([]string{})
	_ = r.ParseFlags([]string{})
	return r
}
