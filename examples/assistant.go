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
	// Load .env for testing
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env not found")
	}

	apiKey := os.Getenv("PIXIGPT_API_KEY")
	baseURL := os.Getenv("PIXIGPT_BASE_URL")

	if apiKey == "" || baseURL == "" {
		log.Fatal("Missing required environment variables")
	}

	// Create client
	c := client.New(apiKey, baseURL)
	ctx := context.Background()

	// List all assistants
	assistants, err := c.ListAssistants(ctx)
	if err != nil {
		log.Fatalf("Failed to list assistants: %v", err)
	}

	fmt.Printf("Found %d assistants:\n", len(assistants))
	for _, a := range assistants {
		fmt.Printf("  - %s (%s)\n", a.Name, a.ID)
		fmt.Printf("    Instructions: %s\n", a.Instructions)
	}

	// Create a new assistant
	newAssistant, err := c.CreateAssistant(ctx,
		"Test Assistant",
		"You are a helpful assistant created via the Go client.",
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to create assistant: %v", err)
	}
	fmt.Printf("\nCreated new assistant: %s (%s)\n", newAssistant.Name, newAssistant.ID)

	// Update the assistant
	updatedAssistant, err := c.UpdateAssistant(ctx,
		newAssistant.ID,
		"Updated Test Assistant",
		"You are a helpful assistant that was just updated.",
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to update assistant: %v", err)
	}
	fmt.Printf("Updated assistant: %s\n", updatedAssistant.Name)

	// Delete the assistant
	if err := c.DeleteAssistant(ctx, newAssistant.ID); err != nil {
		log.Fatalf("Failed to delete assistant: %v", err)
	}
	fmt.Printf("Deleted assistant: %s\n", newAssistant.ID)
}
