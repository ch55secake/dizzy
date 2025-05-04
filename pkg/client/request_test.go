package client

import "testing"

func TestRequest_ToString(t *testing.T) {
	t.Run("request to string should return combined url correctly", func(t *testing.T) {
		request := Request{
			URL:       "http://example.com",
			Subdomain: "banana",
		}

		if request.ToString() != "http://example.com/banana" {
			t.Errorf("Request toString is wrong")
		}
	})
}
