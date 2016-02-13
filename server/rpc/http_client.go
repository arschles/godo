package rpc

import (
	"io"
)

type HTTPClient struct {
	host string
	port uint
}

func NewHTTPClient(host string, port uint) *HTTPClient {
	return &HTTPClient{host: host, port: port}
}

func (h *HTTPClient) Build(ctx io.Reader, crossCompile bool, env []string) (io.ReadCloser, error) {
	return nil, nil
}
