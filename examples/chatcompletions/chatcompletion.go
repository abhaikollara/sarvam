package main

import (
	"fmt"
	"log"
	"os"

	"code.abhai.dev/sarvam"
)

func main() {
	apiKey := os.Getenv("SARVAM_API_KEY")
	if apiKey == "" {
		log.Fatal("SARVAM_API_KEY environment variable is required")
	}

	// Example 1: Using package-level functions (new way)
	fmt.Println("=== Using Package-Level Functions ===")
	sarvam.SetAPIKey(apiKey)

	response, err := sarvam.ChatCompletion(&sarvam.ChatCompletionRequest{
		Model: sarvam.ChatCompletionModelSarvamM,
		Messages: []sarvam.Message{
			{Role: "user", Content: "ഹലോ, നിങ്ങൾക്ക് സുഖമാണോ"},
		},
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Package-level response:", response.Choices[0].Message.Content)

	// Example 2: Using client instance (original way)
	fmt.Println("\n=== Using Client Instance ===")
	client := sarvam.NewClient(apiKey)

	response2, err := client.ChatCompletion(&sarvam.ChatCompletionRequest{
		Model: sarvam.ChatCompletionModelSarvamM,
		Messages: []sarvam.Message{
			{Role: "user", Content: "Hello, how are you?"},
		},
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Client instance response:", response2.Choices[0].Message.Content)
}
