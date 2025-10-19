package client

import (
	"bytes"
	"context"
	"encoding/json"
)

// CreateChatCompletion sends a stateless chat completion request.
//
// This is the simplest way to use PixiGPT - no thread management needed.
// The client manages conversation history.
//
// Chain of thought reasoning is automatically extracted from <think> tags
// and placed in the ReasoningContent field of each choice.
//
// Example:
//
//	resp, err := client.CreateChatCompletion(ctx, ChatCompletionRequest{
//	    AssistantID: "e306844d-be73-4cca-ad29-e1255b97b2aa",
//	    Messages: []Message{
//	        {Role: "user", Content: "Hello!"},
//	    },
//	    Temperature: 0.7,
//	    MaxTokens: 2000,
//	})
//	if resp.Choices[0].ReasoningContent != "" {
//	    fmt.Printf("Reasoning: %s\n", resp.Choices[0].ReasoningContent)
//	}
func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Note: Server defaults temperature to 0.6 if 0
	// Note: Server omits max_tokens if 0 (lets vLLM handle it)
	// No client-side defaults needed - pass values as-is

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp ChatCompletionResponse
	if err := c.doRequest(ctx, "POST", "/chat/completions", bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	// Extract reasoning from each choice
	for i := range resp.Choices {
		mainContent, reasoning := extractReasoning(resp.Choices[i].Message.Content)
		resp.Choices[i].Message.Content = mainContent
		resp.Choices[i].ReasoningContent = reasoning
	}

	return &resp, nil
}
