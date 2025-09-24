package cfg

import "time"

type Hosts struct {
	API     string
	Upload  string
	Gateway string
}

type Config struct {
	APIKey      string
	Hosts       Hosts
	UserAgent   string
	HTTPTimeout time.Duration
}

func Default() Config {
	return Config{
		Hosts: Hosts{
			API:     "https://api.lighthouse.storage",
			Upload:  "https://upload.lighthouse.storage",
			Gateway: "https://gateway.lighthouse.storage",
		},
		UserAgent:   "lighthouse-go-sdk",
		HTTPTimeout: 0,
	}
}
