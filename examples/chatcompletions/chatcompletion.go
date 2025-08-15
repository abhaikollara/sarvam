package main

import (
	"fmt"
	"log"
	"os"

	"code.abhai.dev/sarvam"
)

func main() {
	client := sarvam.NewClient(os.Getenv("SARVAM_API_KEY"))

	response, err := client.ChatCompletion([]sarvam.Message{
		sarvam.NewSystemMessage("You are a helpful assistant that answers in short"),
		sarvam.NewUserMessage("Hello, how are you?"),
	}, sarvam.ChatCompletionModelSarvamM, &sarvam.ChatCompletionParams{})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Response: %v\n", response.GetFirstChoiceContent())
}
