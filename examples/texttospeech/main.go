package main

import (
	"log"
	"os"

	"code.abhai.dev/sarvam"
)

func main() {
	apiKey := os.Getenv("SARVAM_API_KEY")
	client := sarvam.NewClient(apiKey)

	response, err := client.TextToSpeech(sarvam.TextToSpeechParams{
		Text:           "नमस्ते, आप कैसे हैं?",
		Speaker:        &sarvam.SpeakerAnushka,
		TargetLanguage: sarvam.LanguageHindi,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	response.Save("output")
}
