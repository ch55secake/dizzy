package job

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/ch55secake/dizzy/pkg/client"
)

func TestWorker_Start(t *testing.T) {
	t.Run("should be able to start a worker without errors", func(_ *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "success"}`))
			if err != nil {
				return
			}
		}))
		defer mockServer.Close()

		jobChannel := make(chan *Job)
		wg := &sync.WaitGroup{}

		mockRequester := &client.Requester{
			Timeout: 10 * time.Second,
			Method:  "GET",
		}

		worker := &Worker{
			ID:         1,
			JobChannel: jobChannel,
			Requester:  mockRequester,
			wg:         wg,
		}

		worker.Start()

		mockRequest := client.Request{
			URL: mockServer.URL,
		}

		job := NewJob(1, mockRequest)
		wg.Add(1)
		jobChannel <- job

		wg.Wait()

		close(jobChannel)
	})
}
