package sarvam

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
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
func (c *Client) TextToSpeech(text string, targetLanguage Language, params TextToSpeechParams) (*TextToSpeechResponse, error) {
	var payload = map[string]any{
		"text":                 text,
		"target_language_code": targetLanguage,
	}
	if params.Speaker != nil {
		payload["voice"] = *params.Speaker
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

	resp, err := c.makeJsonHTTPRequest(http.MethodPost, c.baseURL+"/text-to-speech", payload)
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
