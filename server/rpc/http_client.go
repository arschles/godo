package rpc

import (
	"fmt"
	"io"
	"io/ioutil"
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
	respBody string
}

func (e errUnexpectedHTTPStatusCode) Error() string {
	return fmt.Sprintf("expected code %d from %s, got %d: %s", e.expected, e.url, e.actual, e.respBody)
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

func (h *HTTPClient) Build(ctx io.Reader, crossCompile bool, packageName string, envs []string) (io.ReadCloser, error) {
	urlStr := h.urlStr("build")
	req, err := http.NewRequest("POST", urlStr, ctx)
	if crossCompile {
		req.Header.Set(common.CrossCompileHeader, common.CrossCompileTrue)
	} else {
		req.Header.Set(common.CrossCompileHeader, common.CrossCompileFalse)
	}
	req.Header.Set(common.PackageNameHeader, packageName)
	for _, env := range envs {
		req.Header.Set(common.EnvHeader, env)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		var body string
		defer resp.Body.Close()
		if bodyBytes, err := ioutil.ReadAll(resp.Body); err == nil {
			body = string(bodyBytes)
		}
		return nil, errUnexpectedHTTPStatusCode{
			expected: http.StatusOK,
			actual:   resp.StatusCode,
			url:      urlStr,
			respBody: body,
		}
	}

	return resp.Body, nil
}
