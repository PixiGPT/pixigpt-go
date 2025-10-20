package client

import (
	"context"
)

// CreateThread creates a new conversation thread.
func (c *Client) CreateThread(ctx context.Context) (*Thread, error) {
	var thread Thread
	if err := c.doRequest(ctx, "POST", "/threads", []byte("{}"), &thread); err != nil {
		return nil, err
	}
	return &thread, nil
}

// GetThread retrieves a thread by ID.
func (c *Client) GetThread(ctx context.Context, threadID string) (*Thread, error) {
	var thread Thread
	if err := c.doRequest(ctx, "GET", "/threads/"+threadID, nil, &thread); err != nil {
		return nil, err
	}
	return &thread, nil
}

// ListThreads retrieves all threads for the authenticated user.
func (c *Client) ListThreads(ctx context.Context) ([]Thread, error) {
	var response struct {
		Object string   `json:"object"`
		Data   []Thread `json:"data"`
	}
	if err := c.doRequest(ctx, "GET", "/threads", nil, &response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

// DeleteThread deletes a thread by ID.
func (c *Client) DeleteThread(ctx context.Context, threadID string) error {
	return c.doRequest(ctx, "DELETE", "/threads/"+threadID, nil, nil)
}
