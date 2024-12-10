package nordigen

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const baseUrl = "bankaccountdata.gocardless.com"
const apiPath = "/api/v2"

type Client struct {
	c         *http.Client
	secretId  string
	secretKey string

	token       *Token
	nextRefresh time.Time

	m        *sync.RWMutex
	stopChan chan struct{}
}

type Transport struct {
	rt  http.RoundTripper
	cli *Client
}

// refreshTokenIfNeeded refreshes the token if refresh time has passed
func (c *Client) refreshTokenIfNeeded(ctx context.Context) error {
	c.m.Lock()
	defer c.m.Unlock()

	if time.Now().Before(c.nextRefresh) {
		return nil
	}

	newToken, err := c.refreshToken(ctx, c.token.Refresh)
	if err != nil {
		return err
	}
	c.updateToken(newToken)
	return nil
}

// updateToken updates the client token and sets the refresh time to half the
// access token lifetime.
func (c *Client) updateToken(t *Token) {
	c.token = t
	c.nextRefresh = time.Now().Add(time.Duration(t.AccessExpires/2) * time.Second)
}

// StartTokenHandler handles token refreshes in the background
func (c *Client) StartTokenHandler(ctx context.Context) error {
	// Initialize the first token
	token, err := c.newToken(ctx)
	if err != nil {
		return errors.New("failed to get initial token: " + err.Error())
	}

	c.m.Lock()
	c.updateToken(token)
	c.m.Unlock()

	go c.tokenRefreshLoop(ctx)
	return nil
}

func (c *Client) tokenRefreshLoop(ctx context.Context) {
	for {
		c.m.RLock()
		refreshTime := c.nextRefresh
		c.m.RUnlock()

		timeToWait := time.Until(refreshTime)
		if timeToWait < 0 {
			timeToWait = 0
		}

		select {
		case <-c.stopChan:
			return
		case <-time.After(timeToWait):
			if err := c.refreshTokenIfNeeded(ctx); err != nil {
				panic(fmt.Sprintf("failed to refresh token: %s", err))
			}
		case <-ctx.Done():
			return
		}
	}
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
	c := &Client{
		c:         &http.Client{Timeout: 60 * time.Second},
		secretId:  secretId,
		secretKey: secretKey,

		m:        &sync.RWMutex{},
		stopChan: make(chan struct{}),
	}

	// Add transport to handle headers, host and path for all requests
	c.c.Transport = Transport{rt: http.DefaultTransport, cli: c}

	// Start token handler
	if err := c.StartTokenHandler(context.Background()); err != nil {
		return nil, err
	}

	return c, nil
}
