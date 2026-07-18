package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/artschekoff/joplin-cli/src/internal/joplin"
	"github.com/spf13/cobra"
)

func newNoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "note",
		Short: "Search, read, and edit Joplin notes",
	}
	cmd.AddCommand(
		newNoteSearchCmd(),
		newNoteGetCmd(),
		newNoteCreateCmd(),
		newNoteUpdateCmd(),
		newNoteDeleteCmd(),
		newNoteImportCmd(),
	)
	return cmd
}

func newNoteSearchCmd() *cobra.Command {
	var limit int
	var queryFlag string
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Full-text search notes",
		Long: `Search all notes with Joplin's search syntax (supports filters like
'tag:work', 'notebook:Journal', 'updated:day-7').

Input:
  <query>        search string (positional; or use --query)
  --query STRING search string (alias for the positional arg)
  --limit INT    max results, 1-100 (default 100)

Output (stdout, JSON):
  {"notes":[<note>...],"count":INT,"has_more":BOOL}
  <note> = {"id","parent_id","title","body","created_time","updated_time","is_todo"}

Exit codes:
  0  results returned (possibly empty)
  1  auth or API error

Example:
  joplin-cli note search "meeting" --limit 20`,
		Annotations: map[string]string{
			"output":  `{"notes":[{"id":"string","parent_id":"string","title":"string","body":"string","created_time":"RFC3339 or null","updated_time":"RFC3339 or null","is_todo":"bool"}],"count":"int","has_more":"bool"}`,
			"example": `joplin-cli note search "meeting" --limit 20`,
		},
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := queryFlag
			if len(args) == 1 {
				query = args[0]
			}
			if query == "" {
				return fmt.Errorf("a search query is required (pass it positionally or with --query)")
			}
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			page, err := client.SearchNotes(query, limit)
			if err != nil {
				return err
			}
			out := toNoteListOut(page)
			return emit(cmd, out, formatNoteList(out))
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 100, "maximum number of results (1-100)")
	cmd.Flags().StringVarP(&queryFlag, "query", "q", "", "search string (alias for the positional arg)")
	return cmd
}

func newNoteGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <note-id>",
		Short: "Retrieve a single note by id",
		Long: `Fetch one note including its Markdown body.

Input:
  <note-id>      (required) positional note id

Output (stdout, JSON):
  {"id","parent_id","title","body","created_time","updated_time","is_todo"}

Exit codes:
  0  note returned
  1  not found, auth, or API error

Example:
  joplin-cli note get 0123456789abcdef0123456789abcdef`,
		Annotations: map[string]string{
			"output":  `{"id":"string","parent_id":"string","title":"string","body":"string","created_time":"RFC3339 or null","updated_time":"RFC3339 or null","is_todo":"bool"}`,
			"example": `joplin-cli note get 0123456789abcdef0123456789abcdef`,
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			n, err := client.GetNote(args[0])
			if err != nil {
				return err
			}
			out := toNoteOut(n)
			return emit(cmd, out, formatNote(out))
		},
	}
}

func newNoteCreateCmd() *cobra.Command {
	var title, body, notebook string
	var todo bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new note",
		Long: `Create a note. Body is Markdown.

Input:
  --title STRING     (required) note title
  --body STRING      Markdown body
  --notebook STRING  parent notebook (folder) id
  --todo             create as a to-do item

Output (stdout, JSON):
  {"id","parent_id","title","body","created_time","updated_time","is_todo"}

Exit codes:
  0  note created
  1  auth or API error

Example:
  joplin-cli note create --title "Idea" --body "- buy milk" --notebook <folder-id>`,
		Annotations: map[string]string{
			"output":  `{"id":"string","parent_id":"string","title":"string","body":"string","created_time":"RFC3339 or null","updated_time":"RFC3339 or null","is_todo":"bool"}`,
			"example": `joplin-cli note create --title "Idea" --body "- buy milk"`,
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			in := joplin.NoteCreate{Title: title, IsTodo: todo}
			if cmd.Flags().Changed("body") {
				in.Body = &body
			}
			if cmd.Flags().Changed("notebook") {
				in.ParentID = &notebook
			}
			n, err := client.CreateNote(in)
			if err != nil {
				return err
			}
			out := toNoteOut(n)
			return emit(cmd, out, formatNote(out))
		},
	}
	f := cmd.Flags()
	f.StringVar(&title, "title", "", "note title (required)")
	f.StringVar(&body, "body", "", "note body in Markdown")
	f.StringVar(&notebook, "notebook", "", "parent notebook (folder) id")
	f.BoolVar(&todo, "todo", false, "create as a to-do item")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newNoteUpdateCmd() *cobra.Command {
	var title, body, notebook string
	var todo bool
	cmd := &cobra.Command{
		Use:   "update <note-id>",
		Short: "Update fields on an existing note",
		Long: `Update one or more fields. Only flags you pass are changed.

Input:
  <note-id>          (required) positional note id
  --title STRING     new title
  --body STRING      new Markdown body
  --notebook STRING  move to this notebook (folder) id
  --todo             set to-do status (pair with --todo=false to unset)

Output (stdout, JSON):
  {"id","parent_id","title","body","created_time","updated_time","is_todo"}

Exit codes:
  0  note updated
  1  not found, auth, or API error

Example:
  joplin-cli note update <note-id> --title "New title"`,
		Annotations: map[string]string{
			"output":  `{"id":"string","parent_id":"string","title":"string","body":"string","created_time":"RFC3339 or null","updated_time":"RFC3339 or null","is_todo":"bool"}`,
			"example": `joplin-cli note update <note-id> --title "New title"`,
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			var in joplin.NoteUpdate
			if cmd.Flags().Changed("title") {
				in.Title = &title
			}
			if cmd.Flags().Changed("body") {
				in.Body = &body
			}
			if cmd.Flags().Changed("notebook") {
				in.ParentID = &notebook
			}
			if cmd.Flags().Changed("todo") {
				in.IsTodo = &todo
			}
			n, err := client.UpdateNote(args[0], in)
			if err != nil {
				return err
			}
			out := toNoteOut(n)
			return emit(cmd, out, formatNote(out))
		},
	}
	f := cmd.Flags()
	f.StringVar(&title, "title", "", "new title")
	f.StringVar(&body, "body", "", "new Markdown body")
	f.StringVar(&notebook, "notebook", "", "new parent notebook (folder) id")
	f.BoolVar(&todo, "todo", false, "set to-do status")
	return cmd
}

