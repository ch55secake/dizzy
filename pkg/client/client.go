// Package client provides a client to make http requests with
package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Requester is the default implementation for the http client
type Requester struct {
	Timeout time.Duration     `json:"timeout"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"header"`
}

// NewRequester will create a new requester object that will allow you to set a timeout
func NewRequester(timeout time.Duration, method string, headers map[string]string) *Requester {
	if timeout == 0 || method == "" {
		log.Infof("Cannot have a timeout of zero, will default to a timeout of ten seconds")
		return &Requester{
			Timeout: 10 * time.Second,
			Method:  "GET",
			Headers: headers,
		}
	}
	return &Requester{
		Timeout: timeout,
		Method:  method,
		Headers: headers,
	}
}

// MakeRequest will return either the error if it occurs or the length of the response body,
// which could indicate that there is something on the path that was just requested for.
func (r *Requester) MakeRequest(request Request) (Response, error) {
	c := http.Client{
		Timeout: r.Timeout,
	}

	c.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}

	bodyLength, statusCode, err := r.sendRequest(request, c)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Errorf("Request timed out")
			return Response{
				StatusCode: 408,
				BodyLength: 0,
			}, err
		}
		// TODO: Retry here on TCP connect/dial errors
		log.WithFields(log.Fields{
			"statusCode": statusCode,
			"bodyLength": bodyLength,
		}).Errorf("Error fetching response, statusCode: %v bodyLength: %v", statusCode, bodyLength)
		return Response{
			StatusCode: statusCode,
			BodyLength: bodyLength,
		}, err
	}

	response := Response{
		StatusCode: statusCode,
		BodyLength: bodyLength,
	}

	log.Infof("Response has body length of %v and a status of %v", response.BodyLength, response.StatusCode)
	return response, nil
}

// sendRequest will send the request with the provided method from the request model.
func (r *Requester) sendRequest(request Request, client http.Client) (int, int, error) {
	valid, invalidError := isValidHTTPMethod(r)
	if valid {
		req, err := http.NewRequest(r.Method, request.ToString(), nil)
		if err != nil {
			log.WithFields(log.Fields{
				"method":  r.Method,
				"request": request.ToString(),
			}).Errorf("Error creating request.")
			return 0, 400, fmt.Errorf("error occurred creating request: %w", err)
		}

		if r.Headers != nil {
			for key, value := range r.Headers {
				req.Header.Set(key, value)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return 0, 400, context.DeadlineExceeded
			}
			return 0, 400, fmt.Errorf("error occurred sending request: %w", err)
		}
		defer func(Body io.ReadCloser) {
			if closeErr := Body.Close(); closeErr != nil {
				log.Warnf("Warning: error occurred closing body: %v", closeErr)
			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, 400, fmt.Errorf("error occurred reading response body: %w", err)
		}

		return len(body), resp.StatusCode, nil
	}
	return 0, 400, invalidError
}

// isValidHTTPMethod will determine whether attempted http method is actually a valid operation
func isValidHTTPMethod(r *Requester) (bool, error) {
	validMethods := map[string]bool{
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodPatch:   true,
		http.MethodDelete:  true,
		http.MethodConnect: true,
		http.MethodOptions: true,
		http.MethodTrace:   true,
	}

	if !validMethods[r.Method] {
		return false, fmt.Errorf("invalid HTTP method: %s", r.Method)
	}
	return true, nil
}
