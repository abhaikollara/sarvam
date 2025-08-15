package main

import (
	"fmt"
	"log"
	"os"

	"code.abhai.dev/sarvam"
)

func main() {
	client := sarvam.NewClient(os.Getenv("SARVAM_API_KEY"))

	response, err := client.Transliterate("enthokke und vishesham", sarvam.LanguageEnglish, sarvam.LanguageMalayalam, &sarvam.TransliterateParams{
		NumeralsFormat: sarvam.Ptr(sarvam.NumeralsFormatInternational),
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Transliterated text: %s\n", response.TransliteratedText) // Transliterated text: എന്തൊക്കെ ഉണ്ട് വിശേഷം
}
