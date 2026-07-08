package cli

import (
	"fmt"

	"github.com/artschekoff/joplin-cli/src/internal/joplin"
	"github.com/spf13/cobra"
)

func newPingCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ping",
		Short: "Check that the Joplin Web Clipper service is reachable",
		Long: `Sends GET /ping to the Joplin Web Clipper service. No token required.

Output (stdout, JSON):
  {"ok": true, "message": "JoplinClipperServer", "base_url": "http://localhost:41184"}

Exit codes:
  0  service reachable
  1  connection refused or unexpected response

Example:
  joplin-cli ping
  joplin-cli ping --base-url http://localhost:41184`,
		Annotations: map[string]string{
			"output":  `{"ok":"bool (true when body is 'JoplinClipperServer')","message":"string (raw /ping body)","base_url":"string"}`,
			"example": `joplin-cli ping`,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := joplin.LoadConfig()
			if bu, _ := cmd.Flags().GetString("base-url"); bu != "" {
				cfg.BaseURL = bu
			}
			msg, err := joplin.NewClient(cfg).Ping()
			if err != nil {
				return err
			}
			ok := msg == "JoplinClipperServer"
			out := map[string]any{"ok": ok, "message": msg, "base_url": cfg.BaseURL}
			return emit(cmd, out, fmt.Sprintf("ok=%v  %s  (%s)", ok, msg, cfg.BaseURL))
		},
	}
}
