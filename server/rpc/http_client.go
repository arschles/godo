package rpc

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/arschles/gci/server/common"
)

const (
	httpScheme  = "http"
	httpsScheme = "https"
)

type errUnexpectedHTTPStatusCode struct {
	url      string
	actual   int
	expected int
}

func (e errUnexpectedHTTPStatusCode) Error() string {
	return fmt.Sprintf("expected code %d from %s, got %d", e.expected, e.url, e.actual)
}

type HTTPClient struct {
	client *http.Client
	scheme string
	host   string
	port   uint
}

func NewHTTPClient(host string, port uint) *HTTPClient {
	return &HTTPClient{client: http.DefaultClient, scheme: httpScheme, host: host, port: port}
}

func (h *HTTPClient) urlStr(path ...string) string {
	return fmt.Sprintf("%s://%s:%d/%s", h.scheme, h.host, h.port, strings.Join(path, "/"))
}

func (h *HTTPClient) Build(ctx io.Reader, crossCompile bool, envs []string) (io.ReadCloser, error) {
	urlStr := h.urlStr("build")
	req, err := http.NewRequest("POST", urlStr, ctx)
	if crossCompile {
		req.Header.Set(common.CrossCompileHeader, common.CrossCompileTrue)
	} else {
		req.Header.Set(common.CrossCompileHeader, common.CrossCompileFalse)
	}
	for _, env := range envs {
		req.Header.Set(common.EnvHeader, env)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errUnexpectedHTTPStatusCode{expected: http.StatusOK, actual: resp.StatusCode, url: urlStr}
	}

	return resp.Body, nil
}
