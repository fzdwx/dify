package dify

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCreateChatApp(t *testing.T) {
	// Create client
	c, err := NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Create a chat app with unique name
	uniqueName := fmt.Sprintf("测试聊天应用_%d", time.Now().Unix())
	resp, err := c.CreateChatApp(ctx, &CreateChatAppRequest{
		Name: uniqueName,
	})

	if err != nil {
		t.Fatal("Failed to create chat app:", err)
	}

	if !resp.IsSuccess() {
		t.Fatal("Create chat app failed:", resp.Message)
	}

	t.Logf("Successfully created chat app: %s (ID: %s)", resp.Result.Name, resp.Result.ID)

	// Test updating the app model config
	updateResp, err := c.UpdateAppModelConfig(ctx, &UpdateAppModelConfigRequest{
		AppID: resp.Result.ID,
		Model: ModelConfig{
			Provider:         "langgenius/deepseek/deepseek",
			Name:             "deepseek-chat",
			Mode:             "chat",
			CompletionParams: map[string]interface{}{},
		},
		DatasetID: "", // No dataset for this test
	})

	if err != nil {
		t.Fatal("Failed to update app model config:", err)
	}

	if !updateResp.IsSuccess() {
		t.Fatal("Update app model config failed:", updateResp.Message)
	}

	t.Log("Successfully updated app model config")
}

func TestCreateChatAppWithDataset(t *testing.T) {
	// Create client
	c, err := NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// First create a dataset
	datasetName := fmt.Sprintf("测试数据集_for_app_%d", time.Now().Unix())
	datasetResp, err := c.CreateEmptyDataset(ctx, &CreateEmptyDatasetRequest{
		Name:              datasetName,
		Description:       "用于测试应用绑定的数据集",
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
		t.Fatal("Failed to create dataset:", err)
	}

	if !datasetResp.IsSuccess() {
		t.Fatal("Create dataset failed:", datasetResp.Message)
	}

	t.Logf("Successfully created dataset: %s (ID: %s)", datasetResp.Result.Name, datasetResp.Result.ID)

	// Create a chat app
	appName := fmt.Sprintf("测试应用_with_dataset_%d", time.Now().Unix())
	appResp, err := c.CreateChatApp(ctx, &CreateChatAppRequest{
		Name: appName,
	})

	if err != nil {
		t.Fatal("Failed to create chat app:", err)
	}

	if !appResp.IsSuccess() {
		t.Fatal("Create chat app failed:", appResp.Message)
	}

	t.Logf("Successfully created chat app: %s (ID: %s)", appResp.Result.Name, appResp.Result.ID)

	// Update the app to bind with the dataset
	updateResp, err := c.UpdateAppModelConfig(ctx, &UpdateAppModelConfigRequest{
		AppID: appResp.Result.ID,
		Model: ModelConfig{
			Provider:         "langgenius/deepseek/deepseek",
			Name:             "deepseek-chat",
			Mode:             "chat",
			CompletionParams: map[string]interface{}{},
		},
		DatasetID: datasetResp.Result.ID, // Bind with the dataset
	})

	if err != nil {
		t.Fatal("Failed to update app model config with dataset:", err)
	}

	if !updateResp.IsSuccess() {
		t.Fatal("Update app model config with dataset failed:", updateResp.Message)
	}

	t.Log("Successfully updated app model config with dataset binding")
}
