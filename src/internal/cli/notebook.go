package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newNotebookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "notebook",
		Aliases: []string{"folder"},
		Short:   "List and manage notebooks (folders)",
	}
	cmd.AddCommand(
		newNotebookListCmd(),
		newNotebookCreateCmd(),
		newNotebookDeleteCmd(),
		newNotebookNotesCmd(),
	)
	return cmd
}

func newNotebookListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all notebooks",
		Long: `List every notebook (folder). Follows pagination internally.

Output (stdout, JSON):
  {"notebooks":[{"id","parent_id","title","created_time","updated_time"}...],"count":INT}

Exit codes:
  0  list returned
  1  auth or API error

Example:
  joplin-cli notebook list`,
		Annotations: map[string]string{
			"output":  `{"notebooks":[{"id":"string","parent_id":"string","title":"string","created_time":"RFC3339 or null","updated_time":"RFC3339 or null"}],"count":"int"}`,
			"example": `joplin-cli notebook list`,
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			folders, err := client.ListFolders()
			if err != nil {
				return err
			}
			items := make([]folderOut, 0, len(folders))
			var b strings.Builder
			fmt.Fprintf(&b, "%d notebook(s)", len(folders))
			for _, f := range folders {
				items = append(items, toFolderOut(f))
				fmt.Fprintf(&b, "\n  %s  %s", f.ID, f.Title)
			}
			out := map[string]any{"notebooks": items, "count": len(items)}
			return emit(cmd, out, b.String())
		},
	}
}

func newNotebookCreateCmd() *cobra.Command {
	var title, parent string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a notebook",
		Long: `Create a notebook (folder), optionally nested under a parent.

Input:
  --title STRING    (required) notebook title
  --parent STRING   parent notebook id (for nesting)

Output (stdout, JSON):
  {"id","parent_id","title","created_time","updated_time"}

Exit codes:
  0  notebook created
  1  auth or API error

Example:
  joplin-cli notebook create --title "Projects"`,
		Annotations: map[string]string{
			"output":  `{"id":"string","parent_id":"string","title":"string","created_time":"RFC3339 or null","updated_time":"RFC3339 or null"}`,
			"example": `joplin-cli notebook create --title "Projects"`,
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			f, err := client.CreateFolder(title, parent)
			if err != nil {
				return err
			}
			out := toFolderOut(f)
			return emit(cmd, out, fmt.Sprintf("%s  %s", out.ID, out.Title))
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "notebook title (required)")
	cmd.Flags().StringVar(&parent, "parent", "", "parent notebook id")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newNotebookDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <notebook-id>",
		Short: "Delete a notebook",
		Long: `Delete a notebook (folder) by id. Notes inside it are moved to the trash by Joplin.

Input:
  <notebook-id>   (required) positional notebook id

Output (stdout, JSON):
  {"deleted": true, "id": "<notebook-id>"}

Exit codes:
  0  notebook deleted
  1  not found, auth, or API error

Example:
  joplin-cli notebook delete <notebook-id>`,
		Annotations: map[string]string{
			"output":  `{"deleted":"bool","id":"string"}`,
			"example": `joplin-cli notebook delete <notebook-id>`,
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			if err := client.DeleteFolder(args[0]); err != nil {
				return err
			}
			out := map[string]any{"deleted": true, "id": args[0]}
			return emit(cmd, out, "deleted "+args[0])
		},
	}
}

func newNotebookNotesCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "notes <notebook-id>",
		Short: "List notes inside a notebook",
		Long: `List the notes contained in a notebook (folder).

Input:
  <notebook-id>   (required) positional notebook id
  --limit INT     max results, 1-100 (default 100)

Output (stdout, JSON):
  {"notes":[<note>...],"count":INT,"has_more":BOOL}

Exit codes:
  0  results returned
  1  not found, auth, or API error

Example:
  joplin-cli notebook notes <notebook-id> --limit 50`,
		Annotations: map[string]string{
			"output":  `{"notes":[{"id":"string","title":"string","body":"string","is_todo":"bool"}],"count":"int","has_more":"bool"}`,
			"example": `joplin-cli notebook notes <notebook-id>`,
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			page, err := client.FolderNotes(args[0], limit)
			if err != nil {
				return err
			}
			out := toNoteListOut(page)
			return emit(cmd, out, formatNoteList(out))
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 100, "maximum number of results (1-100)")
	return cmd
}
