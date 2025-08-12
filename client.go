// Package sarvam provides a Go client for the Sarvam AI API.
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

// Client represents a Sarvam AI API client.
type Client struct {
	baseURL string
	apiKey  string
}

// NewClient creates a new Sarvam AI client with the provided API key.
func NewClient(apiKey string) *Client {
	const baseURL = "https://api.sarvam.ai"
	return &Client{apiKey: apiKey, baseURL: baseURL}
}

// SetBaseURL allows customization of the API endpoint URL.
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// makeJsonHTTPRequest sends a JSON HTTP request to the Sarvam AI API.
func (c *Client) makeJsonHTTPRequest(method, url string, body any) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyBytes := bytes.NewBuffer(jsonBody)
	return c.makeHTTPRequest(method, url, bodyBytes, "application/json")

}

// makeHTTPRequest sends an HTTP request to the Sarvam AI API.
func (c *Client) makeHTTPRequest(method, url string, body *bytes.Buffer, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("api-subscription-key", c.apiKey)

	return http.DefaultClient.Do(req)
}

// makeMultipartRequest sends a multipart form request to the Sarvam AI API.
// TODO: This should be named better. Right now it feels too generic.
func (c *Client) makeMultipartRequest(endpoint, filePath string, model *SpeechToTextModel, languageCode *string, withTimestamps *bool) (*http.Response, error) {
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

// buildSpeechToTextTranslateRequest builds a multipart form request for speech-to-text translation.
func (c *Client) buildSpeechToTextTranslateRequest(endpoint string, params SpeechToTextTranslateParams) (*http.Response, error) {
	// Open the file
	file, err := os.Open(params.FilePath)
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
	if params.Prompt != nil {
		err = writer.WriteField("prompt", *params.Prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to write prompt field: %w", err)
		}
	}

	// Add model parameter if provided
	if params.Model != nil {
		err = writer.WriteField("model", string(*params.Model))
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

// HTTPError represents an error response from the Sarvam AI API.
type HTTPError struct {
	StatusCode int
	Message    string
	Code       string
	RequestID  string
}

// Error implements the error interface for HTTPError.
func (e *HTTPError) Error() string {
	if e.Code != "" && e.RequestID != "" {
		return fmt.Sprintf("status code: %d, code: %s, message: %s, request_id: %s", e.StatusCode, e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("status code: %d, message: %s", e.StatusCode, e.Message)
}

// parseAPIError parses an HTTP error response from the Sarvam AI API.
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

func Ptr[T any](v T) *T {
	return &v
}

// defaultClient is the default client instance used by package-level functions
var defaultClient *Client

// init initializes the default client with the API key from environment variable
func init() {
	if apiKey := os.Getenv("SARVAM_API_KEY"); apiKey != "" {
		defaultClient = NewClient(apiKey)
	}
}

// SetAPIKey sets the API key for the default client and creates a new client instance
func SetAPIKey(apiKey string) {
	defaultClient = NewClient(apiKey)
}

// GetDefaultClient returns the default client instance
func GetDefaultClient() *Client {
	return defaultClient
}

// Package-level convenience functions that use the default client

// SpeechToText is a package-level function that uses the default client
func SpeechToText(params SpeechToTextParams) (*SpeechToTextResponse, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.SpeechToText(params)
}

// SpeechToTextTranslate is a package-level function that uses the default client
func SpeechToTextTranslate(params SpeechToTextTranslateParams) (*SpeechToTextTranslateResponse, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.SpeechToTextTranslate(params)
}

// ChatCompletion is a package-level function that uses the default client
func ChatCompletion(req *ChatCompletionParams) (*ChatCompletionResponse, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.ChatCompletion(req)
}

// SimpleChatCompletion is a package-level function that uses the default client
func SimpleChatCompletion(messages []Message, model ChatCompletionModel) (*ChatCompletionResponse, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.SimpleChatCompletion(messages, model)
}

// ChatCompletionWithParams is a package-level function that uses the default client
func ChatCompletionWithParams(params *ChatCompletionParams) (*ChatCompletionResponse, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.ChatCompletionWithParams(params)
}

// Translate is a package-level function that uses the default client
func Translate(input string, sourceLanguageCode, targetLanguageCode Language) (*Translation, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.Translate(input, sourceLanguageCode, targetLanguageCode)
}

// TranslateWithParams is a package-level function that uses the default client
func TranslateWithParams(params *TranslateParams) (*Translation, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.TranslateWithParams(params)
}

// IdentifyLanguage is a package-level function that uses the default client
func IdentifyLanguage(input string) (*LanguageIdentification, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.IdentifyLanguage(input)
}

// Transliterate is a package-level function that uses the default client
func Transliterate(input string, sourceLanguage Language, targetLanguage Language) (*Transliteration, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.Transliterate(input, sourceLanguage, targetLanguage)
}

// TextToSpeech is a package-level function that uses the default client
func TextToSpeech(params TextToSpeechParams) (*TextToSpeechResponse, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient.TextToSpeech(params)
}
