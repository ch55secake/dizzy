package input

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ch55secake/dizzy/pkg/client"
	"log"
	"os"
	"regexp"
)

// WordList contains the current list of words in a slice of bytes and filepath
type WordList struct {
	data     [][]byte
	filepath string
}

// TransformWordListToRequests transform all list of words to a list of requests that will be handled by the executor
func (w *WordList) TransformWordListToRequests(url string) ([]client.Request, error) {
	var requests []client.Request
	for _, value := range w.data {
		requests = append(requests, client.Request{
			Url:       url,
			Subdomain: bytes.NewBuffer(value).String(),
		})
		//log.Printf("adding word, %s", bytes.NewBuffer(value).String())
	}
	return requests, nil
	//return []client.Request{client.EmptyRequest(url)}, fmt.Errorf("no data to form requests out of, returning empty request")
}

// NewWordList returns the list of data and the attached filepath
func (w *WordList) NewWordList(filepath string) error {
	readable, err := isFileReadable(filepath)
	if err != nil {
		return fmt.Errorf("failed to check if file is readable: %w", err)
	}
	if readable {
		err := w.readFile(filepath)
		if err != nil {
			return fmt.Errorf("error while creating wordlist: %w", err)
		}
	}
	return nil
}

// FilePath returns the attached filepath for the wordlist
func (w *WordList) FilePath() string {
	return w.filepath
}

// Size return the current length of the wordlist in memory
func (w *WordList) Size() int {
	return len(w.data)
}

// isFileReadable will check if the current file is readable before starting to parse the file provided
func isFileReadable(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err != nil { // this will get statistics about the provided file
		return false, err
	}
	f, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	err = f.Close()
	if err != nil {
		log.Fatalf("Error closing file: %v", err)
		return false, err
	}
	return true, nil
}

// openFile will open and return *os.File to the caller of the method, if opening the file fails will return error
// back up to the caller
func openFile(filepath string) (*os.File, error) {
	var file *os.File
	var err error
	if filepath == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(filepath)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
			return nil, err
		}
	}

	return file, nil
}

// closeFile will close the given file
func closeFile(file *os.File) {
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(file)
}

// readFile will read in the file but first check whether the file is a text file
func (w *WordList) readFile(filepath string) error {
	file, _ := openFile(filepath)

	var data [][]byte
	reader := bufio.NewScanner(file)
	re := regexp.MustCompile(`(?i)%ext%`)
	for reader.Scan() {
		if re.MatchString(reader.Text()) {
			data = append(data, []byte(reader.Text()))
		} else {
			text := reader.Text()
			data = append(data, []byte(text))
		}
	}

	closeFile(file)

	w.filepath = filepath
	w.data = data
	return reader.Err()
}
