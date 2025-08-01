package sarvam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SpeechToText represents the result of a speech-to-text operation.
type SpeechToText struct {
	RequestId    string `json:"request_id"`
	Transcript   string `json:"transcript"`
	LanguageCode string `json:"language_code"`
}

func (s *SpeechToText) String() string {
	return s.Transcript
}

type SpeechToTextParams struct {
	FilePath string
	Model    *Model
}

// SpeechToText converts speech from an audio file to text.
func (c *Client) SpeechToText(params SpeechToTextParams) (*SpeechToText, error) {
	resp, err := c.makeMultipartRequest("/speech-to-text", params.FilePath, params.Model)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	// Parse the response
	var response SpeechToText
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
