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
func (c *Client) makeMultipartRequest(endpoint, filePath string, model *Model, languageCode *string, withTimestamps *bool) (*http.Response, error) {
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

	// Add language_code parameter if provided
	if languageCode != nil {
		err = writer.WriteField("language_code", *languageCode)
		if err != nil {
			return nil, fmt.Errorf("failed to write language_code field: %w", err)
		}
	}

	// Add with_timestamps parameter if provided
	if withTimestamps != nil {
		err = writer.WriteField("with_timestamps", fmt.Sprintf("%t", *withTimestamps))
		if err != nil {
			return nil, fmt.Errorf("failed to write with_timestamps field: %w", err)
		}
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return c.makeHTTPRequest(http.MethodPost, c.baseURL+endpoint, &requestBody, writer.FormDataContentType())
}

// makeMultipartRequestTranslate makes a multipart request to the speech-to-text-translate API.
func (c *Client) makeMultipartRequestTranslate(endpoint, filePath string, prompt *string, model *Model) (*http.Response, error) {
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

	// Add prompt parameter if provided
	if prompt != nil {
		err = writer.WriteField("prompt", *prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to write prompt field: %w", err)
		}
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
	Code       string
	RequestID  string
}

func (e *HTTPError) Error() string {
	if e.Code != "" && e.RequestID != "" {
		return fmt.Sprintf("status code: %d, code: %s, message: %s, request_id: %s", e.StatusCode, e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("status code: %d, message: %s", e.StatusCode, e.Message)
}

// parseAPIError parses the error response from the API
func parseAPIError(resp *http.Response) error {
	// Try to read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// If we can't read the body, return a basic error
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Try to parse as API error format
	var apiError struct {
		Error struct {
			Message   string `json:"message"`
			Code      string `json:"code"`
			RequestID string `json:"request_id"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error.Message != "" {
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    apiError.Error.Message,
			Code:       apiError.Error.Code,
			RequestID:  apiError.Error.RequestID,
		}
	}

	// If parsing fails, return the raw body as message
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Message:    string(body),
	}
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

func String(s string) *string {
	return &s
}
