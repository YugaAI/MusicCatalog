package httpclient

import "net/http"

//go:generate mockgen -source=client.go -destination=client_mock.go -package=httpclient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	client HTTPClient
}

func NewClent(client HTTPClient) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
