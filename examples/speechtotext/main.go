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

	// Example 1: Using package-level functions (new way)
	fmt.Println("=== Using Package-Level Functions ===")
	sarvam.SetAPIKey(apiKey)

	// Basic speech-to-text using package-level function
	result, err := sarvam.SpeechToText(sarvam.SpeechToTextParams{
		FilePath: filepath,
		Model:    &sarvam.SpeechToTextModelSaarikaV2dot5,
	})
	if err != nil {
		log.Fatalf("Speech-to-text failed: %v", err)
	}

	fmt.Printf("Request ID: %s\n", result.RequestId)
	fmt.Printf("Transcript: %s\n", result.Transcript)
	fmt.Printf("Language Code: %s\n", result.LanguageCode)

	// Speech-to-text-translate using package-level function
	translateResult, err := sarvam.SpeechToTextTranslate(sarvam.SpeechToTextTranslateParams{
		FilePath: filepath,
		Model:    &sarvam.SpeechToTextTranslateModelSaarasV2dot5,
		Prompt:   sarvam.String("This is a conversation is a greeting"),
	})
	if err != nil {
		log.Fatalf("Speech-to-text-translate failed: %v", err)
	}

	fmt.Printf("Request ID: %s\n", translateResult.RequestId)
	fmt.Printf("Translated Transcript: %s\n", translateResult.Transcript)
	fmt.Printf("Detected Language Code: %s\n", translateResult.LanguageCode)

	// Example 2: Using client instance (original way)
	fmt.Println("\n=== Using Client Instance ===")
	client := sarvam.NewClient(apiKey)

	// Basic speech-to-text
	params := sarvam.SpeechToTextParams{
		FilePath: filepath,
		Model:    &sarvam.SpeechToTextModelSaarikaV2dot5,
	}

	result2, err := client.SpeechToText(params)
	if err != nil {
		log.Fatalf("Speech-to-text failed: %v", err)
	}

	fmt.Printf("Request ID: %s\n", result2.RequestId)
	fmt.Printf("Transcript: %s\n", result2.Transcript)
	fmt.Printf("Language Code: %s\n", result2.LanguageCode)

	// Speech-to-text with timestamps
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
}
