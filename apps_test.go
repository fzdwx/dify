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

func TestCreateAppAccessToken(t *testing.T) {
	// Create client
	c, err := NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// First create a chat app
	appName := fmt.Sprintf("测试应用_for_token_%d", time.Now().Unix())
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

	// Create an access token for the app
	tokenResp, err := c.CreateAppAccessToken(ctx, &CreateAppAccessTokenRequest{
		AppID: appResp.Result.ID,
	})

	if err != nil {
		t.Fatal("Failed to create app access token:", err)
	}

	if !tokenResp.IsSuccess() {
		t.Fatal("Create app access token failed:", tokenResp.Message)
	}

	t.Logf("Successfully created app access token: %s (ID: %s, Type: %s)",
		tokenResp.Result.Token, tokenResp.Result.ID, tokenResp.Result.Type)

	// Verify the token format
	if tokenResp.Result.Type != "app" {
		t.Errorf("Expected token type 'app', got '%s'", tokenResp.Result.Type)
	}

	if len(tokenResp.Result.Token) == 0 {
		t.Error("Token should not be empty")
	}

	if tokenResp.Result.CreatedAt == 0 {
		t.Error("CreatedAt should not be zero")
	}
}

func TestCallWorkflowAppBlocking(t *testing.T) {
	// Create client
	c, err := NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	resp, err := c.CallWorkflowAppBlocking(ctx, &CallWorkflowRequest{
		Inputs: map[string]interface{}{
			"question":           "设备运行状态分类汇总统计",
			"user_id":            0,
			"knowledge_base_ids": "9",
			"internet_search":    1,
			"thinking":           1,
		},
		ResponseMode: ResponseModeBlocking,
		User:         "test_user",
		Token:        "app-SwD6MzpqOwiOB2tEyS074VIi",
	})
	if err != nil {
		t.Fatal("Failed to call workflow app blocking:", err)
	}
	if !resp.IsSuccess() {
		t.Fatal("Call workflow app blocking failed:", resp.Message)
	}
	t.Logf("Successfully called workflow app blocking: %s", resp.Result.Data.Outputs["text"])
}

func TestCallWorkflowAppStreaming(t *testing.T) {
	// Create client
	c, err := NewClient("http://192.168.50.21:88", "likelovec@gmail.com", "Pwd123456")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	resp, err := c.CallWorkflowAppStreaming(ctx, &CallWorkflowRequest{
		Inputs: map[string]interface{}{
			"question":           "设备运行状态分类汇总统计",
			"user_id":            0,
			"knowledge_base_ids": "9",
			"internet_search":    1,
			"thinking":           1,
		},
		ResponseMode: ResponseModeStreaming,
		User:         "test_user",
		Token:        "app-SwD6MzpqOwiOB2tEyS074VIi",
	})
	if err != nil {
		t.Fatal("Failed to call workflow app blocking:", err)
	}

	for chunk := range resp {
		if chunk == nil {
			t.Fatal("Received nil chunk from streaming response")
		}
		t.Logf("Received chunk: %s", chunk.Data.Text)
	}
}
