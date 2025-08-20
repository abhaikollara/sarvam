// Package sarvam provides a Go client for the Sarvam AI API.
package sarvam

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
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
func (c *Client) makeJsonHTTPRequest(ctx context.Context, method, url string, body any) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	bodyBytes := bytes.NewBuffer(jsonBody)
	return c.makeHTTPRequest(ctx, method, url, bodyBytes, "application/json")
}

// makeHTTPRequest sends an HTTP request to the Sarvam AI API.
func (c *Client) makeHTTPRequest(ctx context.Context, method, url string, body *bytes.Buffer, contentType string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("api-subscription-key", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	return resp, nil
}

// buildSpeechToTextRequest builds a multipart form request for speech-to-text.
func (c *Client) buildSpeechToTextRequest(ctx context.Context, endpoint string, speech io.Reader, params SpeechToTextParams) (*http.Response, error) {

	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	// Create a form file field
	part, err := writer.CreateFormFile("file", "speech.wav")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy the file content to the form field
	_, err = io.Copy(part, speech)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add model parameter if provided
	if params.Model != nil {
		err = writer.WriteField("model", string(*params.Model))
		if err != nil {
			return nil, fmt.Errorf("failed to write model field: %w", err)
		}
	}

	// Add language_code parameter if provided
	if params.Language != nil {
		err = writer.WriteField("language_code", string(*params.Language))
		if err != nil {
			return nil, fmt.Errorf("failed to write language_code field: %w", err)
		}
	}

	// Add with_timestamps parameter if provided
	if params.WithTimestamps != nil {
		err = writer.WriteField("with_timestamps", fmt.Sprintf("%t", *params.WithTimestamps))
		if err != nil {
			return nil, fmt.Errorf("failed to write with_timestamps field: %w", err)
		}
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return c.makeHTTPRequest(ctx, http.MethodPost, c.baseURL+endpoint, &requestBody, writer.FormDataContentType())
}

// buildSpeechToTextTranslateRequest builds a multipart form request for speech-to-text translation.
func (c *Client) buildSpeechToTextTranslateRequest(ctx context.Context, endpoint string, speech io.Reader, params SpeechToTextTranslateParams) (*http.Response, error) {
	var err error

	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field
	var part io.Writer
	part, err = writer.CreateFormFile("file", "speech.wav")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, speech)
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

	// Add audio_codec parameter if provided
	if params.AudioCodec != nil {
		err = writer.WriteField("audio_codec", string(*params.AudioCodec))
		if err != nil {
			return nil, fmt.Errorf("failed to write audio_codec field: %w", err)
		}
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return c.makeHTTPRequest(ctx, http.MethodPost, c.baseURL+endpoint, &requestBody, writer.FormDataContentType())
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
var (
	defaultClient *Client
	clientMutex   sync.RWMutex
	clientOnce    sync.Once
)

// initDefaultClient initializes the default client with the API key from environment variable
func initDefaultClient() {
	if apiKey := os.Getenv("SARVAM_API_KEY"); apiKey != "" {
		defaultClient = NewClient(apiKey)
	}
}

// init initializes the default client
func init() {
	clientOnce.Do(initDefaultClient)
}

// SetAPIKey sets the API key for the default client and creates a new client instance
func SetAPIKey(apiKey string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	defaultClient = NewClient(apiKey)
}

// GetDefaultClient returns the default client instance
func GetDefaultClient() *Client {
	clientMutex.RLock()
	defer clientMutex.RUnlock()
	return defaultClient
}

// getDefaultClientSafe safely returns the default client or an error if not initialized
func getDefaultClientSafe() (*Client, error) {
	clientMutex.RLock()
	defer clientMutex.RUnlock()
	if defaultClient == nil {
		return nil, fmt.Errorf("default client not initialized. Call SetAPIKey() or set SARVAM_API_KEY environment variable")
	}
	return defaultClient, nil
}

// Package-level convenience functions that use the default client

// SpeechToText is a package-level function that uses the default client
func SpeechToText(ctx context.Context, speech io.Reader, params SpeechToTextParams) (*SpeechToTextResponse, error) {
	client, err := getDefaultClientSafe()
	if err != nil {
		return nil, err
	}
	return client.SpeechToText(ctx, speech, params)
}

// SpeechToTextTranslate is a package-level function that uses the default client
func SpeechToTextTranslate(ctx context.Context, speech io.Reader, params SpeechToTextTranslateParams) (*SpeechToTextTranslateResponse, error) {
	client, err := getDefaultClientSafe()
	if err != nil {
		return nil, err
	}
	return client.SpeechToTextTranslate(ctx, speech, params)
}

// ChatCompletion is a package-level function that uses the default client
func ChatCompletion(ctx context.Context, messages []Message, model ChatCompletionModel, req *ChatCompletionParams) (*ChatCompletionResponse, error) {
	client, err := getDefaultClientSafe()
	if err != nil {
		return nil, err
	}
	return client.ChatCompletion(ctx, messages, model, req)
}

// Translate is a package-level function that uses the default client
func Translate(ctx context.Context, input string, sourceLanguageCode, targetLanguageCode Language, params *TranslateParams) (*TranslationResponse, error) {
	client, err := getDefaultClientSafe()
	if err != nil {
		return nil, err
	}
	return client.Translate(ctx, input, sourceLanguageCode, targetLanguageCode, params)
}

// IdentifyLanguage is a package-level function that uses the default client
func IdentifyLanguage(ctx context.Context, input string) (*LanguageIdentificationResponse, error) {
	client, err := getDefaultClientSafe()
	if err != nil {
		return nil, err
	}
	return client.IdentifyLanguage(ctx, input)
}

// Transliterate is a package-level function that uses the default client
func Transliterate(ctx context.Context, input string, sourceLanguage Language, targetLanguage Language) (*TransliterationResponse, error) {
	client, err := getDefaultClientSafe()
	if err != nil {
		return nil, err
	}
	return client.Transliterate(ctx, input, sourceLanguage, targetLanguage, nil)
}

// TextToSpeech is a package-level function that uses the default client
func TextToSpeech(ctx context.Context, text string, targetLanguage Language, params TextToSpeechParams) (*TextToSpeechResponse, error) {
	client, err := getDefaultClientSafe()
	if err != nil {
		return nil, err
	}
	return client.TextToSpeech(ctx, text, targetLanguage, params)
}
