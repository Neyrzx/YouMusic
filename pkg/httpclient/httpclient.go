package httpclient

import (
	"io"
	"strings"
)

type HTTPClient struct {
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

func (c *HTTPClient) Get(_ string) (io.ReadCloser, error) {
	// res, err := http.Get(path)
	// if err != nil {
	// 	return nil, err
	// }

	// return res.Body, nil
	return io.NopCloser(strings.NewReader(`{"releaseDate":"01.02.2000","link":"link","text":"text"}`)), nil
}
