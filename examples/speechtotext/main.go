package main

import (
	"fmt"
	"log"
	"os"

	"code.abhai.dev/sarvam"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("SARVAM_API_KEY")
	if apiKey == "" {
		log.Fatal("SARVAM_API_KEY environment variable is required")
	}
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <audio_file_path>")
	}
	filepath := os.Args[1]
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Fatalf("File not found: %s", filepath)
	}

	// Create a new client
	client := sarvam.NewClient(apiKey)

	// Example 1: Basic speech-to-text
	fmt.Println("=== Speech-to-Text Example ===")
	params := sarvam.SpeechToTextParams{
		FilePath: filepath,
		Model:    &sarvam.SpeechToTextModelSaarikaV2dot5,
	}

	result, err := client.SpeechToText(params)
	if err != nil {
		log.Fatalf("Speech-to-text failed: %v", err)
	}

	fmt.Printf("Request ID: %s\n", result.RequestId)
	fmt.Printf("Transcript: %s\n", result.Transcript)
	fmt.Printf("Language Code: %s\n", result.LanguageCode)

	// Example 2: Speech-to-text with timestamps
	fmt.Println("\n=== Speech-to-Text with Timestamps ===")
	paramsWithTimestamps := sarvam.SpeechToTextParams{
		FilePath:       filepath,
		Model:          &sarvam.SpeechToTextModelSaarikaV2dot5,
		WithTimestamps: sarvam.Bool(true),
	}

	resultWithTimestamps, err := client.SpeechToText(paramsWithTimestamps)
	if err != nil {
		log.Fatalf("Speech-to-text with timestamps failed: %v", err)
	}

	fmt.Printf("Request ID: %s\n", resultWithTimestamps.RequestId)
	fmt.Printf("Transcript: %s\n", resultWithTimestamps.Transcript)
	if resultWithTimestamps.Timestamps != nil {
		fmt.Printf("Number of words with timestamps: %d\n", len(resultWithTimestamps.Timestamps.Words))
	}

	// Example 3: Speech-to-text-translate (auto-detect language and translate to English)
	fmt.Println("\n=== Speech-to-Text Translate Example ===")
	translateParams := sarvam.SpeechToTextTranslateParams{
		FilePath: filepath,
		Model:    &sarvam.SpeechToTextTranslateModelSaarasV2dot5,
		Prompt:   sarvam.String("This is a conversation is a greeting"),
	}

	translateResult, err := client.SpeechToTextTranslate(translateParams)
	if err != nil {
		log.Fatalf("Speech-to-text-translate failed: %v", err)
	}

	fmt.Printf("Request ID: %s\n", translateResult.RequestId)
	fmt.Printf("Translated Transcript: %s\n", translateResult.Transcript)
	fmt.Printf("Detected Language Code: %s\n", translateResult.LanguageCode)
}
