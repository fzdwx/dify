package dify

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRefreshToken(t *testing.T) {
	// Create client
	c, err := NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		t.Fatal(err)
	}

	// Cast to concrete type to access internal methods
	client := c.(*client)

	// Test refresh access token (for console API)
	err = client.refreshAccessToken()
	if err != nil {
		t.Fatal("Failed to refresh console access token:", err)
	}

	t.Log("Successfully refreshed console access token")

	// Test refresh dataset API key (this calls console API internally)
	err = client.RefreshDatasetAPIKey()
	if err != nil {
		t.Fatal("Failed to refresh dataset API key:", err)
	}

	t.Log("Successfully refreshed dataset API key")
}

func TestConsoleVsDatasetAPIs(t *testing.T) {
	// Create client
	c, err := NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Test datasets API (uses datasets API key, no refresh token needed)
	uniqueName := fmt.Sprintf("测试API分离_%d", time.Now().Unix())
	resp, err := c.CreateEmptyDataset(ctx, &CreateEmptyDatasetRequest{
		Name:              uniqueName,
		Description:       "测试 datasets API 不需要 refresh token",
		IndexingTechnique: IndexingTechniqueEconomy,
		Permission:        DatasetPermissionAllTeamMembers,
		Provider:          DatasetProviderVendor,
		RetrievalModel: RetrievalModel{
			SearchMethod:    RetrievalModelSearchMethodHybridSearch,
			RerankingEnable: true,
			TopK:            10,
		},
	})

	if err != nil {
		t.Fatal("Datasets API failed:", err)
	}

	if !resp.IsSuccess() {
		t.Fatal("Datasets API failed:", resp.Message)
	}

	t.Log("Datasets API (uses datasets API key) works correctly")

	// Test console API indirectly (RefreshDatasetAPIKey calls console API)
	client := c.(*client)
	err = client.RefreshDatasetAPIKey()
	if err != nil {
		t.Fatal("Console API (refresh dataset API key) failed:", err)
	}

	t.Log("Console API (uses access token + refresh token) works correctly")
}
