package main

import (
	"context"
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

	clientExample(apiKey, filepath)
	defaultClientExample(filepath)
}

func clientExample(apiKey string, filepath string) {
	fmt.Println("Using client")

	client := sarvam.NewClient(apiKey)
	speechFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer speechFile.Close()

	response, err := client.SpeechToText(context.Background(), speechFile, sarvam.SpeechToTextParams{
		Language:       sarvam.Ptr(sarvam.LanguageMalayalam),
		Model:          sarvam.Ptr(sarvam.SpeechToTextModelSaarikaV2),
		WithTimestamps: sarvam.Ptr(true),
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Converting file: %s\n", filepath)
	fmt.Printf("Transcript: %s\n", response.Transcript)
	fmt.Printf("Timestamps: %v\n", response.Timestamps)
}

func defaultClientExample(filepath string) {
	fmt.Println("Using default client")
	speechFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer speechFile.Close()

	response, err := sarvam.SpeechToText(context.Background(), speechFile, sarvam.SpeechToTextParams{
		Language:       sarvam.Ptr(sarvam.LanguageMalayalam),
		Model:          sarvam.Ptr(sarvam.SpeechToTextModelSaarikaV2),
		WithTimestamps: sarvam.Ptr(true),
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Converting file: %s\n", filepath)
	fmt.Printf("Transcript: %s\n", response.Transcript)
	fmt.Printf("Timestamps: %v\n", response.Timestamps)
}
