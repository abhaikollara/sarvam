package sarvam

import (
	"testing"
)

func TestTranslateWithOptions(t *testing.T) {
	client := NewClient("test-key")

	// Test basic translation
	_, err := client.Translate("Hello", LanguageEnglish, LanguageHindi)
	if err == nil {
		t.Error("Expected error for invalid API key, got nil")
	}

	// Test translation with options
	speakerGender := SpeakerGenderMale
	mode := TranslationModeFormal
	model := TranslationModelMayuraV1
	enablePreprocessing := true
	outputScript := OutputScriptRoman
	numeralsFormat := NumeralsFormatInternational

	options := &TranslateOptions{
		SpeakerGender:       &speakerGender,
		Mode:                &mode,
		Model:               &model,
		EnablePreprocessing: &enablePreprocessing,
		OutputScript:        &outputScript,
		NumeralsFormat:      &numeralsFormat,
	}

	_, err = client.TranslateWithOptions("Hello", LanguageEnglish, LanguageHindi, options)
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

	options := &TranslateOptions{
		SpeakerGender: &speakerGender,
		Mode:          &mode,
		Model:         &model,
	}

	_, err := client.TranslateWithOptions(longInput, LanguageEnglish, LanguageHindi, options)
	if err == nil {
		t.Error("Expected error for input too long, got nil")
	}

	// Test input length validation for Sarvam-Translate:v1 (2000 chars)
	veryLongInput := string(make([]byte, 2001))
	sarvamModel := TranslationModelSarvamTranslate

	sarvamOptions := &TranslateOptions{
		Model: &sarvamModel,
	}

	_, err = client.TranslateWithOptions(veryLongInput, LanguageEnglish, LanguageHindi, sarvamOptions)
	if err == nil {
		t.Error("Expected error for input too long, got nil")
	}
}

func TestTranslateOptionsValidation(t *testing.T) {
	// Test that all constants are properly defined
	testCases := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"SpeakerGenderMale", SpeakerGenderMale, "Male"},
		{"SpeakerGenderFemale", SpeakerGenderFemale, "Female"},
		{"TranslationModeFormal", TranslationModeFormal, "formal"},
		{"TranslationModeModernColloquial", TranslationModeModernColloquial, "modern-colloquial"},
		{"TranslationModeClassicColloquial", TranslationModeClassicColloquial, "classic-colloquial"},
		{"TranslationModeCodeMixed", TranslationModeCodeMixed, "code-mixed"},
		{"TranslationModelMayuraV1", TranslationModelMayuraV1, "mayura:v1"},
		{"TranslationModelSarvamTranslate", TranslationModelSarvamTranslate, "sarvam-translate:v1"},
		{"OutputScriptRoman", OutputScriptRoman, "roman"},
		{"OutputScriptFullyNative", OutputScriptFullyNative, "fully-native"},
		{"OutputScriptSpokenFormInNative", OutputScriptSpokenFormInNative, "spoken-form-in-native"},
		{"NumeralsFormatInternational", NumeralsFormatInternational, "international"},
		{"NumeralsFormatNative", NumeralsFormatNative, "native"},
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
	translation := &Translation{
		RequestId:      "test-id",
		TranslatedText: "नमस्ते",
		SourceLanguage: LanguageEnglish,
	}

	expected := "नमस्ते"
	if translation.String() != expected {
		t.Errorf("Expected String() to return %s, got %s", expected, translation.String())
	}
}
