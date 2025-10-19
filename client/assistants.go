package client

import (
	"bytes"
	"context"
	"encoding/json"
)

// ListAssistants retrieves all assistants.
func (c *Client) ListAssistants(ctx context.Context) ([]Assistant, error) {
	var resp struct {
		Object string      `json:"object"`
		Data   []Assistant `json:"data"`
	}

	if err := c.doRequest(ctx, "GET", "/assistants", nil, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// GetAssistant retrieves an assistant by ID.
func (c *Client) GetAssistant(ctx context.Context, assistantID string) (*Assistant, error) {
	var assistant Assistant
	if err := c.doRequest(ctx, "GET", "/assistants/"+assistantID, nil, &assistant); err != nil {
		return nil, err
	}
	return &assistant, nil
}

// CreateAssistant creates a new assistant.
func (c *Client) CreateAssistant(ctx context.Context, name, instructions string, toolsConfig *string) (*Assistant, error) {
	reqBody := map[string]interface{}{
		"name":         name,
		"instructions": instructions,
	}
	if toolsConfig != nil {
		reqBody["tools_config"] = *toolsConfig
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	var assistant Assistant
	if err := c.doRequest(ctx, "POST", "/assistants", bytes.NewReader(body), &assistant); err != nil {
		return nil, err
	}

	return &assistant, nil
}

// UpdateAssistant updates an existing assistant.
func (c *Client) UpdateAssistant(ctx context.Context, assistantID, name, instructions string, toolsConfig *string) (*Assistant, error) {
	reqBody := map[string]interface{}{
		"name":         name,
		"instructions": instructions,
	}
	if toolsConfig != nil {
		reqBody["tools_config"] = *toolsConfig
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	var assistant Assistant
	if err := c.doRequest(ctx, "PUT", "/assistants/"+assistantID, bytes.NewReader(body), &assistant); err != nil {
		return nil, err
	}

	return &assistant, nil
}

// DeleteAssistant deletes an assistant.
func (c *Client) DeleteAssistant(ctx context.Context, assistantID string) error {
	return c.doRequest(ctx, "DELETE", "/assistants/"+assistantID, nil, nil)
}

// ListAssistantThreads retrieves all threads used by an assistant.
func (c *Client) ListAssistantThreads(ctx context.Context, assistantID string, limit int) ([]Thread, error) {
	path := fmt.Sprintf("/assistants/%s/threads", assistantID)
	if limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, limit)
	}

	var response struct {
		Object string   `json:"object"`
		Data   []Thread `json:"data"`
	}
	if err := c.doRequest(ctx, "GET", path, nil, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}
