package nordigen

import (
	"fmt"
	"net/http"
	"strings"
)

const baseUrl = "ob.nordigen.com"
const apiPath = "/api"

type Client struct {
	c       *http.Client
}

type addHeaderTransport struct {
	rt      http.RoundTripper
	headers map[string]string
}

func (t *addHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "https"
	req.URL.Host = baseUrl
	req.URL.Path = strings.Join([]string{apiPath, req.URL.Path}, "/")

	for k, v := range t.headers {
		req.Header.Add(k, v)
	}

	return t.rt.RoundTrip(req)
}

func newAddHeaderTransport(headers map[string]string) *addHeaderTransport {
	return &addHeaderTransport{headers: headers, rt: http.DefaultTransport}
}

func NewClient(token string) Client {
	return Client{
		c: &http.Client{
			Transport: newAddHeaderTransport(map[string]string{
				"Authorization": fmt.Sprintf("Token %s", token),
				"Content-Type": "application/json",
			}),
		},
	}
}
