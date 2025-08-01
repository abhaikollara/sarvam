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

	response, err := client.ChatCompletion(&sarvam.ChatCompletionRequest{
		Model: sarvam.ChatCompletionModelSarvamM,
		Messages: []sarvam.Message{
			{Role: "user", Content: "ഹലോ, നിങ്ങൾക്ക് സുഖമാണോ"},
		},
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(response.Choices[0].Message.Content)
}
