package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// CreateRun starts an async run on a thread.
func (c *Client) CreateRun(ctx context.Context, threadID, assistantID string, enableThinking bool) (*Run, error) {
	reqBody := map[string]interface{}{
		"assistant_id":    assistantID,
		"enable_thinking": enableThinking,
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
