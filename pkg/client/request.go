package client

import log "github.com/sirupsen/logrus"

// Request structure will be used to send requests and later on as flags as part the command
type Request struct {
	URL       string `json:"url"`
	Subdomain string `json:"subdomain"`
}

// EmptyRequest used for when wordlist has no data should be attached to an error
func EmptyRequest(url string) Request {
	return Request{
		URL: url,
	}
}

// isValid is used to determine whether a request can actually reach a valid subdomain/endpoint
func (req Request) isValid() bool {
	return req.Subdomain != ""
}

// ToString will combine the given subdomain with the url unless the subdomain is nil
func (req Request) ToString() string {
	if req.isValid() {
		log.Debugf("Concatenated url request will be made with: %v", req.URL+"/"+req.Subdomain)
		return req.URL + "/" + req.Subdomain
	}
	return req.URL
}
