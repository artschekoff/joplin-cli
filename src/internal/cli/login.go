package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/artschekoff/joplin-cli/src/internal/auth"
	"github.com/spf13/cobra"
)

func newLoginCmd() *cobra.Command {
	var token string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Encrypt your Joplin API token and save it to disk",
		Long: `Save your Joplin Web Clipper token as an AES-256-GCM encrypted file. The
encryption key is derived from this machine's device id (SHA-256 of the
OS-native machine identifier). The credentials file is written to
os.UserConfigDir()/joplin-cli/credentials (mode 0600 on Unix).

Input:
  --token KEY   Pass the token inline (avoid in shell history; prefer stdin).
  stdin         If --token is omitted, read one line from stdin.

Output (stdout, JSON):
  {"status":"ok","path":"/abs/path/to/credentials"}

Exit codes:
  0  token stored
  1  empty/unreadable token, or filesystem/device-id failure

Example:
  echo "$JOPLIN_TOKEN" | joplin-cli login
  joplin-cli login --token abcdef0123456789abcdef0123456789`,
		Annotations: map[string]string{
			"output":  `{"status":"ok","path":"string (absolute path to the credentials file)"}`,
			"example": `echo "$JOPLIN_TOKEN" | joplin-cli login`,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.TrimSpace(token)
			if key == "" {
				k, err := readTokenFromStdin(cmd.InOrStdin())
				if err != nil {
					return err
				}
				key = k
			}
			if key == "" {
				return errors.New("token is empty")
			}
			if err := auth.SaveToken(key); err != nil {
				return err
			}
			path, _ := auth.CredentialsPath()
			fmt.Fprintf(cmd.OutOrStdout(), `{"status":"ok","path":%q}`+"\n", path)
			return nil
		},
	}
	cmd.Flags().StringVar(&token, "token", "", "Joplin token (falls back to stdin if empty)")
	return cmd
}

func readTokenFromStdin(r io.Reader) (string, error) {
	sc := bufio.NewScanner(r)
	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			return "", err
		}
		return "", nil
	}
	return strings.TrimSpace(sc.Text()), nil
}
