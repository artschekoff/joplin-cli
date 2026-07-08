package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/artschekoff/joplin-cli/src/internal/joplin"
	"github.com/spf13/cobra"
)

// emit writes text when --format=text, otherwise pretty JSON of jsonValue.
func emit(cmd *cobra.Command, jsonValue any, text string) error {
	format, _ := cmd.Flags().GetString("format")
	if format == "text" {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), text)
		return err
	}
	b, err := json.MarshalIndent(jsonValue, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), string(b))
	return err
}

func msToISO(ms int64) *string {
	if ms == 0 {
		return nil
	}
	s := time.UnixMilli(ms).UTC().Format(time.RFC3339)
	return &s
}

type noteOut struct {
	ID          string  `json:"id"`
	ParentID    string  `json:"parent_id,omitempty"`
	Title       string  `json:"title"`
	Body        string  `json:"body"`
	CreatedTime *string `json:"created_time"`
	UpdatedTime *string `json:"updated_time"`
	IsTodo      bool    `json:"is_todo"`
}

func toNoteOut(n joplin.Note) noteOut {
	return noteOut{
		ID:          n.ID,
		ParentID:    n.ParentID,
		Title:       n.Title,
		Body:        n.Body,
		CreatedTime: msToISO(n.CreatedTime),
		UpdatedTime: msToISO(n.UpdatedTime),
		IsTodo:      n.IsTodo != 0,
	}
}

type noteListOut struct {
	Notes   []noteOut `json:"notes"`
	Count   int       `json:"count"`
	HasMore bool      `json:"has_more"`
}

func toNoteListOut(page joplin.Page[joplin.Note]) noteListOut {
	out := noteListOut{Notes: make([]noteOut, 0, len(page.Items)), Count: len(page.Items), HasMore: page.HasMore}
	for _, n := range page.Items {
		out.Notes = append(out.Notes, toNoteOut(n))
	}
	return out
}

type folderOut struct {
	ID          string  `json:"id"`
	ParentID    string  `json:"parent_id,omitempty"`
	Title       string  `json:"title"`
	CreatedTime *string `json:"created_time"`
	UpdatedTime *string `json:"updated_time"`
}

func toFolderOut(f joplin.Folder) folderOut {
	return folderOut{ID: f.ID, ParentID: f.ParentID, Title: f.Title, CreatedTime: msToISO(f.CreatedTime), UpdatedTime: msToISO(f.UpdatedTime)}
}

type tagOut struct {
	ID          string  `json:"id"`
	ParentID    string  `json:"parent_id,omitempty"`
	Title       string  `json:"title"`
	CreatedTime *string `json:"created_time"`
	UpdatedTime *string `json:"updated_time"`
}

func toTagOut(tg joplin.Tag) tagOut {
	return tagOut{ID: tg.ID, ParentID: tg.ParentID, Title: tg.Title, CreatedTime: msToISO(tg.CreatedTime), UpdatedTime: msToISO(tg.UpdatedTime)}
}

func formatNote(n noteOut) string {
	todo := "no"
	if n.IsTodo {
		todo = "yes"
	}
	return fmt.Sprintf("ID:    %s\nTitle: %s\nTodo:  %s\n\n%s", n.ID, n.Title, todo, n.Body)
}

func formatNoteList(l noteListOut) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d note(s) (has_more=%v)", l.Count, l.HasMore)
	for _, n := range l.Notes {
		fmt.Fprintf(&b, "\n  %s  %s", n.ID, n.Title)
	}
	return b.String()
}
