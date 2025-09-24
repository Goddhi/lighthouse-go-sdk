package lighthouse

import (
	"net/http"

	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/files"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/cfg"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/httpx"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/storage"
)

type Client struct {
	http *http.Client
	cfg  Config

	storage StorageService
	files   FilesService
}

func NewClient(h *http.Client, options ...Option) *Client {
	c := &Client{cfg: DefaultConfig()}

	for _, opt := range options {
		opt(c)
	}
	if c.http == nil {
		if h != nil {
			c.http = h
		} else {
			c.http = &http.Client{Timeout: c.cfg.HTTPTimeout}
		}
	}

	hx := httpx.New(c.http, httpx.Options{
		UserAgent: c.cfg.UserAgent,
		APIKey:    c.cfg.APIKey,
	})

	// adapt to cfg.Config for subpackages
	cc := cfg.Config(c.cfg)

	c.storage = storage.New(hx, cc)
	c.files = files.New(hx, cc)
	return c
}

func (c *Client) Storage() StorageService { return c.storage }
func (c *Client) Files() FilesService     { return c.files }
