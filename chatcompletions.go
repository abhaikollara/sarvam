package sarvam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Message represents a message in the chat conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MessageRole string

const (
	MessageRoleSystem    MessageRole = "system"
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

func NewMessage(role MessageRole, content string) Message {
	return Message{
		Role:    string(role),
		Content: content,
	}
}

func NewSystemMessage(content string) Message {
	return NewMessage(MessageRoleSystem, content)
}

func NewUserMessage(content string) Message {
	return NewMessage(MessageRoleUser, content)
}

func NewAssistantMessage(content string) Message {
	return NewMessage(MessageRoleAssistant, content)
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
	Temperature      *float64
	TopP             *float64
	ReasoningEffort  *ReasoningEffort
	MaxTokens        *int
	Stream           *bool
	Stop             []string // string or []string. TODO: Find a way to make this more type safe.
	N                *int
	Seed             *int64
	FrequencyPenalty *float64
	PresencePenalty  *float64
	WikiGrounding    *bool
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
func (c *Client) ChatCompletion(ctx context.Context, messages []Message, model ChatCompletionModel, req *ChatCompletionParams) (*ChatCompletionResponse, error) {
	// Input validation
	if err := validateChatCompletionInput(messages, model, req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	type chatCompletionRequest struct {
		Model            ChatCompletionModel `json:"model"`
		Messages         []Message           `json:"messages"`
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

	var payload chatCompletionRequest
	payload.Model = model
	payload.Messages = messages

	if req != nil {
		if req.Temperature != nil {
			payload.Temperature = req.Temperature
		}
		if req.TopP != nil {
			payload.TopP = req.TopP
		}
		if req.ReasoningEffort != nil {
			payload.ReasoningEffort = req.ReasoningEffort
		}
		if req.MaxTokens != nil {
			payload.MaxTokens = req.MaxTokens
		}
		if req.Stream != nil {
			payload.Stream = req.Stream
		}
		if len(req.Stop) > 0 {
			if len(req.Stop) == 1 {
				payload.Stop = req.Stop[0]
			} else {
				payload.Stop = req.Stop
			}
		}
		if req.N != nil {
			payload.N = req.N
		}
		if req.Seed != nil {
			payload.Seed = req.Seed
		}
		if req.FrequencyPenalty != nil {
			payload.FrequencyPenalty = req.FrequencyPenalty
		}
		if req.PresencePenalty != nil {
			payload.PresencePenalty = req.PresencePenalty
		}
		if req.WikiGrounding != nil {
			payload.WikiGrounding = req.WikiGrounding
		}
	}

	resp, err := c.makeJsonHTTPRequest(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", payload)
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

// validateChatCompletionInput validates inputs for chat completion requests
func validateChatCompletionInput(messages []Message, model ChatCompletionModel, params *ChatCompletionParams) error {
	if len(messages) == 0 {
		return fmt.Errorf("messages cannot be empty")
	}

	if model == "" {
		return fmt.Errorf("model is required")
	}

	// Validate messages
	for i, msg := range messages {
		if strings.TrimSpace(msg.Role) == "" {
			return fmt.Errorf("message %d: role cannot be empty", i)
		}
		if strings.TrimSpace(msg.Content) == "" {
			return fmt.Errorf("message %d: content cannot be empty", i)
		}

		// Validate role values
		switch msg.Role {
		case string(MessageRoleSystem), string(MessageRoleUser), string(MessageRoleAssistant):
			// Valid roles
		default:
			return fmt.Errorf("message %d: invalid role '%s'", i, msg.Role)
		}
	}

	// Validate optional parameters
	if params != nil {
		if params.Temperature != nil && (*params.Temperature < 0.0 || *params.Temperature > 2.0) {
			return fmt.Errorf("temperature must be between 0.0 and 2.0")
		}
		if params.TopP != nil && (*params.TopP < 0.0 || *params.TopP > 1.0) {
			return fmt.Errorf("top_p must be between 0.0 and 1.0")
		}
		if params.MaxTokens != nil && *params.MaxTokens <= 0 {
			return fmt.Errorf("max_tokens must be greater than 0")
		}
		if params.N != nil && *params.N <= 0 {
			return fmt.Errorf("n must be greater than 0")
		}
		if params.FrequencyPenalty != nil && (*params.FrequencyPenalty < -2.0 || *params.FrequencyPenalty > 2.0) {
			return fmt.Errorf("frequency_penalty must be between -2.0 and 2.0")
		}
		if params.PresencePenalty != nil && (*params.PresencePenalty < -2.0 || *params.PresencePenalty > 2.0) {
			return fmt.Errorf("presence_penalty must be between -2.0 and 2.0")
		}
	}

	return nil
}
