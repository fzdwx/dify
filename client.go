package dify

import (
	"context"
	"resty.dev/v3"
)

type Client interface {
	// Datasets

	// CreateEmptyDataset 创建空知识库
	CreateEmptyDataset(ctx context.Context, req *CreateEmptyDatasetRequest) (*Response[CreateEmptyDatasetResponse], error)
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
	Result *T
}

type client struct {
	c *resty.Client
}

func (c *client) r() *resty.Request {
	return c.c.R()
}

func buildResponse[T any](response *resty.Response, val *T) *Response[T] {
	resp := &Response[T]{
		Response: response,
		Result:   val,
	}
	return resp
}
