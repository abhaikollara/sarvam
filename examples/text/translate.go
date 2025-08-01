package main

import (
	"fmt"
	"log"
	"os"

	"code.abhai.dev/sarvam"
)

func main() {
	apiKey := os.Getenv("SARVAM_API_KEY")
	client := sarvam.NewClient(apiKey)

	// Basic translation
	response, err := client.Translate("Hello, how are you?", sarvam.LanguageEnglish, sarvam.LanguageHindi)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Basic Translation: %s\n", response)

	// Advanced translation with options
	speakerGender := sarvam.SpeakerGenderMale
	mode := sarvam.TranslationModeModernColloquial
	model := sarvam.TranslationModelMayuraV1
	enablePreprocessing := true
	outputScript := sarvam.OutputScriptRoman
	numeralsFormat := sarvam.NumeralsFormatNative

	options := &sarvam.TranslateOptions{
		SpeakerGender:       &speakerGender,
		Mode:                &mode,
		Model:               &model,
		EnablePreprocessing: &enablePreprocessing,
		OutputScript:        &outputScript,
		NumeralsFormat:      &numeralsFormat,
	}

	advancedResponse, err := client.TranslateWithOptions(
		"Your EMI of Rs. 3000 is pending.",
		sarvam.LanguageEnglish,
		sarvam.LanguageHindi,
		options,
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Advanced Translation: %s\n", advancedResponse)

	// Translation with Sarvam-Translate model (supports more languages)
	sarvamModel := sarvam.TranslationModelSarvamTranslate
	sarvamMode := sarvam.TranslationModeFormal
	sarvamPreprocessing := true

	sarvamOptions := &sarvam.TranslateOptions{
		Model:               &sarvamModel,
		Mode:                &sarvamMode,
		EnablePreprocessing: &sarvamPreprocessing,
	}

	sarvamResponse, err := client.TranslateWithOptions(
		"मैं ऑफिस जा रहा हूँ",
		sarvam.LanguageHindi,
		sarvam.LanguageEnglish,
		sarvamOptions,
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Sarvam-Translate: %s\n", sarvamResponse)

	// identify language
	identifyResponse, err := client.IdentifyLanguage("namaskaaram")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Identified Language: %s, Script: %s\n", identifyResponse.Language, identifyResponse.Script)

	// transliterate
	transliterateResponse, err := client.Transliterate("namaskaaram", sarvam.LanguageAuto, sarvam.LanguageMalayalam)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Transliterated Text: %s\n", transliterateResponse)
}
