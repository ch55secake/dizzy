package client

import (
	"net/http"
	"net/http/httptest"
	"reflect"
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

	request := Request{
		Url:       mockServer.URL,
		Subdomain: "/banana",
	}

	r := Requester{
		Timeout: 1 * time.Second,
		Method:  "GET",
	}

	err, response := r.MakeRequest(request)

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
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		request := Request{
			Url: mockServer.URL,
		}

		r := Requester{
			Timeout: 1 * time.Second,
			Method:  "GET",
		}

		err, response := r.MakeRequest(request)

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
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}))
		defer mockServer.Close()

		request := Request{
			Url: mockServer.URL,
		}

		r := Requester{
			Timeout: 1,
			Method:  "GET",
		}

		err, _ := r.MakeRequest(request)

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

	request := Request{
		Url: mockServer.URL,
	}

	r := Requester{
		Timeout: 5 * time.Second,
		Method:  "GET",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	err, response := r.MakeRequest(request)
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
				Url: "htp://invalid-url", // Malformed URL
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
				Url: "http://example.com",
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
				Url: "http://example.com",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, response := tt.requester.MakeRequest(tt.request)

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

func Test_isValidHttpMethod(t *testing.T) {
	type args struct {
		r *Requester
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := isValidHttpMethod(tt.args.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("isValidHttpMethod() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("isValidHttpMethod() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
