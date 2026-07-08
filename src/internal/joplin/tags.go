package joplin

import (
	"net/http"
	"net/url"
	"strconv"
)

const tagFields = "id,parent_id,title,created_time,updated_time"

// Tag mirrors the Joplin tag item.
type Tag struct {
	ID          string `json:"id"`
	ParentID    string `json:"parent_id"`
	Title       string `json:"title"`
	CreatedTime int64  `json:"created_time"`
	UpdatedTime int64  `json:"updated_time"`
}

func (c *Client) ListTags() ([]Tag, error) {
	var all []Tag
	for page := 1; ; page++ {
		q := url.Values{}
		q.Set("fields", tagFields)
		q.Set("limit", "100")
		q.Set("page", strconv.Itoa(page))
		var pg Page[Tag]
		if err := c.request(http.MethodGet, "/tags", q, nil, &pg); err != nil {
			return nil, err
		}
		all = append(all, pg.Items...)
		if !pg.HasMore {
			return all, nil
		}
	}
}

func (c *Client) CreateTag(title string) (Tag, error) {
	q := url.Values{}
	q.Set("fields", tagFields)
	var tg Tag
	err := c.request(http.MethodPost, "/tags", q, map[string]any{"title": title}, &tg)
	return tg, err
}

func (c *Client) DeleteTag(id string) error {
	return c.request(http.MethodDelete, "/tags/"+url.PathEscape(id), nil, nil, nil)
}

func (c *Client) AddTagToNote(tagID, noteID string) error {
	return c.request(http.MethodPost, "/tags/"+url.PathEscape(tagID)+"/notes", nil, map[string]any{"id": noteID}, nil)
}

func (c *Client) RemoveTagFromNote(tagID, noteID string) error {
	return c.request(http.MethodDelete, "/tags/"+url.PathEscape(tagID)+"/notes/"+url.PathEscape(noteID), nil, nil, nil)
}

func (c *Client) TagNotes(id string, limit int) (Page[Note], error) {
	q := url.Values{}
	q.Set("fields", noteFields)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var pg Page[Note]
	err := c.request(http.MethodGet, "/tags/"+url.PathEscape(id)+"/notes", q, nil, &pg)
	return pg, err
}
