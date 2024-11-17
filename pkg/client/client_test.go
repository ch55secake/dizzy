package client

import (
	"github.com/ch55secake/dizzy/pkg/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMakeRequest(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "success"}`))
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	request := model.Request{
		Method:    "GET",
		Url:       mockServer.URL,
		Subdomain: "/banana",
		Timeout:   5,
	}

	err, response := MakeRequest(request)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	expectedResponse := model.Response{
		BodyLength: 22,
		StatusCode: http.StatusOK,
	}

	if response.BodyLength != expectedResponse.BodyLength {
		t.Errorf("Expected body length %d, but got %d", expectedResponse.BodyLength, response.BodyLength)
	}
}

func TestMakeRequest_WhenBodyHasNoLength(t *testing.T) {
	t.Run("should gracefully handle a request body when it has no length", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		request := model.Request{
			Method:  "GET",
			Url:     mockServer.URL,
			Timeout: 5,
		}

		err, response := MakeRequest(request)

		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		expectedResponse := model.Response{
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
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}))
		defer mockServer.Close()

		request := model.Request{
			Method:  "GET",
			Url:     mockServer.URL,
			Timeout: 1,
		}

		err, _ := MakeRequest(request)

		if err == nil {
			t.Error("Expected a timeout error, but got nil")
		}
	})
}

func TestMakeRequest_WithHeaders(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "success"}`))
		if err != nil {
			return
		}
	}))

	defer mockServer.Close()

	request := model.Request{
		Method:  "GET",
		Url:     mockServer.URL,
		Timeout: 5,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	err, response := MakeRequest(request)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	expectedResponse := model.Response{
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
		request   model.Request
		wantError bool
	}{
		{
			name: "Request should return an error when the url is invalid",
			request: model.Request{
				Method: "GET",
				Url:    "htp://invalid-url", // Malformed URL
			},
			wantError: true,
		},
		{
			name: "Request should return an error when the http method is invalid",
			request: model.Request{
				Method: "INVALID", // Unsupported HTTP method
				Url:    "http://example.com",
			},
			wantError: true,
		},
		{
			name: "Request should not error as it is a valid request, as a base test case",
			request: model.Request{
				Method: "GET", // Valid GET request
				Url:    "http://example.com",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, response := MakeRequest(tt.request)

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
