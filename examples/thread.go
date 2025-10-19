package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PixiGPT/pixigpt-go/client"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env for testing
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env not found")
	}

	apiKey := os.Getenv("PIXIGPT_API_KEY")
	baseURL := os.Getenv("PIXIGPT_BASE_URL")
	assistantID := os.Getenv("DEFAULT_ASSISTANT_ID")

	if apiKey == "" || baseURL == "" || assistantID == "" {
		log.Fatal("Missing required environment variables")
	}

	// Create client
	c := client.New(apiKey, baseURL)
	ctx := context.Background()

	// 1. Create thread
	thread, err := c.CreateThread(ctx)
	if err != nil {
		log.Fatalf("Failed to create thread: %v", err)
	}
	fmt.Printf("Created thread: %s\n", thread.ID)

	// 2. Add user message
	msg, err := c.CreateMessage(ctx, thread.ID, "user", "What's the meaning of life?")
	if err != nil {
		log.Fatalf("Failed to create message: %v", err)
	}
	fmt.Printf("Added message: %s\n", msg.ID)

	// 3. Create run (async)
	run, err := c.CreateRun(ctx, thread.ID, assistantID, true)
	if err != nil {
		log.Fatalf("Failed to create run: %v", err)
	}
	fmt.Printf("Created run: %s (status: %s)\n", run.ID, run.Status)

	// 4. Wait for completion
	fmt.Println("Waiting for run to complete...")
	completedRun, err := c.WaitForRun(ctx, thread.ID, run.ID)
	if err != nil {
		log.Fatalf("Run failed: %v", err)
	}
	fmt.Printf("Run completed: %s\n", completedRun.Status)

	// 5. Get messages
	messages, err := c.ListMessages(ctx, thread.ID, 10)
	if err != nil {
		log.Fatalf("Failed to list messages: %v", err)
	}

	// Print conversation (reversed - oldest first)
	fmt.Println("\n=== Conversation ===")
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if len(msg.Content) > 0 {
			fmt.Printf("%s: %s\n", msg.Role, msg.Content[0].Text.Value)
		}
	}
}
