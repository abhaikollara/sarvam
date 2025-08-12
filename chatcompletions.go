package sarvam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Message represents a message in the chat conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ReasoningEffort represents the reasoning effort level for chat completions.
type ReasoningEffort string

const (
	ReasoningEffortLow    ReasoningEffort = "low"
	ReasoningEffortMedium ReasoningEffort = "medium"
	ReasoningEffortHigh   ReasoningEffort = "high"
)

// ChatCompletionParams represents the parameters for chat completions.
type ChatCompletionParams struct {
	Messages         []Message           `json:"messages"`
	Model            ChatCompletionModel `json:"model"`
	Temperature      *float64            `json:"temperature,omitempty"`
	TopP             *float64            `json:"top_p,omitempty"`
	ReasoningEffort  *ReasoningEffort    `json:"reasoning_effort,omitempty"`
	MaxTokens        *int                `json:"max_tokens,omitempty"`
	Stream           *bool               `json:"stream,omitempty"`
	Stop             interface{}         `json:"stop,omitempty"` // string or []string. TODO: Find a way to make this more type safe.
	N                *int                `json:"n,omitempty"`
	Seed             *int64              `json:"seed,omitempty"`
	FrequencyPenalty *float64            `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64            `json:"presence_penalty,omitempty"`
	WikiGrounding    *bool               `json:"wiki_grounding,omitempty"`
}

// ChatCompletionChoice represents a single completion choice.
type ChatCompletionChoice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

// Usage represents token usage information for the API call.
type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionResponse represents the response from the chat completions API.
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Choices []ChatCompletionChoice `json:"choices"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Object  string                 `json:"object"`
	Usage   *Usage                 `json:"usage"`
}

// ChatCompletion creates a chat completion using the Sarvam AI API.
func (c *Client) ChatCompletion(req *ChatCompletionParams) (*ChatCompletionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("messages cannot be empty")
	}

	if req.Model == "" {
		return nil, fmt.Errorf("model is required")
	}
	// TODO: Include constraints as per the API docs

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

// SimpleChatCompletion is a convenience function for simple chat completions.
func (c *Client) SimpleChatCompletion(messages []Message, model ChatCompletionModel) (*ChatCompletionResponse, error) {
	req := &ChatCompletionParams{
		Messages: messages,
		Model:    model,
	}
	return c.ChatCompletion(req)
}

// ChatCompletionWithParams creates a chat completion with custom parameters.
func (c *Client) ChatCompletionWithParams(params *ChatCompletionParams) (*ChatCompletionResponse, error) {
	return c.ChatCompletion(params)
}

// GetFirstChoiceContent returns the content of the first choice from the response.
func (r *ChatCompletionResponse) GetFirstChoiceContent() string {
	if len(r.Choices) > 0 {
		return r.Choices[0].Message.Content
	}
	return ""
}

// GetChoiceContent returns the content of a specific choice by index.
func (r *ChatCompletionResponse) GetChoiceContent(index int) string {
	if index >= 0 && index < len(r.Choices) {
		return r.Choices[index].Message.Content
	}
	return ""
}
