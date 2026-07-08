package cli

import (
	"fmt"

	"github.com/artschekoff/joplin-cli/src/internal/auth"
	"github.com/spf13/cobra"
)

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove the encrypted credentials file",
		Long: "Delete the credentials file written by `login`. Safe to run when no token is stored — no error is returned.\n\n" +
			"Output (stdout, JSON):\n  {\"status\":\"ok\"}\n\n" +
			"Exit codes:\n  0  file removed or was already absent\n  1  filesystem error",
		Annotations: map[string]string{
			"output":  `{"status":"ok"}`,
			"example": `joplin-cli logout`,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.DeleteToken(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), `{"status":"ok"}`)
			return nil
		},
	}
}
