package lighthouse

import (
	"net/http"
	"time"
)

type Option func(*Client)

func WithAPIKey(key string) Option {
	return func(c *Client) { c.cfg.APIKey = key }
}

func WithHosts(api, upload, gateway string) Option {
	return func(c *Client) {
		c.cfg.Hosts.API = api
		c.cfg.Hosts.Upload = upload
		c.cfg.Hosts.Gateway = gateway
	}
}

func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.http = h }
}

func WithUserAgent(ua string) Option {
	return func(c *Client) { c.cfg.UserAgent = ua }
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.cfg.HTTPTimeout = d }
}
