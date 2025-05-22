package dify

import "resty.dev/v3"

type Client interface {
	// Datasets

	// CreateEmptyDataset 创建空知识库
	CreateEmptyDataset()
}

func NewClient(baseUrl, key string) Client {
	c := resty.New().SetBaseURL(baseUrl)
	c.Header().Set("Authorization", "Bearer "+key)
	return &client{
		c: c,
	}
}

type client struct {
	c *resty.Client
}

func (c *client) r() *resty.Request {
	return c.c.R()
}
