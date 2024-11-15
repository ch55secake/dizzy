package client

import (
	"errors"
	"fmt"
	"github.com/ch55secake/dizzy/src/model"
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

	err, bodyLength := MakeRequest(request)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	expectedLength := 22
	if bodyLength != expectedLength {
		t.Errorf("Expected body length %d, but got %d", expectedLength, bodyLength)
	}
}

func TestMakeRequest_WhenBodyHasNoLength(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	request := model.Request{
		Method:  "GET",
		Url:     mockServer.URL,
		Timeout: 5,
	}

	err, bodyLength := MakeRequest(request)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	expectedLength := 0
	if bodyLength != expectedLength {
		t.Errorf("Expected body length %d, but got %d", expectedLength, bodyLength)
	}
}

func TestMakeRequest_Timeout(t *testing.T) {
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

	err, bodyLength := MakeRequest(request)

	if err == nil {
		t.Error("Expected a timeout error, but got nil")
	}

	if bodyLength != 0 {
		t.Errorf("Expected body length 0, but got %d", bodyLength)
	}
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

	err, bodyLength := MakeRequest(request)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	expectedLength := 22
	if bodyLength != expectedLength {
		t.Errorf("Expected body length %d, but got %d", expectedLength, bodyLength)
	}
}

func TestMakeRequest_Error(t *testing.T) {
	tests := []struct {
		name      string
		request   model.Request
		wantError bool
	}{
		{
			name: "Invalid URL",
			request: model.Request{
				Method: "GET",
				Url:    "htp://invalid-url",
			},
			wantError: true,
		},
		{
			name: "Invalid HTTP Method",
			request: model.Request{
				Method: "INVALID",
				Url:    "http://example.com",
			},
			wantError: true,
		},
		{
			name: "Valid Request",
			request: model.Request{
				Method: "GET",
				Url:    "http://example.com",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.wantError {
				err, _ = MakeRequest(tt.request)
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else {
					if !errors.Is(err, fmt.Errorf("error sending request: %v", err)) {
						t.Errorf("Expected error message but got: %v", err)
					}
				}
			} else {
				err, _ = MakeRequest(tt.request)
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}
