package executor

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDefaultExecutor_Execute(t *testing.T) {
	t.Run("should execute without error", func(t *testing.T) {

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
		err := os.WriteFile(mockFile, content, 0600)
		if err != nil {
			t.Fatalf("failed to create mock file: %v", err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Fatalf("failed to remove mock file: %v", err)
			}
		}(mockFile)

		ctx := ExecutionContext{
			Filepath:       "/Users/oscar/Projects/dizzy/pkg/testdata/testlist.txt",
			URL:            mockServer.URL,
			ResponseLength: 0,
			Timeout:        0,
			Method:         "GET",
			Headers:        nil,
		}

		Execute(ctx)

	})
}
