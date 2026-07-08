package joplin

import (
	"net/http"
	"net/url"
	"strconv"
)

const folderFields = "id,parent_id,title,created_time,updated_time"

// Folder mirrors the Joplin folder (notebook) item.
type Folder struct {
	ID          string `json:"id"`
	ParentID    string `json:"parent_id"`
	Title       string `json:"title"`
	CreatedTime int64  `json:"created_time"`
	UpdatedTime int64  `json:"updated_time"`
}

func (c *Client) ListFolders() ([]Folder, error) {
	var all []Folder
	for page := 1; ; page++ {
		q := url.Values{}
		q.Set("fields", folderFields)
		q.Set("limit", "100")
		q.Set("page", strconv.Itoa(page))
		var pg Page[Folder]
		if err := c.request(http.MethodGet, "/folders", q, nil, &pg); err != nil {
			return nil, err
		}
		all = append(all, pg.Items...)
		if !pg.HasMore {
			return all, nil
		}
	}
}

func (c *Client) CreateFolder(title, parentID string) (Folder, error) {
	body := map[string]any{"title": title}
	if parentID != "" {
		body["parent_id"] = parentID
	}
	q := url.Values{}
	q.Set("fields", folderFields)
	var f Folder
	err := c.request(http.MethodPost, "/folders", q, body, &f)
	return f, err
}

func (c *Client) DeleteFolder(id string) error {
	return c.request(http.MethodDelete, "/folders/"+url.PathEscape(id), nil, nil, nil)
}

func (c *Client) FolderNotes(id string, limit int) (Page[Note], error) {
	q := url.Values{}
	q.Set("fields", noteFields)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var pg Page[Note]
	err := c.request(http.MethodGet, "/folders/"+url.PathEscape(id)+"/notes", q, nil, &pg)
	return pg, err
}
