package jira

import (
	"io"
	"net/http"
	"strings"
)

type Client struct {
	baseURL    string
	email      string
	apiToken   string
	httpClient *http.Client
}

func NewClient(baseURL, email, apiToken string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		email:      email,
		apiToken:   apiToken,
		httpClient: &http.Client{},
	}
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	// Jira Cloud API token auth requires Basic auth: base64(email:token)
	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}
