package sarvam

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
)

// TextToSpeech represents the result of a text-to-speech operation.
type TextToSpeech struct {
	RequestId string
	Audios    []string
	Data      []byte
}

// Save saves the text-to-speech data as a WAV file
func (t *TextToSpeech) Save(filename string) error {
	return os.WriteFile(filename+".wav", t.Data, 0644)
}

type TextToSpeechParams struct {
	Text                string
	TargetLanguage      Language
	Speaker             *Speaker
	Pitch               *float64
	Pace                *float64
	Loudness            *float64
	SpeechSampleRate    *SpeechSampleRate
	EnablePreprocessing *bool
	Model               *Model
}

type SpeechSampleRate int

var (
	SpeechSampleRate8000  SpeechSampleRate = 8000
	SpeechSampleRate16000 SpeechSampleRate = 16000
	SpeechSampleRate22050 SpeechSampleRate = 22050
	SpeechSampleRate24000 SpeechSampleRate = 24000
)

// TextToSpeech converts text to speech in the specified language.
func (c *Client) TextToSpeech(params TextToSpeechParams) (*TextToSpeech, error) {
	var payload = map[string]any{
		"text":                 params.Text,
		"target_language_code": params.TargetLanguage,
	}
	if params.Speaker != nil {
		payload["voice"] = *params.Speaker
	}
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

	resp, err := c.makeHTTPRequest(http.MethodPost, c.baseURL+"/text-to-speech", payload)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	type textToSpeechResponse struct {
		RequestId string   `json:"request_id"`
		Audios    []string `json:"audios"`
	}

	var response textToSpeechResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	data, err := convertAndConcatBase64ToBytes(response.Audios)
	if err != nil {
		return nil, err
	}
	return &TextToSpeech{
		RequestId: response.RequestId,
		Audios:    response.Audios,
		Data:      data,
	}, nil
}

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

func convertBase64ToBytes(base64Str string) ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}
