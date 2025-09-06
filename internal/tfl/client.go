package tfl

import (
	"fmt"
	"net/http"
	"strings"
)

const defaultBaseURL = "https://api.tfl.gov.uk"

type Client struct {
	*http.Client

	baseURL string
	appKey  *string
}

func NewClient(opts ...opt) *Client {
	c := &Client{
		Client:  http.DefaultClient,
		baseURL: defaultBaseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) get(route string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", strings.TrimRight(c.baseURL, "/"), strings.TrimLeft(route, "/"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if c.appKey != nil {
		q := req.URL.Query()
		q.Add("app_key", *c.appKey)
		req.URL.RawQuery = q.Encode()
	}

	// need to set this header to avoid cloudflare blocks
	req.Header.Set("User-Agent", "https://github.com/prnvbn/bq")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 status code: %s", resp.Status)
	}

	return resp, nil
}
