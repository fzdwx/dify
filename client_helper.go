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

// getOrCreateDatasetAPIKeyWithRetry gets existing dataset API keys or creates a new one with retry on console token expiry
// This method uses console API (/console/api/datasets/api-keys) so it needs refresh token support
func (c *client) getOrCreateDatasetAPIKeyWithRetry() (string, error) {
	var result string
	var resultErr error

	_, err := c.executeConsoleWithRetry(func() (*resty.Response, error) {
		// First, try to get existing API keys
		var keysResp DatasetAPIKeysResponse
		resp, err := c.consoleClient.R().
			SetContentType("application/json").
			SetResult(&keysResp).
			Get("/console/api/datasets/api-keys")

		if err != nil {
			resultErr = fmt.Errorf("failed to get dataset API keys: %w", err)
			return resp, err
		}

		if resp.IsError() {
			resultErr = fmt.Errorf("failed to get dataset API keys with status %d: %s", resp.StatusCode(), resp.String())
			return resp, nil // Don't return error here, let executeWithRetry handle 401
		}

		// If we have existing keys, use the first one
		if len(keysResp.Data) > 0 {
			result = keysResp.Data[0].Token
			return resp, nil
		}

		// No existing keys, create a new one
		var newKeyResp DatasetAPIKey
		resp, err = c.consoleClient.R().
			SetContentType("application/json").
			SetResult(&newKeyResp).
			Post("/console/api/datasets/api-keys")

		if err != nil {
			resultErr = fmt.Errorf("failed to create dataset API key: %w", err)
			return resp, err
		}

		if resp.IsError() {
			resultErr = fmt.Errorf("failed to create dataset API key with status %d: %s", resp.StatusCode(), resp.String())
			return resp, nil // Don't return error here, let executeWithRetry handle 401
		}

		result = newKeyResp.Token
		return resp, nil
	})

	if err != nil {
		return "", err
	}

	if resultErr != nil {
		return "", resultErr
	}

	return result, nil
}

// refreshAccessToken refreshes the access token using the refresh token
func (c *client) refreshAccessToken() error {
	refreshReq := &RefreshTokenRequest{
		RefreshToken: c.refreshToken,
	}

	var refreshResp RefreshTokenResponse
	response, err := c.consoleClient.R().
		SetContentType("application/json").
		SetBody(refreshReq).
		SetResult(&refreshResp).
		Post("/console/api/refresh-token")

	if err != nil {
		return fmt.Errorf("refresh token request failed: %w", err)
	}

	if response.IsError() {
		return fmt.Errorf("refresh token failed with status %d: %s", response.StatusCode(), response.String())
	}

	if refreshResp.Result != "success" {
		return fmt.Errorf("refresh token failed: %s", refreshResp.Result)
	}

	// Update the console client with new access token
	c.consoleClient.Header().Set("Authorization", "Bearer "+refreshResp.Data.AccessToken)

	// Update the refresh token
	c.refreshToken = refreshResp.Data.RefreshToken

	return nil
}

// executeConsoleWithRetry executes a console API request with automatic token refresh on 401 errors
// This should only be used for /console/api/ endpoints
func (c *client) executeConsoleWithRetry(requestFunc func() (*resty.Response, error)) (*resty.Response, error) {
	// First attempt
	response, err := requestFunc()
	if err != nil {
		return response, err
	}

	// Check if we got a 401 unauthorized error (only for console API)
	if response.StatusCode() == 401 {
		// Try to refresh the access token
		if refreshErr := c.refreshAccessToken(); refreshErr != nil {
			return response, fmt.Errorf("failed to refresh console access token: %w", refreshErr)
		}

		// Retry the request with the new token
		response, err = requestFunc()
	}

	return response, err
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

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	Result string `json:"result"`
	Data   struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}
