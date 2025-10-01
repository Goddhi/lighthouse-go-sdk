package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Options struct {
	UserAgent string
	APIKey    string
}

type Client struct {
	inner *http.Client
	opt   Options
}

func New(h *http.Client, opt Options) *Client { 
	return &Client{
		inner: h,
		 opt: opt}
}

// executes a prepared *http.Request (used for streaming/multipart).
func (c *Client) Inject(req *http.Request) (*http.Response, error) {
	if c.opt.UserAgent != "" && req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", c.opt.UserAgent)
	}
	if c.opt.APIKey != "" && req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", "Bearer "+c.opt.APIKey)
	}
	return c.inner.Do(req)
}

// JSON sends optional JSON body and decodes JSON response into out.
func (c *Client) WriteJSON(ctx context.Context, method, url string, in any, out any) (*http.Response, error) {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.opt.UserAgent != "" {
		req.Header.Set("User-Agent", c.opt.UserAgent)
	}
	if c.opt.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.opt.APIKey)
	}

	res, err := c.inner.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return res, fmt.Errorf("http %d: %s", res.StatusCode, string(b))
	}
	if out != nil {
		return res, json.NewDecoder(res.Body).Decode(out)
	}
	return res, nil
}
