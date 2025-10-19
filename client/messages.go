package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// CreateMessage adds a message to a thread.
func (c *Client) CreateMessage(ctx context.Context, threadID string, role, content string) (*ThreadMessage, error) {
	reqBody := map[string]string{
		"role":    role,
		"content": content,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	var msg ThreadMessage
	if err := c.doRequest(ctx, "POST", "/threads/"+threadID+"/messages", bytes.NewReader(body), &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// ListMessages retrieves messages from a thread.
func (c *Client) ListMessages(ctx context.Context, threadID string, limit int) ([]ThreadMessage, error) {
	if limit == 0 {
		limit = 20
	}

	var resp struct {
		Object  string          `json:"object"`
		Data    []ThreadMessage `json:"data"`
		FirstID string          `json:"first_id,omitempty"`
		LastID  string          `json:"last_id,omitempty"`
		HasMore bool            `json:"has_more"`
	}

	path := fmt.Sprintf("/threads/%s/messages?limit=%d", threadID, limit)
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
