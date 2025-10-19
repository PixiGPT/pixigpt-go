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
func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Set defaults
	if req.Temperature == 0 {
		req.Temperature = 0.7
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 2000
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp ChatCompletionResponse
	if err := c.doRequest(ctx, "POST", "/chat/completions", bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
