package main

import (
	"log"
	"os"

	"code.abhai.dev/sarvam"
)

func main() {
	apiKey := os.Getenv("SARVAM_API_KEY")
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <audio_file_path>", os.Args[0])
	}
	filepath := os.Args[1]
	client := sarvam.NewClient(apiKey)

	response, err := client.SpeechToText(sarvam.SpeechToTextParams{
		FilePath: filepath,
		Model:    &sarvam.ModelSarvamM, // Optional: specify a model
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Request ID: %s", response.RequestId)
	log.Printf("Transcribed Text: %v", response)
}
