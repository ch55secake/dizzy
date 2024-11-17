package client

import "time"

// Request structure will be used to send requests and later on as flags as part the command
type Request struct {
	Url       string            `json:"url"`
	Subdomain string            `json:"subdomain"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Timeout   time.Duration     `json:"timeout"`
}
