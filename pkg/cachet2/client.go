package cachet2

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

const (
	xCachetToken    = "X-Cachet-Token"
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

type Client struct {
	h     *http.Client
	url   string
	token string
}

func New(url, token string) *Client {
	return &Client{h: &http.Client{Timeout: time.Second * 10}, url: url, token: token}
}

func (c *Client) do(r *http.Request) (*http.Response, error) {
	r.Header.Add(xCachetToken, c.token)
	r.Header.Add(contentType, applicationJson)
	return c.h.Do(r)
}

func (c *Client) get(path string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.url, path), nil)
	if err != nil {
		return nil, err
	}
	return c.do(request)
}

func (c *Client) delete(path string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", c.url, path), nil)
	if err != nil {
		return nil, err
	}
	return c.do(request)
}

func (c *Client) post(path string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.url, path), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return c.do(request)
}
