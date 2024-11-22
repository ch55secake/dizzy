package executor

import (
	"github.com/ch55secake/dizzy/pkg/job"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDefaultExecutor_Execute(t *testing.T) {
	t.Run("should execute without error", func(t *testing.T) {

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"message": "success"}`))
			if err != nil {
				return
			}
		}))

		defer mockServer.Close()

		mockFile := "mockfile.txt"
		var content []byte
		content = append(content, []byte("Hello")...)
		content = append(content, []byte("World")...)
		err := os.WriteFile(mockFile, content, 0644)
		if err != nil {
			t.Fatalf("failed to create mock file: %v", err)
		}
		defer os.Remove(mockFile)

		testExecutor := &DefaultExecutor{
			Dispatcher:  *job.NewDispatcher(600, 1400),
			WorkerCount: 600,
			QueueSize:   1400,
		}

		testExecutor.Execute("~/singlewordtestlist.txt", mockServer.URL)

	})
}
