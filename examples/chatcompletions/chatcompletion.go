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

	// Example 3: Using ChatCompletionWithOptions with proper params struct
	fmt.Println("\n=== Using ChatCompletionWithOptions ===")

	// Create params with custom options
	optionsParams := sarvam.ChatCompletionWithOptionsParams{
		Model: sarvam.ChatCompletionModelSarvamM,
		Messages: []sarvam.Message{
			{Role: "user", Content: "Explain quantum computing in simple terms"},
		},
		Temperature: sarvam.Float64(0.7),
		MaxTokens:   sarvam.Int(150),
		TopP:        sarvam.Float64(0.9),
	}

	response3, err := sarvam.ChatCompletionWithOptions(optionsParams)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("ChatCompletionWithOptions response:", response3.Choices[0].Message.Content)
}
