package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Requester struct {
	Timeout time.Duration `json:"timeout"`
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
			return fmt.Errorf("request timed out: %w", context.DeadlineExceeded), Response{
				StatusCode: 408, // 408 indicates timeout
				BodyLength: 0,
			}
		}
		return fmt.Errorf("error fetching response, statusCode: %v bodyLength: %v, err: %v", statusCode, bodyLength, err), Response{
			StatusCode: statusCode,
			BodyLength: bodyLength,
		}
	}

	response := Response{
		StatusCode: statusCode,
		BodyLength: bodyLength,
	}

	log.Printf("Response body length: %d \n", response.BodyLength)
	log.Printf("Status code: %d \n", response.StatusCode)

	return nil, response
}

// sendRequest will send the request with the provided method from the request model.
func (r *Requester) sendRequest(request Request, client http.Client) (error, int, int) {
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

	if !validMethods[request.Method] {
		return fmt.Errorf("invalid HTTP method: %s", request.Method), 0, 405
	}

	req, err := http.NewRequest(request.Method, request.Url, nil)
	if err != nil {
		return fmt.Errorf("error occurred creating request: %w", err), 0, 400
	}

	if request.Headers != nil {
		for key, value := range request.Headers {
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
			log.Printf("Warning: error occurred closing body: %v", closeErr)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error occurred reading response body: %w", err), 0, 400
	}

	return nil, len(body), resp.StatusCode
}
