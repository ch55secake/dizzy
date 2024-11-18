package client

// Request structure will be used to send requests and later on as flags as part the command
type Request struct {
	Url       string `json:"url"`
	Subdomain string `json:"subdomain"`
}

// EmptyRequest used for when wordlist has no data should be attached to an error
func EmptyRequest(url string) Request {
	return Request{
		Url: url,
	}
}

// isValid is used to determine whether a request can actually reach a valid subdomain/endpoint
func (req Request) isValid() bool {
	if req.Subdomain == "" {
		return false
	}
	return true
}
