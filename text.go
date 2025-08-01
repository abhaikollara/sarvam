package sarvam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Translation represents the result of a translation operation.
type Translation struct {
	RequestId      string
	TranslatedText string
	SourceLanguage Language
}

// String returns the translated text.
func (t *Translation) String() string {
	return string(t.TranslatedText)
}

// Translate converts text from one language to another while preserving its meaning.
func (c *Client) Translate(input string, sourceLanguageCode, targetLanguageCode Language) (*Translation, error) {
	var reqBody = map[string]any{
		"input":                input,
		"source_language_code": sourceLanguageCode,
		"target_language_code": targetLanguageCode,
	}
	resp, err := c.makeJsonHTTPRequest(http.MethodPost, c.baseURL+"/translate", reqBody)
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

	type translateResponse struct {
		RequestId      string `json:"request_id"`
		TranslatedText string `json:"translated_text"`
		SourceLanguage string `json:"source_language_code"`
	}

	var response translateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &Translation{
		RequestId:      response.RequestId,
		TranslatedText: response.TranslatedText,
		SourceLanguage: mapLanguageCodeToLanguage(response.SourceLanguage),
	}, nil
}

// LanguageIdentification represents the result of a language identification operation.
type LanguageIdentification struct {
	RequestId string
	Language  Language
	Script    string
}

// IdentifyLanguage identifies the language (e.g., en-IN, hi-IN) and script (e.g., Latin, Devanagari) of the input text, supporting multiple languages.
func (c *Client) IdentifyLanguage(input string) (*LanguageIdentification, error) {
	var payload = map[string]string{
		"input": input,
	}
	resp, err := c.makeJsonHTTPRequest(http.MethodPost, c.baseURL+"/text-lid", payload)
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

	type identifyLanguageResponse struct {
		RequestId    string `json:"request_id"`
		LanguageCode string `json:"language_code"`
		ScriptCode   string `json:"script_code"`
	}

	var response identifyLanguageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &LanguageIdentification{
		RequestId: response.RequestId,
		Language:  mapLanguageCodeToLanguage(response.LanguageCode),
		Script:    response.ScriptCode,
	}, nil
}

// Transliteration represents the result of a transliteration operation.
type Transliteration struct {
	RequestId          string
	TransliteratedText string
	SourceLanguage     Language
}

// String returns the transliterated text.
func (t *Transliteration) String() string {
	return string(t.TransliteratedText)
}

// Transliterate converts text from one script to another while preserving the original pronunciation.
func (c *Client) Transliterate(input string, sourceLanguage Language, targetLanguage Language) (*Transliteration, error) {
	if l := len(input); l >= 1000 {
		return nil, &ErrInputTooLong{
			InputLength: l,
		}
	}

	var payload = map[string]any{
		"input":                input,
		"source_language_code": sourceLanguage,
		"target_language_code": targetLanguage,
	}
	resp, err := c.makeJsonHTTPRequest(http.MethodPost, c.baseURL+"/transliterate", payload)
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

	type transliterationResponse struct {
		RequestId          string `json:"request_id"`
		TransliteratedText string `json:"transliterated_text"`
		SourceLanguage     string `json:"source_language_code"`
	}

	var response transliterationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &Transliteration{
		RequestId:          response.RequestId,
		TransliteratedText: response.TransliteratedText,
		SourceLanguage:     mapLanguageCodeToLanguage(response.SourceLanguage),
	}, nil
}

// ErrInputTooLong is returned when the input length is greater than or equal to 1000 characters.
type ErrInputTooLong struct {
	InputLength int
}

func (e *ErrInputTooLong) Error() string {
	return fmt.Sprintf("input length must be less than 1000 characters, got %d", e.InputLength)
}
