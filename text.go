package sarvam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SpeakerGender represents the gender of the speaker for better translations.
type SpeakerGender string

const (
	SpeakerGenderMale   SpeakerGender = "Male"
	SpeakerGenderFemale SpeakerGender = "Female"
)

// TranslationMode specifies the tone or style of the translation.
type TranslationMode string

const (
	// TranslationModeFormal represents formal translation style.
	TranslationModeFormal TranslationMode = "formal"
	// TranslationModeModernColloquial represents modern colloquial translation style.
	TranslationModeModernColloquial TranslationMode = "modern-colloquial"
	// TranslationModeClassicColloquial represents classic colloquial translation style.
	TranslationModeClassicColloquial TranslationMode = "classic-colloquial"
	// TranslationModeCodeMixed represents code-mixed translation style.
	TranslationModeCodeMixed TranslationMode = "code-mixed"
)

// OutputScript controls the transliteration style applied to the output text.
type OutputScript string

const (
	// OutputScriptRoman represents Roman script output.
	OutputScriptRoman OutputScript = "roman"
	// OutputScriptFullyNative represents fully native script output.
	OutputScriptFullyNative OutputScript = "fully-native"
	// OutputScriptSpokenFormInNative represents spoken form in native script output.
	OutputScriptSpokenFormInNative OutputScript = "spoken-form-in-native"
)

// NumeralsFormat specifies the format for numerals in the translation.
type NumeralsFormat string

const (
	// NumeralsFormatInternational represents international numeral format.
	NumeralsFormatInternational NumeralsFormat = "international"
	// NumeralsFormatNative represents native numeral format.
	NumeralsFormatNative NumeralsFormat = "native"
)

// TranslationResponse represents the result of a translation operation.
type TranslationResponse struct {
	RequestId      string
	TranslatedText string
	SourceLanguage Language
}

// String returns the translated text.
func (t *TranslationResponse) String() string {
	return string(t.TranslatedText)
}

// TranslateParams contains all optional parameters for translation.
type TranslateParams struct {
	SpeakerGender       *SpeakerGender
	Mode                *TranslationMode
	Model               *TranslationModel
	EnablePreprocessing *bool
	OutputScript        *OutputScript
	NumeralsFormat      *NumeralsFormat
}

// TranslateWithParams converts text from one language to another with custom parameters.
func (c *Client) Translate(input string, sourceLanguageCode, targetLanguageCode Language, params *TranslateParams) (*TranslationResponse, error) {
	// Validate input length based on model
	maxLength := 2000 // Default for sarvam-translate:v1
	if params != nil && params.Model != nil && *params.Model == TranslationModelMayuraV1 {
		maxLength = 1000
	}

	if l := len(input); l > maxLength {
		return nil, &ErrInputTooLong{
			InputLength: l,
			MaxLength:   maxLength,
		}
	}

	var reqBody = map[string]any{
		"input":                input,
		"source_language_code": sourceLanguageCode,
		"target_language_code": targetLanguageCode,
	}

	// Add optional parameters if provided
	if params != nil {
		if params.SpeakerGender != nil {
			reqBody["speaker_gender"] = *params.SpeakerGender
		}
		if params.Mode != nil {
			reqBody["mode"] = *params.Mode
		}
		if params.Model != nil {
			reqBody["model"] = *params.Model
		}
		if params.EnablePreprocessing != nil {
			reqBody["enable_preprocessing"] = *params.EnablePreprocessing
		}
		if params.OutputScript != nil {
			reqBody["output_script"] = *params.OutputScript
		}
		if params.NumeralsFormat != nil {
			reqBody["numerals_format"] = *params.NumeralsFormat
		}
	}

	resp, err := c.makeJsonHTTPRequest(http.MethodPost, c.baseURL+"/translate", reqBody)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp)
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

	return &TranslationResponse{
		RequestId:      response.RequestId,
		TranslatedText: response.TranslatedText,
		SourceLanguage: mapLanguageCodeToLanguage(response.SourceLanguage),
	}, nil
}

// LanguageIdentification represents the result of language identification.
type LanguageIdentificationResponse struct {
	RequestId string
	Language  Language
	Script    Script
}

// IdentifyLanguage identifies the language (e.g., en-IN, hi-IN) and script (e.g., Latin, Devanagari) of the input text, supporting multiple languages.
func (c *Client) IdentifyLanguage(input string) (*LanguageIdentificationResponse, error) {
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
		return nil, parseAPIError(resp)
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

	return &LanguageIdentificationResponse{
		RequestId: response.RequestId,
		Language:  mapLanguageCodeToLanguage(response.LanguageCode),
		Script:    mapScriptCodeToScript(response.ScriptCode),
	}, nil
}

// TransliterationResponse represents the result of a transliteration operation.
type TransliterationResponse struct {
	RequestId          string
	TransliteratedText string
	SourceLanguage     Language
}

// String returns the transliterated text.
func (t *TransliterationResponse) String() string {
	return string(t.TransliteratedText)
}

// Transliterate converts text from one script to another while preserving the original pronunciation.
// TODO: There are more params. See docs. Add them. This would change the signature I guess.
func (c *Client) Transliterate(input string, sourceLanguage Language, targetLanguage Language) (*TransliterationResponse, error) {
	if l := len(input); l >= 1000 {
		return nil, &ErrInputTooLong{
			InputLength: l,
			MaxLength:   1000,
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
		return nil, parseAPIError(resp)
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

	return &TransliterationResponse{
		RequestId:          response.RequestId,
		TransliteratedText: response.TransliteratedText,
		SourceLanguage:     mapLanguageCodeToLanguage(response.SourceLanguage),
	}, nil
}

// ErrInputTooLong is returned when the input length is greater than or equal to 1000 characters.
type ErrInputTooLong struct {
	InputLength int
	MaxLength   int
}

func (e *ErrInputTooLong) Error() string {
	return fmt.Sprintf("input length must be less than %d characters, got %d", e.MaxLength, e.InputLength)
}
