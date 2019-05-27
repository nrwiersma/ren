package reader

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	path2 "path"
)

// HTTPReader is an http file reader.
type HTTPReader struct {
	base string
}

// NewHTTPReader returns an http file reader.
func NewHTTPReader(path string) *HTTPReader {
	return &HTTPReader{
		base: path,
	}
}

// Read reads the file at the given path.
func (r *HTTPReader) Read(path string) (string, error) {
	u, err := url.Parse(r.base)
	if err != nil {
		return "", nil
	}
	u.Path = path2.Join(u.Path, path)

	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("resolver: file not found: " + u.String())
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
