package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "List and manage tags, and tag/untag notes",
	}
	cmd.AddCommand(
		newTagListCmd(),
		newTagCreateCmd(),
		newTagDeleteCmd(),
		newTagAddCmd(),
		newTagRemoveCmd(),
		newTagNotesCmd(),
	)
	return cmd
}

func newTagListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tags",
		Long: `List every tag. Follows pagination internally.

Output (stdout, JSON):
  {"tags":[{"id","parent_id","title","created_time","updated_time"}...],"count":INT}

Exit codes:
  0  list returned
  1  auth or API error

Example:
  joplin-cli tag list`,
		Annotations: map[string]string{
			"output":  `{"tags":[{"id":"string","parent_id":"string","title":"string","created_time":"RFC3339 or null","updated_time":"RFC3339 or null"}],"count":"int"}`,
			"example": `joplin-cli tag list`,
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			tags, err := client.ListTags()
			if err != nil {
				return err
			}
			items := make([]tagOut, 0, len(tags))
			var b strings.Builder
			fmt.Fprintf(&b, "%d tag(s)", len(tags))
			for _, tg := range tags {
				items = append(items, toTagOut(tg))
				fmt.Fprintf(&b, "\n  %s  %s", tg.ID, tg.Title)
			}
			out := map[string]any{"tags": items, "count": len(items)}
			return emit(cmd, out, b.String())
		},
	}
}

func newTagCreateCmd() *cobra.Command {
	var title string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a tag",
		Long: `Create a new tag.

Input:
  --title STRING   (required) tag title

Output (stdout, JSON):
  {"id","parent_id","title","created_time","updated_time"}

Exit codes:
  0  tag created
  1  auth or API error

Example:
  joplin-cli tag create --title important`,
		Annotations: map[string]string{
			"output":  `{"id":"string","parent_id":"string","title":"string","created_time":"RFC3339 or null","updated_time":"RFC3339 or null"}`,
			"example": `joplin-cli tag create --title important`,
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			tg, err := client.CreateTag(title)
			if err != nil {
				return err
			}
			out := toTagOut(tg)
			return emit(cmd, out, fmt.Sprintf("%s  %s", out.ID, out.Title))
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "tag title (required)")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newTagDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <tag-id>",
		Short: "Delete a tag",
		Long: `Delete a tag by id. The notes themselves are not deleted.

Input:
  <tag-id>   (required) positional tag id

Output (stdout, JSON):
  {"deleted": true, "id": "<tag-id>"}

Exit codes:
  0  tag deleted
  1  not found, auth, or API error

Example:
  joplin-cli tag delete <tag-id>`,
		Annotations: map[string]string{
			"output":  `{"deleted":"bool","id":"string"}`,
			"example": `joplin-cli tag delete <tag-id>`,
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			if err := client.DeleteTag(args[0]); err != nil {
				return err
			}
			out := map[string]any{"deleted": true, "id": args[0]}
			return emit(cmd, out, "deleted "+args[0])
		},
	}
}

func newTagAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <tag-id> <note-id>",
		Short: "Attach a tag to a note",
		Long: `Associate an existing tag with a note.

Input:
  <tag-id>    (required) positional tag id
  <note-id>   (required) positional note id

Output (stdout, JSON):
  {"tagged": true, "tag_id": "<tag-id>", "note_id": "<note-id>"}

Exit codes:
  0  tag attached
  1  not found, auth, or API error

Example:
  joplin-cli tag add <tag-id> <note-id>`,
		Annotations: map[string]string{
			"output":  `{"tagged":"bool","tag_id":"string","note_id":"string"}`,
			"example": `joplin-cli tag add <tag-id> <note-id>`,
		},
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			if err := client.AddTagToNote(args[0], args[1]); err != nil {
				return err
			}
			out := map[string]any{"tagged": true, "tag_id": args[0], "note_id": args[1]}
			return emit(cmd, out, fmt.Sprintf("tagged note %s with %s", args[1], args[0]))
		},
	}
}

func newTagRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <tag-id> <note-id>",
		Short: "Detach a tag from a note",
		Long: `Remove the association between a tag and a note.

Input:
  <tag-id>    (required) positional tag id
  <note-id>   (required) positional note id

Output (stdout, JSON):
  {"untagged": true, "tag_id": "<tag-id>", "note_id": "<note-id>"}

Exit codes:
  0  tag detached
  1  not found, auth, or API error

Example:
  joplin-cli tag remove <tag-id> <note-id>`,
		Annotations: map[string]string{
			"output":  `{"untagged":"bool","tag_id":"string","note_id":"string"}`,
			"example": `joplin-cli tag remove <tag-id> <note-id>`,
		},
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			if err := client.RemoveTagFromNote(args[0], args[1]); err != nil {
				return err
			}
			out := map[string]any{"untagged": true, "tag_id": args[0], "note_id": args[1]}
			return emit(cmd, out, fmt.Sprintf("untagged note %s from %s", args[1], args[0]))
		},
	}
}

func newTagNotesCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "notes <tag-id>",
		Short: "List notes carrying a tag",
		Long: `List the notes that have a given tag.

Input:
  <tag-id>      (required) positional tag id
  --limit INT   max results, 1-100 (default 100)

Output (stdout, JSON):
  {"notes":[<note>...],"count":INT,"has_more":BOOL}

Exit codes:
  0  results returned
  1  not found, auth, or API error

Example:
  joplin-cli tag notes <tag-id>`,
		Annotations: map[string]string{
			"output":  `{"notes":[{"id":"string","title":"string","body":"string","is_todo":"bool"}],"count":"int","has_more":"bool"}`,
			"example": `joplin-cli tag notes <tag-id>`,
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			page, err := client.TagNotes(args[0], limit)
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
