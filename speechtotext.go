package sarvam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Timestamps represents word-level timing information
type Timestamps struct {
	Words            []string  `json:"words"`
	StartTimeSeconds []float64 `json:"start_time_seconds"`
	EndTimeSeconds   []float64 `json:"end_time_seconds"`
}

// DiarizedEntry represents a single speaker's transcript segment
type DiarizedEntry struct {
	Transcript       string  `json:"transcript"`
	StartTimeSeconds float64 `json:"start_time_seconds"`
	EndTimeSeconds   float64 `json:"end_time_seconds"`
	SpeakerID        string  `json:"speaker_id"`
}

// DiarizedTranscript represents the complete diarized transcript
type DiarizedTranscript struct {
	Entries []DiarizedEntry `json:"entries"`
}

// SpeechToText represents the result of a speech-to-text operation.
type SpeechToText struct {
	RequestId          string              `json:"request_id"`
	Transcript         string              `json:"transcript"`
	Timestamps         *Timestamps         `json:"timestamps,omitempty"`
	DiarizedTranscript *DiarizedTranscript `json:"diarized_transcript,omitempty"`
	LanguageCode       string              `json:"language_code"`
}

func (s *SpeechToText) String() string {
	return s.Transcript
}

// SpeechToTextParams contains parameters for speech-to-text conversion
type SpeechToTextParams struct {
	FilePath       string  // Required: Path to the audio file
	Model          *Model  // Optional: Model to use (default: saarika:v2.5)
	LanguageCode   *string // Optional: Language code for the input audio
	WithTimestamps *bool   // Optional: Whether to include timestamps in response
}

// SpeechToText converts speech from an audio file to text.
func (c *Client) SpeechToText(params SpeechToTextParams) (*SpeechToText, error) {
	resp, err := c.makeMultipartRequest("/speech-to-text", params.FilePath, params.Model, params.LanguageCode, params.WithTimestamps)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	// Parse the response
	var response SpeechToText
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
