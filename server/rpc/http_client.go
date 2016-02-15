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

// Build requests to build packageName on the server. It uploads the contents of ctx (which should be tar formatted) to the server, and tells the server whether or not to cross compile and which environment variables to build with (e.g. GOOS, GOARCH, etc...). On successful return, it writes the results to target. On failure, returns an error.
//
// Writes may have occured to target even if a non-nil error was returned.
func (h *HTTPClient) Build(ctx io.Reader, target io.Writer, crossCompile bool, packageName string, envs []string) error {
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
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody string
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			errBody = string(bodyBytes)
		}
		return errUnexpectedHTTPStatusCode{
			expected: http.StatusOK,
			actual:   resp.StatusCode,
			url:      urlStr,
			respBody: errBody,
		}
	}

	if _, err := io.Copy(target, resp.Body); err != nil {
		return err
	}

	return nil
}
