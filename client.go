package sarvam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
	apiKey  string
}

func NewClient(apiKey string) *Client {
	const baseURL = "https://api.sarvam.ai"
	return &Client{apiKey: apiKey, baseURL: baseURL}
}

func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

func (c *Client) makeHTTPRequest(method, url string, body any) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("api-subscription-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("status code: %d, message: %s", e.StatusCode, e.Message)
}

func Bool(b bool) *bool {
	return &b
}

func Float64(f float64) *float64 {
	return &f
}

func Int(i int) *int {
	return &i
}
