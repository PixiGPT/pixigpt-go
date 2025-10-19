package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CreateRun starts an async run on a thread.
//
// Temperature and MaxTokens are optional - if 0, server uses defaults.
func (c *Client) CreateRun(ctx context.Context, threadID, assistantID string, temperature float32, maxTokens int, enableThinking bool) (*Run, error) {
	reqBody := map[string]interface{}{
		"assistant_id":    assistantID,
		"enable_thinking": enableThinking,
	}

	// Only include if > 0 (server handles defaults)
	if temperature > 0 {
		reqBody["temperature"] = temperature
	}
	if maxTokens > 0 {
		reqBody["max_tokens"] = maxTokens
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	var run Run
	path := fmt.Sprintf("/threads/%s/runs", threadID)
	if err := c.doRequest(ctx, "POST", path, bytes.NewReader(body), &run); err != nil {
		return nil, err
	}

	return &run, nil
}

// GetRun retrieves run status.
func (c *Client) GetRun(ctx context.Context, threadID, runID string) (*Run, error) {
	var run Run
	path := fmt.Sprintf("/threads/%s/runs/%s", threadID, runID)
	if err := c.doRequest(ctx, "GET", path, nil, &run); err != nil {
		return nil, err
	}
	return &run, nil
}

// CreateRunSimple starts an async run with defaults (no temp/max_tokens).
func (c *Client) CreateRunSimple(ctx context.Context, threadID, assistantID string, enableThinking bool) (*Run, error) {
	return c.CreateRun(ctx, threadID, assistantID, 0, 0, enableThinking)
}

// WaitForRun polls until run completes (or fails).
//
// Returns the completed run or error. Context can be used to cancel polling.
func (c *Client) WaitForRun(ctx context.Context, threadID, runID string) (*Run, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			run, err := c.GetRun(ctx, threadID, runID)
			if err != nil {
				return nil, err
			}

			switch run.Status {
			case "completed":
				return run, nil
			case "failed":
				return run, fmt.Errorf("run failed")
			case "cancelled":
				return run, fmt.Errorf("run cancelled")
			// Continue polling for "queued" or "in_progress"
			}
		}
	}
}
