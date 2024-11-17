package input

import (
	"log"
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
	t.Run("should return the file size", func(t *testing.T) {})
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
		err := os.WriteFile(mockFile, content, 0644)
		if err != nil {
			t.Fatalf("failed to create mock file: %v", err)
		}
		defer os.Remove(mockFile)

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

}

func Test_isFileReadable(t *testing.T) {
}

func Test_openFile(t *testing.T) {
}