func newNoteDeleteCmd() *cobra.Command {
	var permanent bool
	cmd := &cobra.Command{
		Use:   "delete <note-id>",
		Short: "Delete a note (to trash, or permanently)",
		Long: `Move a note to the trash, or delete it permanently with --permanent.

Input:
  <note-id>     (required) positional note id
  --permanent   delete permanently instead of moving to trash

Output (stdout, JSON):
  {"deleted": true, "id": "<note-id>", "permanent": BOOL}

Exit codes:
  0  note deleted
  1  not found, auth, or API error

Example:
  joplin-cli note delete <note-id>
  joplin-cli note delete <note-id> --permanent`,
		Annotations: map[string]string{
			"output":  `{"deleted":"bool","id":"string","permanent":"bool"}`,
			"example": `joplin-cli note delete <note-id> --permanent`,
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			if err := client.DeleteNote(args[0], permanent); err != nil {
				return err
			}
			out := map[string]any{"deleted": true, "id": args[0], "permanent": permanent}
			text := fmt.Sprintf("deleted %s (permanent=%v)", args[0], permanent)
			return emit(cmd, out, text)
		},
	}
	cmd.Flags().BoolVar(&permanent, "permanent", false, "delete permanently instead of trashing")
	return cmd
}

func newNoteImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <file.md>",
		Short: "Import a Markdown file as a new note",
		Long: `Read a Markdown file and create a note from it. If the first line is an
'# H1' heading it becomes the title and is stripped from the body; otherwise
the filename (without extension) is the title.

Input:
  <file.md>     (required) path to a UTF-8 Markdown file

Output (stdout, JSON):
  {"note": <note>, "imported_from": "<path>"}

Exit codes:
  0  note created
  1  file missing/empty, auth, or API error

Example:
  joplin-cli note import ./meeting.md`,
		Annotations: map[string]string{
			"output":  `{"note":{"id":"string","parent_id":"string","title":"string","body":"string","created_time":"RFC3339 or null","updated_time":"RFC3339 or null","is_todo":"bool"},"imported_from":"string (path)"}`,
			"example": `joplin-cli note import ./meeting.md`,
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			title, body, err := parseMarkdownFile(args[0])
			if err != nil {
				return err
			}
			n, err := client.CreateNote(joplin.NoteCreate{Title: title, Body: &body})
			if err != nil {
				return err
			}
			note := toNoteOut(n)
			out := struct {
				Note         noteOut `json:"note"`
				ImportedFrom string  `json:"imported_from"`
			}{Note: note, ImportedFrom: args[0]}
			return emit(cmd, out, formatNote(note)+"\n\nimported_from: "+args[0])
		},
	}
	return cmd
}

// parseMarkdownFile extracts a title and body from a Markdown file, matching
// the behaviour of the original MCP's MarkdownContent.from_file.
func parseMarkdownFile(path string) (string, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", "", err
	}
	if info.IsDir() {
		return "", "", fmt.Errorf("not a file: %s", path)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	content := string(raw)
	if strings.TrimSpace(content) == "" {
		return "", "", errors.New("file is empty: " + path)
	}
	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimSpace(line[2:])
			content = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
			break
		}
	}
	return title, content, nil
}
