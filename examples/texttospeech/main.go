package main

import (
	"context"
	"log"
	"os"

	"code.abhai.dev/sarvam"
)

func main() {
	apiKey := os.Getenv("SARVAM_API_KEY")
	client := sarvam.NewClient(apiKey)

	response, err := client.TextToSpeech(context.Background(), "नमस्ते, आप कैसे हैं?", sarvam.LanguageHindi, sarvam.TextToSpeechParams{
		Speaker: &sarvam.SpeakerHitesh,
		Model:   &sarvam.TextToSpeechModelBulbulV2,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	response.Save("output.wav")
}
