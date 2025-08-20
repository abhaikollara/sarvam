package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"code.abhai.dev/sarvam"
)

func main() {
	client := sarvam.NewClient(os.Getenv("SARVAM_API_KEY"))

	response, err := client.Translate(context.Background(), "All men must serve", sarvam.LanguageEnglish, sarvam.LanguageHindi, &sarvam.TranslateParams{
		OutputScript: sarvam.Ptr(sarvam.OutputScriptFullyNative),
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Translated text: %s\n", response.TranslatedText)
}
