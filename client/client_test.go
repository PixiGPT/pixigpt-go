package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var (
	testClient      *Client
	testAssistantID string
)

func init() {
	// Load .env from parent directory
	_ = godotenv.Load("../.env")

	apiKey := os.Getenv("PIXIGPT_API_KEY")
	baseURL := os.Getenv("PIXIGPT_BASE_URL")
	testAssistantID = os.Getenv("DEFAULT_ASSISTANT_ID")

	if apiKey == "" || baseURL == "" || testAssistantID == "" {
		panic("Missing test environment variables: PIXIGPT_API_KEY, PIXIGPT_BASE_URL, DEFAULT_ASSISTANT_ID")
	}

	testClient = New(apiKey, baseURL)
}

func TestChatCompletion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := testClient.CreateChatCompletion(ctx, ChatCompletionRequest{
		AssistantID: testAssistantID,
		Messages: []Message{
			{Role: "user", Content: "Say only 'TEST PASS' and nothing else."},
		},
		Temperature: 0.1, // Low temperature for consistent response
		MaxTokens:   50,
	})

	if err != nil {
		t.Fatalf("CreateChatCompletion failed: %v", err)
	}

	// Validate response structure
	if resp.ID == "" {
		t.Error("Response missing ID")
	}
	if resp.Object != "chat.completion" {
		t.Errorf("Expected object 'chat.completion', got %s", resp.Object)
	}
	if len(resp.Choices) == 0 {
		t.Fatal("Response has no choices")
	}

	content := resp.Choices[0].Message.Content
	t.Logf("Assistant response: %s", content)
	t.Logf("Token usage: %d input + %d output = %d total",
		resp.Usage.PromptTokens,
		resp.Usage.CompletionTokens,
		resp.Usage.TotalTokens)

	// Validate we got actual content
	if content == "" {
		t.Error("Response content is empty")
	}
	if resp.Usage.TotalTokens == 0 {
		t.Error("Token usage is zero")
	}
}

func TestChatCompletionWithThinkingDisabled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	enableThinking := false
	resp, err := testClient.CreateChatCompletion(ctx, ChatCompletionRequest{
		AssistantID:    testAssistantID,
		Messages:       []Message{{Role: "user", Content: "What is 2+2?"}},
		EnableThinking: &enableThinking,
		MaxTokens:      100,
	})

	if err != nil {
		t.Fatalf("CreateChatCompletion (no thinking) failed: %v", err)
	}

	t.Logf("Response (thinking disabled): %s", resp.Choices[0].Message.Content)
	t.Logf("Tokens: %d total", resp.Usage.TotalTokens)

	if resp.Choices[0].Message.Content == "" {
		t.Error("Response content is empty")
	}
}

func TestThreadWorkflow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 1. Create thread
	t.Log("Creating thread...")
	thread, err := testClient.CreateThread(ctx)
	if err != nil {
		t.Fatalf("CreateThread failed: %v", err)
	}
	t.Logf("Created thread: %s", thread.ID)

	if thread.ID == "" {
		t.Fatal("Thread ID is empty")
	}
	if thread.Object != "thread" {
		t.Errorf("Expected object 'thread', got %s", thread.Object)
	}

	// 2. Add user message
	t.Log("Adding message...")
	msg, err := testClient.CreateMessage(ctx, thread.ID, "user", "What is the capital of France? Answer in one word.")
	if err != nil {
		t.Fatalf("CreateMessage failed: %v", err)
	}
	t.Logf("Created message: %s", msg.ID)

	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got %s", msg.Role)
	}

	// 3. Create run
	t.Log("Creating run...")
	run, err := testClient.CreateRunSimple(ctx, thread.ID, testAssistantID, true)
	if err != nil {
		t.Fatalf("CreateRun failed: %v", err)
	}
	t.Logf("Created run: %s (status: %s)", run.ID, run.Status)

	if run.Status != "queued" && run.Status != "in_progress" {
		t.Errorf("Unexpected initial run status: %s", run.Status)
	}

	// 4. Wait for completion
	t.Log("Waiting for run completion...")
	completedRun, err := testClient.WaitForRun(ctx, thread.ID, run.ID)
	if err != nil {
		t.Fatalf("WaitForRun failed: %v", err)
	}
	t.Logf("Run completed with status: %s", completedRun.Status)

	if completedRun.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", completedRun.Status)
	}

	// 5. List messages
	t.Log("Retrieving messages...")
	messages, err := testClient.ListMessages(ctx, thread.ID, 10)
	if err != nil {
		t.Fatalf("ListMessages failed: %v", err)
	}
	t.Logf("Retrieved %d messages", len(messages))

	// Should have at least 2 messages (user + assistant)
	if len(messages) < 2 {
		t.Errorf("Expected at least 2 messages, got %d", len(messages))
	}

	// Print conversation
	t.Log("\n=== Conversation ===")
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if len(msg.Content) > 0 {
			t.Logf("%s: %s", msg.Role, msg.Content[0].Text.Value)
		}
	}

	// Validate assistant response exists
	hasAssistantResponse := false
	for _, msg := range messages {
		if msg.Role == "assistant" && len(msg.Content) > 0 {
			hasAssistantResponse = true
			content := msg.Content[0].Text.Value
			if content == "" {
				t.Error("Assistant response content is empty")
			}
		}
	}
	if !hasAssistantResponse {
		t.Error("No assistant response found in thread")
	}
}

func TestGetThread(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a thread first
	thread, err := testClient.CreateThread(ctx)
	if err != nil {
		t.Fatalf("CreateThread failed: %v", err)
	}

	// Retrieve it
	retrieved, err := testClient.GetThread(ctx, thread.ID)
	if err != nil {
		t.Fatalf("GetThread failed: %v", err)
	}

	if retrieved.ID != thread.ID {
		t.Errorf("Retrieved thread ID mismatch: expected %s, got %s", thread.ID, retrieved.ID)
	}
}

func TestErrorHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test with invalid assistant ID
	_, err := testClient.CreateChatCompletion(ctx, ChatCompletionRequest{
		AssistantID: "invalid-uuid",
		Messages:    []Message{{Role: "user", Content: "test"}},
	})

	if err == nil {
		t.Error("Expected error for invalid assistant ID, got nil")
	}
	t.Logf("Error (expected): %v", err)
}

func TestContextCancellation(t *testing.T) {
	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Wait for context to expire
	time.Sleep(5 * time.Millisecond)

	_, err := testClient.CreateChatCompletion(ctx, ChatCompletionRequest{
		AssistantID: testAssistantID,
		Messages:    []Message{{Role: "user", Content: "test"}},
	})

	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}
	t.Logf("Cancellation error (expected): %v", err)
}
