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
		t.cli.expiration.Add(time.Duration(t.cli.token.RefreshExpires-60) * time.Second)
	}
	t.cli.m.Unlock()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.cli.token.Access))

	return t.rt.RoundTrip(req)
}

func NewClient(secretId, secretKey string) (*Client, error) {
	var err error

	secretId = "3d981463-2299-4ede-913c-8608d4f312be"
	secretKey = "3e8edc8714b6efcd5f32e92fb54a0140de6bf989dbbbba26eb4c40ed019af15fa42dc7698d567acda79193acfb2dfed4fe0ace8e1283f18cf9a21ef554335ace"

	c := &Client{c: &http.Client{}, m: &sync.Mutex{}}
	c.token, err = c.newToken(secretId, secretKey)

	if err != nil {
		return nil, err
	}
	c.c.Transport = refreshTokenTransport{rt: http.DefaultTransport, cli: c}
	c.expiration = time.Now().Add(time.Duration(c.token.AccessExpires-60) * time.Second)

	return c, nil
}
