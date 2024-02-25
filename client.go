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
	stopChan   chan struct{}
}

type Transport struct {
	rt  http.RoundTripper
	cli *Client
}

func (c *Client) refreshTokenIfNeeded() error {
	c.m.Lock()
	defer c.m.Unlock()

	if time.Now().Add(time.Minute).Before(c.expiration) {
		return nil
	} else {
		// Refresh the token if its expiration is less than a minute away
		newToken, err := c.refreshToken(c.token.Refresh)
		if err != nil {
			return err
		}
		c.token = newToken
		c.expiration = time.Now().Add(time.Duration(newToken.RefreshExpires-60) * time.Second)
		return nil
	}
}

func (c *Client) StartTokenHandler() {
	c.stopChan = make(chan struct{})

	// Initialize the first token and start the token handler
	token, err := c.newToken()
	if err != nil {
		panic("Failed to get initial token: " + err.Error())
	}
	c.token = token

	go func() {
		for {
			timeToWait := time.Until(c.expiration) - time.Minute
			if timeToWait < 0 {
				// If the token is already expired, try to refresh immediately
				timeToWait = 0
			}

			select {
			case <-c.stopChan:
				return
			case <-time.After(timeToWait):
				if err := c.refreshTokenIfNeeded(); err != nil {
					// TODO(Martin): add retry logic
					panic("Failed to refresh token: " + err.Error())
				}
			}
		}
	}()
}

func (c *Client) StopTokenHandler() {
	close(c.stopChan)
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "https"
	req.URL.Host = baseUrl
	req.URL.Path = strings.Join([]string{apiPath, req.URL.Path}, "/")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Add the access token to the request if it exists
	if t.cli.token != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.cli.token.Access))
	}

	return t.rt.RoundTrip(req)
}

// NewClient creates a new Nordigen client that handles token refreshes and adds
// the necessary headers, host, and path to all requests.
func NewClient(secretId, secretKey string) (*Client, error) {
	c := &Client{c: &http.Client{Timeout: 60 * time.Second}, m: &sync.Mutex{},
		secretId:  secretId,
		secretKey: secretKey,
	}

	// Add transport to handle headers, host and path for all requests
	c.c.Transport = Transport{rt: http.DefaultTransport, cli: c}

	// Start token handler
	c.StartTokenHandler()
	defer c.StopTokenHandler()

	return c, nil
}
