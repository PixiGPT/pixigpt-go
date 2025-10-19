package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fycat/pixigpt-go/client"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env for testing convenience
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env not found, using environment variables")
	}

	apiKey := os.Getenv("PIXIGPT_API_KEY")
	baseURL := os.Getenv("PIXIGPT_BASE_URL")
	assistantID := os.Getenv("DEFAULT_ASSISTANT_ID")

	if apiKey == "" || baseURL == "" || assistantID == "" {
		log.Fatal("Missing required environment variables: PIXIGPT_API_KEY, PIXIGPT_BASE_URL, DEFAULT_ASSISTANT_ID")
	}

	// Create client
	c := client.New(apiKey, baseURL)

	// Send chat completion request
	ctx := context.Background()
	resp, err := c.CreateChatCompletion(ctx, client.ChatCompletionRequest{
		AssistantID: assistantID,
		Messages: []client.Message{
			{Role: "user", Content: "Hello! What's your name?"},
		},
		Temperature: 0.7,
		MaxTokens:   2000,
	})
	if err != nil {
		log.Fatalf("Chat completion failed: %v", err)
	}

	// Print response
	if len(resp.Choices) > 0 {
		fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
		fmt.Printf("\nUsage: %d input + %d output = %d total tokens\n",
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens,
			resp.Usage.TotalTokens)
	}
}
