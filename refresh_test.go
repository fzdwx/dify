package dify

import (
	"testing"
)

func TestRefreshToken(t *testing.T) {
	// Create client
	c, err := NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		t.Fatal(err)
	}

	// Cast to concrete type to access internal methods
	client := c.(*client)

	// Test refresh access token
	err = client.refreshAccessToken()
	if err != nil {
		t.Fatal("Failed to refresh access token:", err)
	}

	t.Log("Successfully refreshed access token")

	// Test refresh dataset API key
	err = client.RefreshDatasetAPIKey()
	if err != nil {
		t.Fatal("Failed to refresh dataset API key:", err)
	}

	t.Log("Successfully refreshed dataset API key")
}
