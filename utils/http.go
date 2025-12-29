package utils

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

var (
	httpClient = &fasthttp.Client{}
)

type Client struct {
	BaseURL string
	APIKey  string
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

func (c *Client) Get(endpoint string) ([]byte, int, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	url := c.BaseURL + endpoint

	req.Header.SetMethod("GET")
	req.SetRequestURI(url)
	req.Header.Set("x-api-key", c.APIKey)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := httpClient.Do(req, resp)
	if err != nil {
		return nil, 0, err
	}

	bodyCopy := make([]byte, len(resp.Body()))
	copy(bodyCopy, resp.Body())

	return bodyCopy, resp.StatusCode(), nil
}

func (c *Client) Post(endpoint string, body any) ([]byte, int, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	url := c.BaseURL + endpoint

	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	req.SetBody(jsonBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = httpClient.Do(req, resp)
	if err != nil {
		return nil, 0, err
	}

	bodyCopy := make([]byte, len(resp.Body()))
	copy(bodyCopy, resp.Body())

	return bodyCopy, resp.StatusCode(), nil
}
