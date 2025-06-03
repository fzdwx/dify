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

func NewClient(baseUrl, email, password string) (Client, error) {
	c := resty.New().SetBaseURL(baseUrl)

	// Perform login to get access token
	loginReq := &LoginRequest{
		Email:      email,
		Password:   password,
		Language:   "zh-Hans",
		RememberMe: true,
	}

	var loginResp LoginResponse
	response, err := c.R().
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

	// Set the authorization header with the access token
	c.Header().Set("Authorization", "Bearer "+loginResp.Data.AccessToken)

	return &client{
		c: c,
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
	c *resty.Client
}

func (c *client) r() *resty.Request {
	return c.c.R()
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
