package nordigen

import (
	"context"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	sharedClient *Client
	initOnce     sync.Once
)

func initTestClient(t *testing.T) *Client {
	id, idExists := os.LookupEnv("NORDIGEN_SECRET_ID")
	key, keyExists := os.LookupEnv("NORDIGEN_SECRET_KEY")
	if !idExists || !keyExists {
		t.Skip("NORDIGEN_SECRET_ID and NORDIGEN_SECRET_KEY not set")
	}

	initOnce.Do(func() {
		c := &Client{
			c:         &http.Client{Timeout: 60 * time.Second},
			secretId:  id,
			secretKey: key,

			m: &sync.RWMutex{},
		}
		c.c.Transport = Transport{rt: http.DefaultTransport, cli: c}

		// Initialize the first token
		token, err := c.newToken(context.Background())
		if err != nil {
			t.Fatalf("newToken: %s", err)
		}

		c.token = token
		sharedClient = c
	})

	return sharedClient
}

func TestAccessRefresh(t *testing.T) {
	c := initTestClient(t)

	// Expire token immediately
	c.token.AccessExpires = 0

	ctx, cancel := context.WithCancel(context.Background())
	go c.tokenHandler(ctx)
	_, err := c.ListRequisitions()
	if err != nil {
		t.Fatalf("ListRequisitions: %s", err)
	}
	cancel() // Stop handler again
}

func TestRefreshRefresh(t *testing.T) {
	c := initTestClient(t)

	// Expire token immediately
	c.token.RefreshExpires = 0

	ctx, cancel := context.WithCancel(context.Background())
	go c.tokenHandler(ctx)
	_, err := c.ListRequisitions()
	if err != nil {
		t.Fatalf("ListRequisitions: %s", err)
	}
	cancel() // Stop handler again
}
