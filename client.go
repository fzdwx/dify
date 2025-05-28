package dify

import (
	"context"
	"encoding/json"
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

func NewClient(baseUrl, key string) Client {
	c := resty.New().SetBaseURL(baseUrl)
	c.Header().Set("Authorization", "Bearer "+key)
	return &client{
		c: c,
	}
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
