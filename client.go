package dify

type Client interface {
}

type client struct {
	baseUrl string
	key     string
}

func NewClient(baseUrl, key string) Client {
	return &client{
		baseUrl: baseUrl,
		key:     key,
	}
}
