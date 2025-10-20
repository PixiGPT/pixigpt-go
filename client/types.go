package client

// Message represents a chat message.
type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall represents a function call made by the assistant.
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction represents the function details within a tool call.
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
}

// Tool represents a function tool definition (OpenAI format).
type Tool map[string]interface{}

// ChatCompletionRequest represents a request to the chat completions endpoint.
// AssistantID is optional - if omitted, messages[0] must be a system message.
// Tools can be provided to override assistant's configured tools.
type ChatCompletionRequest struct {
	AssistantID    string    `json:"assistant_id,omitempty"`
	Messages       []Message `json:"messages"`
	Temperature    float32   `json:"temperature,omitempty"`
	MaxTokens      int       `json:"max_tokens,omitempty"`
	EnableThinking *bool     `json:"enable_thinking,omitempty"`
	Tools          []Tool    `json:"tools,omitempty"`
}

// ChatCompletionChoice represents a single choice in the response.
type ChatCompletionChoice struct {
	Index            int    `json:"index"`
	Message          Message `json:"message"`
	FinishReason     string `json:"finish_reason"`
	ReasoningContent string `json:"reasoning_content,omitempty"` // Chain of thought reasoning
}

// ChatCompletionResponse represents the response from chat completions.
type ChatCompletionResponse struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int64                    `json:"created"`
	Model   string                   `json:"model"`
	Choices []ChatCompletionChoice   `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Thread represents a conversation thread.
type Thread struct {
	ID        string            `json:"id"`
	Object    string            `json:"object"`
	CreatedAt int64             `json:"created_at"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// MessageContentText represents text content in a message.
type MessageContentText struct {
	Value       string        `json:"value"`
	Annotations []interface{} `json:"annotations"`
}

// MessageContent represents content in a thread message.
type MessageContent struct {
	Type string              `json:"type"`
	Text MessageContentText `json:"text"`
}

// ThreadMessage represents a message in a thread.
type ThreadMessage struct {
	ID               string           `json:"id"`
	Object           string           `json:"object"`
	CreatedAt        int64            `json:"created_at"`
	ThreadID         string           `json:"thread_id"`
	Role             string           `json:"role"`
	Content          []MessageContent `json:"content"`
	ReasoningContent string           `json:"reasoning_content,omitempty"` // Chain of thought reasoning
}

// Run represents an async run.
type Run struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	CreatedAt   int64  `json:"created_at"`
	ThreadID    string `json:"thread_id"`
	AssistantID string `json:"assistant_id"`
	Status      string `json:"status"` // queued, in_progress, completed, failed
	Model       string `json:"model"`
}

// Assistant represents an AI assistant.
type Assistant struct {
	ID           string  `json:"id"`
	Object       string  `json:"object"`
	CreatedAt    int64   `json:"created_at"`
	Name         string  `json:"name"`
	Instructions string  `json:"instructions"`
	ToolsConfig  *string `json:"tools_config,omitempty"`
}
