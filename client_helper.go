package dify

import (
	"encoding/json"
	"fmt"
	"resty.dev/v3"
)

// getOrCreateDatasetAPIKey gets existing dataset API keys or creates a new one if none exist
func getOrCreateDatasetAPIKey(consoleClient *resty.Client) (string, error) {
	// First, try to get existing API keys
	var keysResp DatasetAPIKeysResponse
	response, err := consoleClient.R().
		SetContentType("application/json").
		SetResult(&keysResp).
		Get("/console/api/datasets/api-keys")

	if err != nil {
		return "", fmt.Errorf("failed to get dataset API keys: %w", err)
	}

	if response.IsError() {
		return "", fmt.Errorf("failed to get dataset API keys with status %d: %s", response.StatusCode(), response.String())
	}

	// If we have existing keys, use the first one
	if len(keysResp.Data) > 0 {
		return keysResp.Data[0].Token, nil
	}

	// No existing keys, create a new one
	var newKeyResp DatasetAPIKey
	response, err = consoleClient.R().
		SetContentType("application/json").
		SetResult(&newKeyResp).
		Post("/console/api/datasets/api-keys")

	if err != nil {
		return "", fmt.Errorf("failed to create dataset API key: %w", err)
	}

	if response.IsError() {
		return "", fmt.Errorf("failed to create dataset API key with status %d: %s", response.StatusCode(), response.String())
	}

	return newKeyResp.Token, nil
}

func buildResponse[T any](response *resty.Response, val *T) *Response[T] {
	if response.IsError() {
		var errResp struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Status  int    `json:"status"`
		}
		if err := json.Unmarshal(response.Bytes(), &errResp); err != nil {
			errResp.Code = "unknown_error"
			errResp.Message = "An unknown error occurred"
		}
		return &Response[T]{
			Response: response,
			Result:   nil,
			Code:     errResp.Code,
			Message:  errResp.Message,
			Status:   errResp.Status,
		}
	}
	resp := &Response[T]{
		Response: response,
		Result:   val,
	}
	return resp
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Language   string `json:"language"`
	RememberMe bool   `json:"remember_me"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Result string `json:"result"`
	Data   struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

// DatasetAPIKey represents a dataset API key
type DatasetAPIKey struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Token      string `json:"token"`
	LastUsedAt *int64 `json:"last_used_at"`
	CreatedAt  int64  `json:"created_at"`
}

// DatasetAPIKeysResponse represents the response when getting dataset API keys
type DatasetAPIKeysResponse struct {
	Data []DatasetAPIKey `json:"data"`
}
