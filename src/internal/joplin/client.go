package joplin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Page is the generic paginated envelope returned by Joplin list endpoints.
type Page[T any] struct {
	Items   []T  `json:"items"`
	HasMore bool `json:"has_more"`
}

// Client talks to the Joplin Web Clipper REST API.
type Client struct {
	cfg  Config
	http *http.Client
}

func NewClient(cfg Config) *Client {
	to := cfg.TimeoutSeconds
	if to <= 0 {
		to = 30
	}
	return &Client{cfg: cfg, http: &http.Client{Timeout: time.Duration(to) * time.Second}}
}

func transient(code int) bool {
	switch code {
	case 408, 425, 429, 500, 502, 503, 504:
		return true
	}
	return false
}

func backoff(attempt int, base float64) time.Duration {
	if base < 0.01 {
		base = 0.01
	}
	return time.Duration(base * float64(attempt) * float64(time.Second))
}

func apiError(status int, raw []byte) error {
	var e struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(raw, &e) == nil && e.Error != "" {
		return fmt.Errorf("joplin API error (HTTP %d): %s", status, e.Error)
	}
	return fmt.Errorf("joplin API error (HTTP %d): %s", status, strings.TrimSpace(string(raw)))
}

// request performs a JSON API call. The token is always added as a query param.
// out may be nil (response body ignored). body may be nil (no request body).
func (c *Client) request(method, path string, query url.Values, body any, out any) error {
	if query == nil {
		query = url.Values{}
	}
	query.Set("token", c.cfg.Token)
	full := strings.TrimRight(c.cfg.BaseURL, "/") + path
	if enc := query.Encode(); enc != "" {
		full += "?" + enc
	}

	var payload []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = b
	}

	retries := c.cfg.HTTPRetries
	if retries < 1 {
		retries = 1
	}
	var lastErr error
	for attempt := 1; attempt <= retries; attempt++ {
		var reader io.Reader
		if payload != nil {
			reader = bytes.NewReader(payload)
		}
		req, err := http.NewRequest(method, full, reader)
		if err != nil {
			return err
		}
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.http.Do(req)
		if err != nil {
			// Strip the *url.Error wrapper: its Error() embeds the request URL,
			// which carries ?token=<secret>. Never let that reach stderr/logs.
			if ue := new(url.Error); errors.As(err, &ue) {
				err = ue.Err
			}
			lastErr = fmt.Errorf("network error contacting Joplin: %w", err)
			if attempt < retries {
				time.Sleep(backoff(attempt, c.cfg.HTTPRetryBackoff))
			}
			continue
		}
		raw, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = fmt.Errorf("read response body: %w", readErr)
			if attempt < retries {
				time.Sleep(backoff(attempt, c.cfg.HTTPRetryBackoff))
			}
			continue
		}
		if resp.StatusCode >= 300 {
			lastErr = apiError(resp.StatusCode, raw)
			if attempt < retries && transient(resp.StatusCode) {
				time.Sleep(backoff(attempt, c.cfg.HTTPRetryBackoff))
				continue
			}
			return lastErr
		}
		if out != nil && len(raw) > 0 {
			if err := json.Unmarshal(raw, out); err != nil {
				return fmt.Errorf("parse response: %w", err)
			}
		}
		return nil
	}
	return lastErr
}

// Ping checks the Web Clipper service. Returns the raw body (expected
// "JoplinClipperServer"). No token required by the /ping endpoint.
func (c *Client) Ping() (string, error) {
	full := strings.TrimRight(c.cfg.BaseURL, "/") + "/ping"
	resp, err := c.http.Get(full)
	if err != nil {
		return "", fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", apiError(resp.StatusCode, raw)
	}
	return strings.TrimSpace(string(raw)), nil
}
