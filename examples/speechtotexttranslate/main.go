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

	response, err := client.SpeechToTextTranslate(speechFile, sarvam.SpeechToTextTranslateParams{
		Model:      sarvam.Ptr(sarvam.SpeechToTextTranslateModelSaarasV2),
		AudioCodec: sarvam.Ptr(sarvam.AudioCodecWav),
		Prompt:     sarvam.Ptr("This is a greeting"),
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Converting file: %s\n", filepath)
	fmt.Printf("Transcript (translated to English): %s\n", response.Transcript)
	fmt.Printf("Detected source language: %s\n", response.Language)
	if response.DiarizedTranscript != nil {
		fmt.Printf("Diarized transcript: %v\n", response.DiarizedTranscript)
	}
}

func defaultClientExample(filepath string) {
	fmt.Println("Using default client")
	speechFile, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer speechFile.Close()

	response, err := sarvam.SpeechToTextTranslate(speechFile, sarvam.SpeechToTextTranslateParams{
		Model:      sarvam.Ptr(sarvam.SpeechToTextTranslateModelSaarasV2dot5),
		AudioCodec: sarvam.Ptr(sarvam.AudioCodecWav),
		Prompt:     sarvam.Ptr("This is a greeting"),
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Converting file: %s\n", filepath)
	fmt.Printf("Transcript (translated to English): %s\n", response.Transcript)
	fmt.Printf("Detected source language: %s\n", response.Language)
	if response.DiarizedTranscript != nil {
		fmt.Printf("Diarized transcript: %v\n", response.DiarizedTranscript)
	}
}
