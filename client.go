package sarvam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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

func (c *Client) makeJsonHTTPRequest(method, url string, body any) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyBytes := bytes.NewBuffer(jsonBody)
	return c.makeHTTPRequest(method, url, bodyBytes, "application/json")

}

func (c *Client) makeHTTPRequest(method, url string, body *bytes.Buffer, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("api-subscription-key", c.apiKey)

	return http.DefaultClient.Do(req)
}

// makeMultipartRequest makes a multipart request to the API.
// but this either does not belong here or should be a more generic function
func (c *Client) makeMultipartRequest(endpoint, filePath string, model *Model) (*http.Response, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	// Create a form file field
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy the file content to the form field
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add model parameter if provided
	if model != nil {
		err = writer.WriteField("model", string(*model))
		if err != nil {
			return nil, fmt.Errorf("failed to write model field: %w", err)
		}
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return c.makeHTTPRequest(http.MethodPost, c.baseURL+endpoint, &requestBody, writer.FormDataContentType())
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
