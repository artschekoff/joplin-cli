package joplin

import (
	"net/http"
	"net/url"
	"strconv"
)

const noteFields = "id,parent_id,title,body,created_time,updated_time,is_todo"

// Note mirrors the Joplin note item. Timestamps are epoch milliseconds.
type Note struct {
	ID          string `json:"id"`
	ParentID    string `json:"parent_id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	CreatedTime int64  `json:"created_time"`
	UpdatedTime int64  `json:"updated_time"`
	IsTodo      int    `json:"is_todo"`
}

type NoteCreate struct {
	Title    string
	Body     *string
	ParentID *string
	IsTodo   bool
}

type NoteUpdate struct {
	Title    *string
	Body     *string
	ParentID *string
	IsTodo   *bool
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (c *Client) SearchNotes(query string, limit int) (Page[Note], error) {
	q := url.Values{}
	q.Set("query", query)
	q.Set("type", "note")
	q.Set("fields", noteFields)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var page Page[Note]
	err := c.request(http.MethodGet, "/search", q, nil, &page)
	return page, err
}

func (c *Client) GetNote(id string) (Note, error) {
	q := url.Values{}
	q.Set("fields", noteFields)
	var n Note
	err := c.request(http.MethodGet, "/notes/"+url.PathEscape(id), q, nil, &n)
	return n, err
}

func (c *Client) CreateNote(in NoteCreate) (Note, error) {
	body := map[string]any{"title": in.Title, "is_todo": boolToInt(in.IsTodo)}
	if in.Body != nil {
		body["body"] = *in.Body
	}
	if in.ParentID != nil {
		body["parent_id"] = *in.ParentID
	}
	q := url.Values{}
	q.Set("fields", noteFields)
	var n Note
	err := c.request(http.MethodPost, "/notes", q, body, &n)
	return n, err
}

func (c *Client) UpdateNote(id string, in NoteUpdate) (Note, error) {
	body := map[string]any{}
	if in.Title != nil {
		body["title"] = *in.Title
	}
	if in.Body != nil {
		body["body"] = *in.Body
	}
	if in.ParentID != nil {
		body["parent_id"] = *in.ParentID
	}
	if in.IsTodo != nil {
		body["is_todo"] = boolToInt(*in.IsTodo)
	}
	q := url.Values{}
	q.Set("fields", noteFields)
	var n Note
	err := c.request(http.MethodPut, "/notes/"+url.PathEscape(id), q, body, &n)
	return n, err
}

func (c *Client) DeleteNote(id string, permanent bool) error {
	q := url.Values{}
	if permanent {
		q.Set("permanent", "1")
	}
	return c.request(http.MethodDelete, "/notes/"+url.PathEscape(id), q, nil, nil)
}
