package tfl

type opt func(*Client)

func WithBaseURL(baseURL string) opt {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithAppKey(appKey string) opt {
	return func(c *Client) {
		c.appKey = &appKey
	}
}
