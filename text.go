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
	TranslationModeFormal            TranslationMode = "formal"
	TranslationModeModernColloquial  TranslationMode = "modern-colloquial"
	TranslationModeClassicColloquial TranslationMode = "classic-colloquial"
	TranslationModeCodeMixed         TranslationMode = "code-mixed"
)

// TranslationModel specifies the translation model to use.
type TranslationModel string

const (
	TranslationModelMayuraV1        TranslationModel = "mayura:v1"
	TranslationModelSarvamTranslate TranslationModel = "sarvam-translate:v1"
)

// OutputScript controls the transliteration style applied to the output text.
type OutputScript string

const (
	OutputScriptRoman              OutputScript = "roman"
	OutputScriptFullyNative        OutputScript = "fully-native"
	OutputScriptSpokenFormInNative OutputScript = "spoken-form-in-native"
)

// NumeralsFormat specifies the format for numerals in the translation.
type NumeralsFormat string

const (
	NumeralsFormatInternational NumeralsFormat = "international"
	NumeralsFormatNative        NumeralsFormat = "native"
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
// This is the simple version that uses default options.
func (c *Client) Translate(input string, sourceLanguageCode, targetLanguageCode Language) (*Translation, error) {
	return c.TranslateWithOptions(&TranslateParams{
		Input:          input,
		SourceLanguage: sourceLanguageCode,
		TargetLanguage: targetLanguageCode,
	})
}

// TranslateParams contains all optional parameters for translation.
type TranslateParams struct {
	Input               string
	SourceLanguage      Language
	TargetLanguage      Language
	SpeakerGender       *SpeakerGender
	Mode                *TranslationMode
	Model               *TranslationModel
	EnablePreprocessing *bool
	OutputScript        *OutputScript
	NumeralsFormat      *NumeralsFormat
}

// TranslateWithOptions converts text from one language to another with custom options.
func (c *Client) TranslateWithOptions(params *TranslateParams) (*Translation, error) {
	// Validate input length based on model
	maxLength := 2000 // Default for sarvam-translate:v1
	if params != nil && params.Model != nil && *params.Model == TranslationModelMayuraV1 {
		maxLength = 1000
	}

	if l := len(params.Input); l > maxLength {
		return nil, &ErrInputTooLong{
			InputLength: l,
			MaxLength:   maxLength,
		}
	}

	var reqBody = map[string]any{
		"input":                params.Input,
		"source_language_code": params.SourceLanguage,
		"target_language_code": params.TargetLanguage,
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

	return &Transliteration{
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
