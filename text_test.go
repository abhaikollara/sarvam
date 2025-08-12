package sarvam

import (
	"testing"
)

func TestTranslateWithParams(t *testing.T) {
	client := NewClient("test-key")

	// Test basic translation
	_, err := client.Translate("Hello", LanguageEnglish, LanguageHindi, nil)
	if err == nil {
		t.Error("Expected error for invalid API key, got nil")
	}

	// Test translation with parameters
	speakerGender := SpeakerGenderMale
	mode := TranslationModeFormal
	model := TranslationModelMayuraV1
	enablePreprocessing := true
	outputScript := OutputScriptRoman
	numeralsFormat := NumeralsFormatInternational

	options := &TranslateParams{
		SpeakerGender:       &speakerGender,
		Mode:                &mode,
		Model:               &model,
		EnablePreprocessing: &enablePreprocessing,
		OutputScript:        &outputScript,
		NumeralsFormat:      &numeralsFormat,
	}

	_, err = client.Translate("Hello", LanguageEnglish, LanguageHindi, options)
	if err == nil {
		t.Error("Expected error for invalid API key, got nil")
	}
}

func TestTranslateInputLengthValidation(t *testing.T) {
	client := NewClient("test-key")

	// Test input length validation for Mayura:v1 (1000 chars)
	longInput := string(make([]byte, 1001))
	speakerGender := SpeakerGenderMale
	mode := TranslationModeFormal
	model := TranslationModelMayuraV1

	options := &TranslateParams{
		SpeakerGender: &speakerGender,
		Mode:          &mode,
		Model:         &model,
	}

	_, err := client.Translate(longInput, LanguageEnglish, LanguageHindi, options)
	if err == nil {
		t.Error("Expected error for input too long, got nil")
	}

	// Test input length validation for Sarvam-Translate:v1 (2000 chars)
	veryLongInput := string(make([]byte, 2001))
	sarvamModel := TranslationModelSarvamTranslate

	sarvamOptions := &TranslateParams{
		Model: &sarvamModel,
	}

	_, err = client.Translate(veryLongInput, LanguageEnglish, LanguageHindi, sarvamOptions)
	if err == nil {
		t.Error("Expected error for input too long, got nil")
	}
}

func TestTranslateParamsValidation(t *testing.T) {
	// Test that all constants are properly defined
	testCases := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"SpeakerGenderMale", string(SpeakerGenderMale), "Male"},
		{"SpeakerGenderFemale", string(SpeakerGenderFemale), "Female"},
		{"TranslationModeFormal", string(TranslationModeFormal), "formal"},
		{"TranslationModeModernColloquial", string(TranslationModeModernColloquial), "modern-colloquial"},
		{"TranslationModeClassicColloquial", string(TranslationModeClassicColloquial), "classic-colloquial"},
		{"TranslationModeCodeMixed", string(TranslationModeCodeMixed), "code-mixed"},
		{"TranslationModelMayuraV1", string(TranslationModelMayuraV1), "mayura:v1"},
		{"TranslationModelSarvamTranslate", string(TranslationModelSarvamTranslate), "sarvam-translate:v1"},
		{"OutputScriptRoman", string(OutputScriptRoman), "roman"},
		{"OutputScriptFullyNative", string(OutputScriptFullyNative), "fully-native"},
		{"OutputScriptSpokenFormInNative", string(OutputScriptSpokenFormInNative), "spoken-form-in-native"},
		{"NumeralsFormatInternational", string(NumeralsFormatInternational), "international"},
		{"NumeralsFormatNative", string(NumeralsFormatNative), "native"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.value != tc.expected {
				t.Errorf("Expected %s to be %s, got %s", tc.name, tc.expected, tc.value)
			}
		})
	}
}

func TestTranslationString(t *testing.T) {
	translation := &TranslationResponse{
		RequestId:      "test-id",
		TranslatedText: "नमस्ते",
		SourceLanguage: LanguageEnglish,
	}

	expected := "नमस्ते"
	if translation.String() != expected {
		t.Errorf("Expected String() to return %s, got %s", expected, translation.String())
	}
}
