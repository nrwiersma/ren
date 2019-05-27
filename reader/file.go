package reader

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileReader is a local file reader.
type FileReader struct {
	base string
}

// NewFileReader returns a file reader.
func NewFileReader(path string) *FileReader {
	if strings.HasPrefix(path, "/.") {
		path = path[1:]
	}

	return &FileReader{
		base: path,
	}
}

// Read reads the file at the given path.
func (r *FileReader) Read(path string) (string, error) {
	path = filepath.Join(r.base, path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
