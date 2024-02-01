package prometheus

import "net/url"

type Client struct {
	baseUrl *url.URL
}

func New(baseUrl *url.URL) *Client {
	return &Client{baseUrl: baseUrl}
}

func FromUrl(baseUrl string) (*Client, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	return New(u), nil
}
