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

	response, err := client.IdentifyLanguage(context.Background(), "सभी पुरुषों को सेवा करनी चाहिए")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Language: %s\n", response.Language)
}
