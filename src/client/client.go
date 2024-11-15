package client

import (
	"fmt"
	"github.com/ch55secake/dizzy/src/model"
	"io"
	"log"
	"net/http"
	"time"
)

// MakeRequest will return either the request if it errors or the length of the response body,
// which could indicate that there is something on the path that was just requested for
func MakeRequest(request model.Request) (error, int) {
	c := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       request.Timeout * time.Second,
	}

	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	err, length := sendRequest(request, c)
	if err != nil {
		err := fmt.Errorf("error sending request: %v", err)
		return err, 0
	}

	fmt.Printf("Body length returned: %v", length)
	return err, length
}

// sendRequest will send the request with the provided method from the request model
func sendRequest(request model.Request, client http.Client) (error, int) {
	req, err := http.NewRequest(request.Method, request.Url, nil)
	if err != nil {
		log.Fatalf("Error occured creating request: %v", err)
		return err, 0
	}

	if request.Headers != nil {
		for key, value := range request.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error occured whilst making request: %v", err), 0
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return err, 0
	}

	fmt.Printf("Length of response body: %d \n", len(body))

	return err, len(body)
}
