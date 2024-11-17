package client

type Response struct {
	StatusCode int `json:"status_code"`
	BodyLength int `json:"body_length"`
}
