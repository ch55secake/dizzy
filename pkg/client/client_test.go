package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMakeRequest(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "success"}`))
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	request := Request{
		URL:       mockServer.URL,
		Subdomain: "/banana",
	}

	r := Requester{
		Timeout: 1 * time.Second,
		Method:  "GET",
	}

	response, err := r.MakeRequest(request)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	expectedResponse := Response{
		BodyLength: 22,
		StatusCode: http.StatusOK,
	}

	if response.BodyLength != expectedResponse.BodyLength {
		t.Errorf("Expected body length %d, but got %d", expectedResponse.BodyLength, response.BodyLength)
	}
}

func TestMakeRequest_WhenBodyHasNoLength(t *testing.T) {
	t.Run("should gracefully handle a request body when it has no length", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		request := Request{
			URL: mockServer.URL,
		}

		r := Requester{
			Timeout: 1 * time.Second,
			Method:  "GET",
		}

		response, err := r.MakeRequest(request)

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		expectedResponse := Response{
			BodyLength: 0,
			StatusCode: http.StatusInternalServerError,
		}

		if response.BodyLength != expectedResponse.BodyLength {
			t.Errorf("Expected body length %d, but got %d", expectedResponse.BodyLength, response.BodyLength)
		}
	})
}

func TestMakeRequest_Timeout(t *testing.T) {
	t.Run("should return a timeout error when the request timed out", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "success"}`))
			if err != nil {
				return
			}
		}))
		defer mockServer.Close()

		request := Request{
			URL: mockServer.URL,
		}

		r := Requester{
			Timeout: 1,
			Method:  "GET",
		}

		_, err := r.MakeRequest(request)

		if err == nil {
			t.Error("Expected a timeout error, but got nil")
		}
	})
}

func TestMakeRequest_WithHeaders(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "success"}`))
		if err != nil {
			return
		}
	}))

	defer mockServer.Close()

	request := Request{
		URL: mockServer.URL,
	}

	r := Requester{
		Timeout: 5 * time.Second,
		Method:  "GET",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	response, err := r.MakeRequest(request)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	expectedResponse := Response{
		BodyLength: 22,
		StatusCode: http.StatusOK,
	}

	if response.BodyLength != expectedResponse.BodyLength {
		t.Errorf("Expected body length %d, but got %d", expectedResponse.BodyLength, response.BodyLength)
	}
}

func TestMakeRequest_Error(t *testing.T) {
	tests := []struct {
		name      string
		request   Request
		requester Requester
		wantError bool
	}{
		{
			name: "Request should return an error when the url is invalid",
			requester: Requester{
				Timeout: 5 * time.Second,
				Method:  "GET",
			},
			request: Request{
				URL: "htp://invalid-url", // Malformed URL
			},
			wantError: true,
		},
		{
			name: "Request should return an error when the http method is invalid",
			requester: Requester{
				Timeout: 5 * time.Second,
				// Unsupported HTTP method
				Method: "INVALID",
			},
			request: Request{
				URL: "http://example.com",
			},
			wantError: true,
		},
		{
			name: "Request should not error as it is a valid request, as a base test case",
			requester: Requester{
				Timeout: 5 * time.Second,
				Method:  "GET",
			},
			request: Request{
				// Valid GET request
				URL: "http://example.com",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := tt.requester.MakeRequest(tt.request)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, got statusCode: %d and err: %v", response.StatusCode, err)
				}

			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}
