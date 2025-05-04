package input

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestWordList_FilePath(t *testing.T) {
	t.Run("should return the file path", func(t *testing.T) {
		expectedPath := "/path/to/file.txt"
		w := &WordList{
			filepath: expectedPath,
		}

		actualPath := w.FilePath()

		if actualPath != expectedPath {
			t.Errorf("FilePath() = %q; want %q", actualPath, expectedPath)
		}
	})
}

func TestWordList_Size(t *testing.T) {
	t.Run("should return the file size", func(_ *testing.T) {})
	expectedSize := len(make([][]byte, 0))
	w := &WordList{
		data: make([][]byte, 0),
	}

	actualSize := w.Size()

	if actualSize != expectedSize {
		t.Errorf("Size() = %q; want %q", actualSize, expectedSize)
	}

}

func TestWordList_NewWordList(t *testing.T) {
	t.Run("should populate the word list of given mockfile", func(t *testing.T) {
		mockFile := "mockfile.txt"
		content := []byte("word1\nword2\nword3\n")
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

		wl := &WordList{}

		err = wl.NewWordList(mockFile)
		if err != nil {
			t.Errorf("NewWordList returned an unexpected error: %v", err)
		}

		expectedWords := [][]byte{
			[]byte("word1"),
			[]byte("word2"),
			[]byte("word3"),
		}

		if len(wl.data) != len(expectedWords) {
			t.Errorf("Expected %d words, got %d", len(expectedWords), len(wl.data))
		}

		for i, word := range wl.data {
			log.Printf("Word at index %d: got %q", i, string(word))
			if string(word) != string(expectedWords[i]) {
				t.Errorf("Word at index %d: got %q, want %q", i, word, expectedWords[i])
			}
		}
	})

	t.Run("should return an error if file does not exist", func(t *testing.T) {
		wl := &WordList{}
		err := wl.NewWordList("skibidi-rizz-ohio-file.txt")
		if err == nil {
			t.Errorf("Expected an error for a file that doesnt exist, but got nil")
		}
	})
}

func TestWordList_readFile(t *testing.T) {
	t.Run("should read the file content without returning an error", func(t *testing.T) {
		mockFile := "mockfile.txt"
		content := []byte("word1\nword2\nword3\n")
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

		wl := &WordList{}
		err = wl.readFile(mockFile)

		if err != nil {
			t.Errorf("readFile returned an unexpected error: %v", err)
		}
	})
}

func Test_isFileReadable(t *testing.T) {
	tests := []struct {
		name           string
		wantError      bool
		permissionCode os.FileMode
		fileName       string
		readable       bool
		shouldExist    bool
	}{
		{
			name:           "should return true if file exists and is readable",
			wantError:      false,
			permissionCode: 0644,
			fileName:       "mockfile.txt",
			readable:       true,
			shouldExist:    true,
		},
		{
			name:           "should return false if there is no permission to read the file",
			wantError:      true,
			permissionCode: 0333,
			fileName:       "mockfile.txt",
			readable:       false,
			shouldExist:    true,
		},
		{
			name:           "should return false and error if there is no permission to read the file",
			wantError:      true,
			permissionCode: 0644,
			fileName:       "i-dont-exist.txt",
			readable:       false,
			shouldExist:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldExist {
				content := []byte("word1\nword2\nword3\n")
				err := os.WriteFile(tt.fileName, content, tt.permissionCode)
				if err != nil {
					t.Fatalf("failed to create mock file: %v", err)
				}
				defer func(name string) {
					err := os.Remove(name)
					if err != nil {
						t.Fatalf("failed to remove mock file: %v", err)
					}
				}(tt.fileName)
			}

			result, err := isFileReadable(tt.fileName)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, got result: %v and err: %v", result, err)
				}

			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}

			if result != tt.readable {
				t.Errorf("isFileReadable returned %t, want %t", result, tt.readable)
			}
		})
	}
}

func Test_openFile(t *testing.T) {
	t.Run("should open the file without error", func(t *testing.T) {
		mockFile := "mockfile.txt"
		content := []byte("word1\nword2\nword3\n")
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

		file, err := openFile(mockFile)

		if err != nil {
			t.Errorf("openFile returned an unexpected error: %v", err)
		}

		if file == nil {
			t.Errorf("openFile returned nil file")
		}
	})
}

func TestWordList_TransformWordListToRequests(t *testing.T) {
	t.Run("should transform wordlist to requests", func(_ *testing.T) {
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
		content := []byte("word1\nword2\nword3\n")
		wl := &WordList{}
		err := wl.NewWordList(mockFile)
		if err != nil {
			return
		}
		err = os.WriteFile(mockFile, content, 0600)
		if err != nil {
			t.Fatalf("failed to create mock file: %v", err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Fatalf("failed to remove mock file: %v", err)
			}
		}(mockFile)

		if err != nil {
			t.Errorf("unexpected error returned when creating a new wordlist: %v", err)
		}

		requests, err := wl.TransformWordListToRequests(mockServer.URL)
		if err != nil {
			return
		}

		for _, request := range requests {
			if request.Subdomain == "" {
				t.Errorf("TransformWordListToRequests returned nil request")
			}
		}
	})
}
