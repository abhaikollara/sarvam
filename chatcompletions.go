package sarvam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Message represents a message in the chat conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ReasoningEffort represents the reasoning effort level
type ReasoningEffort string

const (
	ReasoningEffortLow    ReasoningEffort = "low"
	ReasoningEffortMedium ReasoningEffort = "medium"
	ReasoningEffortHigh   ReasoningEffort = "high"
)

// ChatCompletionRequest represents the request payload for chat completions
type ChatCompletionRequest struct {
	Messages         []Message        `json:"messages"`
	Model            Model            `json:"model"`
	Temperature      *float64         `json:"temperature,omitempty"`
	TopP             *float64         `json:"top_p,omitempty"`
	ReasoningEffort  *ReasoningEffort `json:"reasoning_effort,omitempty"`
	MaxTokens        *int             `json:"max_tokens,omitempty"`
	Stream           *bool            `json:"stream,omitempty"`
	Stop             interface{}      `json:"stop,omitempty"` // string or []string
	N                *int             `json:"n,omitempty"`
	Seed             *int64           `json:"seed,omitempty"`
	FrequencyPenalty *float64         `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64         `json:"presence_penalty,omitempty"`
	WikiGrounding    *bool            `json:"wiki_grounding,omitempty"`
}

// ChatCompletionChoice represents a single completion choice
type ChatCompletionChoice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

// Usage represents token usage information
type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionResponse represents the response from the chat completions API
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Choices []ChatCompletionChoice `json:"choices"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Object  string                 `json:"object"`
	Usage   *Usage                 `json:"usage"`
}

// ChatCompletion creates a chat completion using the Sarvam AI API
func (c *Client) ChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("messages cannot be empty")
	}

	if req.Model == "" {
		return nil, fmt.Errorf("model is required")
	}

	resp, err := c.makeJsonHTTPRequest(http.MethodPost, c.baseURL+"/v1/chat/completions", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SimpleChatCompletion is a convenience function for simple chat completions
func (c *Client) SimpleChatCompletion(messages []Message, model Model) (*ChatCompletionResponse, error) {
	req := &ChatCompletionRequest{
		Messages: messages,
		Model:    model,
	}
	return c.ChatCompletion(req)
}

// ChatCompletionWithOptions creates a chat completion with custom options
func (c *Client) ChatCompletionWithOptions(messages []Message, model Model, options map[string]interface{}) (*ChatCompletionResponse, error) {
	req := &ChatCompletionRequest{
		Messages: messages,
		Model:    model,
	}

	// Apply options
	if temp, ok := options["temperature"].(float64); ok {
		req.Temperature = &temp
	}
	if topP, ok := options["top_p"].(float64); ok {
		req.TopP = &topP
	}
	if reasoningEffort, ok := options["reasoning_effort"].(ReasoningEffort); ok {
		req.ReasoningEffort = &reasoningEffort
	}
	if maxTokens, ok := options["max_tokens"].(int); ok {
		req.MaxTokens = &maxTokens
	}
	if stream, ok := options["stream"].(bool); ok {
		req.Stream = &stream
	}
	if stop, ok := options["stop"]; ok {
		req.Stop = stop
	}
	if n, ok := options["n"].(int); ok {
		req.N = &n
	}
	if seed, ok := options["seed"].(int64); ok {
		req.Seed = &seed
	}
	if freqPenalty, ok := options["frequency_penalty"].(float64); ok {
		req.FrequencyPenalty = &freqPenalty
	}
	if presPenalty, ok := options["presence_penalty"].(float64); ok {
		req.PresencePenalty = &presPenalty
	}
	if wikiGrounding, ok := options["wiki_grounding"].(bool); ok {
		req.WikiGrounding = &wikiGrounding
	}

	return c.ChatCompletion(req)
}

// GetFirstChoiceContent returns the content of the first choice from the response
func (r *ChatCompletionResponse) GetFirstChoiceContent() string {
	if len(r.Choices) > 0 {
		return r.Choices[0].Message.Content
	}
	return ""
}

// GetChoiceContent returns the content of a specific choice by index
func (r *ChatCompletionResponse) GetChoiceContent(index int) string {
	if index >= 0 && index < len(r.Choices) {
		return r.Choices[index].Message.Content
	}
	return ""
}
