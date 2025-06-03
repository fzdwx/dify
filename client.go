package dify

import (
	"context"
	"encoding/json"
	"fmt"
	"resty.dev/v3"
)

type Client interface {
	// Datasets

	// CreateEmptyDataset 创建空知识库
	CreateEmptyDataset(ctx context.Context, req *CreateEmptyDatasetRequest) (*Response[CreateEmptyDatasetResponse], error)
	// CreateByFile 通过文件创建文档
	// 此接口基于已存在知识库，在此知识库的基础上通过文件创建新的文档
	CreateByFile(ctx context.Context, req *CreateByFileRequest) (*Response[CreateByFileResponse], error)

	// RefreshDatasetAPIKey 刷新 datasets API key（当 console token 过期时可能需要）
	RefreshDatasetAPIKey() error
}

func NewClient(baseUrl, email, password string) (Client, error) {
	// Create console client for authentication and API key management
	consoleClient := resty.New().SetBaseURL(baseUrl)

	// Perform login to get access token
	loginReq := &LoginRequest{
		Email:      email,
		Password:   password,
		Language:   "zh-Hans",
		RememberMe: true,
	}

	var loginResp LoginResponse
	response, err := consoleClient.R().
		SetContentType("application/json").
		SetBody(loginReq).
		SetResult(&loginResp).
		Post("/console/api/login")

	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}

	if response.IsError() {
		return nil, fmt.Errorf("login failed with status %d: %s", response.StatusCode(), response.String())
	}

	if loginResp.Result != "success" {
		return nil, fmt.Errorf("login failed: %s", loginResp.Result)
	}

	// Set the authorization header with the access token for console operations
	consoleClient.Header().Set("Authorization", "Bearer "+loginResp.Data.AccessToken)

	// Get or create datasets API key
	datasetAPIKey, err := getOrCreateDatasetAPIKey(consoleClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get dataset API key: %w", err)
	}

	// Create datasets client with API key
	datasetsClient := resty.New().SetBaseURL(baseUrl + "/v1")
	datasetsClient.Header().Set("Authorization", "Bearer "+datasetAPIKey)

	return &client{
		datasetsClient: datasetsClient,
		consoleClient:  consoleClient,
		datasetAPIKey:  datasetAPIKey,
		refreshToken:   loginResp.Data.RefreshToken,
		baseUrl:        baseUrl,
	}, nil
}

type Response[T any] struct {
	*resty.Response
	Result  *T
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type client struct {
	datasetsClient *resty.Client
	consoleClient  *resty.Client // 用于 console API 调用
	datasetAPIKey  string        // datasets API key
	refreshToken   string        // refresh token for console API
	baseUrl        string        // base URL for API calls
}

// datasets returns a request for datasets API (/v1/datasets/)
// These APIs use datasets API key and do NOT need refresh token
func (c *client) datasets() *resty.Request {
	return c.datasetsClient.R()
}

// console returns a request for console API (/console/api/)
// These APIs use access token and NEED refresh token support
func (c *client) console() *resty.Request {
	return c.consoleClient.R()
}

// RefreshDatasetAPIKey implements the Client interface
func (c *client) RefreshDatasetAPIKey() error {
	// Get or create a new dataset API key with retry
	newAPIKey, err := c.getOrCreateDatasetAPIKeyWithRetry()
	if err != nil {
		return fmt.Errorf("failed to refresh dataset API key: %w", err)
	}

	// Update the dataset API key
	c.datasetAPIKey = newAPIKey

	// Update the datasets client with the new API key
	c.datasetsClient.Header().Set("Authorization", "Bearer "+newAPIKey)

	return nil
}

func (r *Response[T]) String() string {
	var v any
	if r.Result != nil {
		v = *r.Result
	} else {
		v = r
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "error marshaling response: " + err.Error()
	}
	return string(data)
}
