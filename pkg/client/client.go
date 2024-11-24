package client

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

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
func (r *Requester) MakeRequest(request Request) (error, Response) {
	c := http.Client{
		Timeout: r.Timeout,
	}

	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	err, bodyLength, statusCode := r.sendRequest(request, c)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Errorf("Request timed out")
			return err, Response{
				StatusCode: 408,
				BodyLength: 0,
			}
		}
		// TODO: Retry here on TCP connect/dial errors
		log.WithFields(log.Fields{
			"statusCode": statusCode,
			"bodyLength": bodyLength,
		}).Errorf("Error fetching response, statusCode: %v bodyLength: %v", statusCode, bodyLength)
		return err, Response{
			StatusCode: statusCode,
			BodyLength: bodyLength,
		}
	}

	response := Response{
		StatusCode: statusCode,
		BodyLength: bodyLength,
	}

	log.Infof("Response has body length of %v and a status of %v", response.BodyLength, response.StatusCode)
	return nil, response
}

// sendRequest will send the request with the provided method from the request model.
func (r *Requester) sendRequest(request Request, client http.Client) (error, int, int) {
	invalidError, valid := isValidHttpMethod(r)
	if valid {
		req, err := http.NewRequest(r.Method, request.ToString(), nil)
		if err != nil {
			log.WithFields(log.Fields{
				"method":  r.Method,
				"request": request.ToString(),
			}).Errorf("Error creating request.")
			return fmt.Errorf("error occurred creating request: %w", err), 0, 400
		}

		if r.Headers != nil {
			for key, value := range r.Headers {
				req.Header.Set(key, value)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return context.DeadlineExceeded, 0, 400
			}
			return fmt.Errorf("error occurred sending request: %w", err), 0, 400
		}
		defer func(Body io.ReadCloser) {
			if closeErr := Body.Close(); closeErr != nil {
				log.Warnf("Warning: error occurred closing body: %v", closeErr)
			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error occurred reading response body: %w", err), 0, 400
		}

		return nil, len(body), resp.StatusCode
	}
	return invalidError, 0, 400
}

// isValidHttpMethod will determine whether attempted http method is actually a valid operation
func isValidHttpMethod(r *Requester) (error, bool) {
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
		return fmt.Errorf("invalid HTTP method: %s", r.Method), false
	}
	return nil, true
}
