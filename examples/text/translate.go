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

	response, err := client.Translate("Hello, how are you?", sarvam.LanguageEnglish, sarvam.LanguageHindi)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Translated Text: %s\n", response)

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
