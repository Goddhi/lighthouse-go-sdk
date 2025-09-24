package lighthouse

import (
	"time"

	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/cfg"
)

// Re-export for public API convenience.
type Hosts = cfg.Hosts
type Config = cfg.Config

// DefaultConfig returns SDK defaults (override as needed via options).
func DefaultConfig() Config {
	c := cfg.Default()
	if c.UserAgent == "" {
		c.UserAgent = "lighthouse-go-sdk"
	}
	if c.HTTPTimeout == 0 {
		c.HTTPTimeout = 30 * time.Second
	}
	return c
}
