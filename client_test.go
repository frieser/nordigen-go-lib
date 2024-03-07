package nordigen

import (
	"os"
	"testing"
	"time"
)

// TestClientTokenRefresh should do a successful token refresh. We force this by
// setting the expiration to a time in the past and then calling any method.
// This test will only run if you have a valid secretId and secretKey in your
// environment.
func TestClientTokenRefresh(t *testing.T) {
	id, id_exists := os.LookupEnv("NORDIGEN_SECRET_ID")
	key, key_exists := os.LookupEnv("NORDIGEN_SECRET_KEY")
	if !id_exists || !key_exists {
		t.Skip("NORDIGEN_SECRET_ID and NORDIGEN_SECRET_KEY not set")
	}

	c, err := NewClient(id, key)
	if err != nil {
		t.Fatalf("NewClient: %s", err)
	}

	c.expiration = time.Now().Add(-time.Hour)
	_, err = c.ListRequisitions()
	if err != nil {
		t.Fatalf("ListRequisitions: %s", err)
	}
}
