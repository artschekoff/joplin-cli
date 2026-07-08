package cli

import (
	"errors"
	"os"

	"github.com/artschekoff/joplin-cli/src/internal/auth"
	"github.com/artschekoff/joplin-cli/src/internal/joplin"
	"github.com/spf13/cobra"
)

// resolveToken picks the token: --token flag > JOPLIN_TOKEN env > stored credentials.
func resolveToken(cmd *cobra.Command) (string, error) {
	if t, _ := cmd.Flags().GetString("token"); t != "" {
		return t, nil
	}
	if t := os.Getenv("JOPLIN_TOKEN"); t != "" {
		return t, nil
	}
	t, err := auth.LoadToken()
	if err == nil {
		return t, nil
	}
	if errors.Is(err, auth.ErrNoCredentials) {
		return "", errors.New("no Joplin token found — set JOPLIN_TOKEN, pass --token, or run `joplin-cli login`")
	}
	return "", err
}

// newClient builds a configured client, requiring a resolvable token.
func newClient(cmd *cobra.Command) (*joplin.Client, error) {
	tok, err := resolveToken(cmd)
	if err != nil {
		return nil, err
	}
	cfg := joplin.LoadConfig()
	cfg.Token = tok
	if bu, _ := cmd.Flags().GetString("base-url"); bu != "" {
		cfg.BaseURL = bu
	}
	return joplin.NewClient(cfg), nil
}
