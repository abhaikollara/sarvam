package sarvam

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

// TextToSpeechResponse represents the result of a text-to-speech operation.
type TextToSpeechResponse struct {
	RequestId string
	Audios    []string
}

func (t *TextToSpeechResponse) Bytes() ([]byte, error) {
	return convertAndConcatBase64ToBytes(t.Audios)
}

// Save saves the text-to-speech data as a WAV file.
func (t *TextToSpeechResponse) Save(filename string) error {
	data, err := t.Bytes()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// TextToSpeechParams contains all parameters for text-to-speech conversion.
type TextToSpeechParams struct {
	Speaker             *Speaker
	Pitch               *float64
	Pace                *float64
	Loudness            *float64
	SpeechSampleRate    *SpeechSampleRate
	EnablePreprocessing *bool
	Model               *TextToSpeechModel
}

// SpeechSampleRate represents the audio sample rate for text-to-speech output.
type SpeechSampleRate int

var (
	SpeechSampleRate8000  SpeechSampleRate = 8000
	SpeechSampleRate16000 SpeechSampleRate = 16000
	SpeechSampleRate22050 SpeechSampleRate = 22050
	SpeechSampleRate24000 SpeechSampleRate = 24000
)

// TextToSpeech converts text to speech in the specified language.
func (c *Client) TextToSpeech(ctx context.Context, text string, targetLanguage Language, params TextToSpeechParams) (*TextToSpeechResponse, error) {
	// Input validation
	if err := validateTextToSpeechInput(text, targetLanguage, params); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var payload = map[string]any{
		"text":                 text,
		"target_language_code": targetLanguage,
	}
	if params.Speaker != nil {
		payload["speaker"] = *params.Speaker
	}
	// TODO: Add constraints as per the API docs for pitch, pace, etc...
	if params.Pitch != nil {
		payload["pitch"] = *params.Pitch
	}
	if params.Pace != nil {
		payload["pace"] = *params.Pace
	}
	if params.Loudness != nil {
		payload["loudness"] = *params.Loudness
	}
	if params.SpeechSampleRate != nil {
		payload["speech_sample_rate"] = *params.SpeechSampleRate
	}
	if params.EnablePreprocessing != nil {
		payload["enable_preprocessing"] = *params.EnablePreprocessing
	}
	if params.Model != nil {
		payload["model"] = *params.Model
	}

	resp, err := c.makeJsonHTTPRequest(ctx, http.MethodPost, c.baseURL+"/text-to-speech", payload)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
	}

	type textToSpeechResponse struct {
		RequestId string   `json:"request_id"`
		Audios    []string `json:"audios"`
	}

	var response textToSpeechResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &TextToSpeechResponse{
		RequestId: response.RequestId,
		Audios:    response.Audios,
	}, nil
}

// convertAndConcatBase64ToBytes converts multiple base64-encoded audio chunks to a single byte array.
func convertAndConcatBase64ToBytes(base64Strs []string) ([]byte, error) {
	var data []byte
	for _, base64Str := range base64Strs {
		decodedBytes, err := convertBase64ToBytes(base64Str)
		if err != nil {
			return nil, err
		}
		data = append(data, decodedBytes...)
	}
	return data, nil
}

// convertBase64ToBytes converts a single base64-encoded string to bytes.
func convertBase64ToBytes(base64Str string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}

// validateTextToSpeechInput validates inputs for text-to-speech requests
func validateTextToSpeechInput(text string, targetLanguage Language, params TextToSpeechParams) error {
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("text cannot be empty")
	}
	if !utf8.ValidString(text) {
		return fmt.Errorf("text must be valid UTF-8")
	}
	if targetLanguage == "" {
		return fmt.Errorf("target language cannot be empty")
	}

	// Validate text length (typical limit for TTS systems)
	if utf8.RuneCountInString(text) > 5000 {
		return fmt.Errorf("text is too long (max 5000 characters)")
	}

	// Validate optional parameters
	if params.Pitch != nil && (*params.Pitch < 0.5 || *params.Pitch > 2.0) {
		return fmt.Errorf("pitch must be between 0.5 and 2.0")
	}
	if params.Pace != nil && (*params.Pace < 0.5 || *params.Pace > 2.0) {
		return fmt.Errorf("pace must be between 0.5 and 2.0")
	}
	if params.Loudness != nil && (*params.Loudness < 0.5 || *params.Loudness > 2.0) {
		return fmt.Errorf("loudness must be between 0.5 and 2.0")
	}

	return nil
}
