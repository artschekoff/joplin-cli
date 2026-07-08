package cli

import "github.com/spf13/cobra"

var version = "dev"

// SetVersion is called from main() with the ldflags-injected value.
func SetVersion(v string) { version = v }

// NewRootCmd builds a fresh command tree — used by tests and Execute.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "joplin-cli",
		Short:   "CLI for the Joplin note API",
		Long:    "joplin-cli — read and write Joplin notes, notebooks and tags via the Web Clipper REST API. Authenticate with `login` (or JOPLIN_TOKEN); every command emits JSON on stdout for machine consumption. Run `describe` for the full machine-readable schema.",
		Version: version,
		// Runtime errors print to stderr from main(); do not dump usage.
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	// Keep the built-in shell-completion command out of `describe` and help:
	// it carries no output schema and is noise in the machine-readable contract.
	root.CompletionOptions.HiddenDefaultCmd = true
	root.PersistentFlags().String("format", "json", "output format: json|text")
	root.PersistentFlags().String("token", "", "Joplin API token (overrides JOPLIN_TOKEN)")
	root.PersistentFlags().String("base-url", "", "Joplin API base URL (overrides JOPLIN_BASE_URL)")
	root.AddCommand(newLoginCmd())
	root.AddCommand(newLogoutCmd())
	root.AddCommand(newDescribeCmd())
	root.AddCommand(newPingCmd())
	root.AddCommand(newNoteCmd())
	root.AddCommand(newNotebookCmd())
	root.AddCommand(newTagCmd())
	return root
}

// Execute is the entry point called from main().
func Execute() error {
	return NewRootCmd().Execute()
}
