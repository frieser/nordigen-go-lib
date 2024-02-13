package nordigen

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const baseUrl = "ob.nordigen.com"
const apiPath = "/api/v2"

type Client struct {
	c          *http.Client
	secretId   string
	secretKey  string
	expiration time.Time
	token      *Token
	m          *sync.Mutex
}

type refreshTokenTransport struct {
	rt  http.RoundTripper
	cli *Client
}

func (t refreshTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var err error

	req.URL.Scheme = "https"
	req.URL.Host = baseUrl
	req.URL.Path = strings.Join([]string{apiPath, req.URL.Path}, "/")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	t.cli.m.Lock()

	if t.cli.expiration.Before(time.Now()) {
		t.cli.token, err = t.cli.refreshToken(t.cli.token.Refresh)

		if err != nil {
			return nil, err
		}
		t.cli.expiration = t.cli.expiration.Add(time.Duration(t.cli.token.RefreshExpires-60) * time.Second)
	}
	t.cli.m.Unlock()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.cli.token.Access))

	return t.rt.RoundTrip(req)
}

func NewClient(secretId, secretKey string) (*Client, error) {
	var err error

	c := &Client{c: &http.Client{Timeout: 60 * time.Second}, m: &sync.Mutex{}}
	c.token, err = c.newToken(secretId, secretKey)

	if err != nil {
		return nil, err
	}
	c.c.Transport = refreshTokenTransport{rt: http.DefaultTransport, cli: c}
	c.expiration = time.Now().Add(time.Duration(c.token.AccessExpires-60) * time.Second)

	return c, nil
}
