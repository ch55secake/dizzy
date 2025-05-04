package client

// Response that the client will map too, this tool only cares about statusCode and bodyLength so that is all that is
// mapped
type Response struct {
	StatusCode int    `json:"status_code"`
	BodyLength int    `json:"body_length"`
	Subdomain  string `json:"subdomain"`
}
