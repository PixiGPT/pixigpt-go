package client

// Message represents a chat message.
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID *string    `json:"tool_call_id,omitempty"` // For role="tool" messages
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
	// Attachments from tool execution (Pixi tools only)
	Sources []MessageSource `json:"sources,omitempty"`
	Media   []MessageMedia  `json:"media,omitempty"`
	Code    []MessageCode   `json:"code,omitempty"`
}

// MessageSource represents a source attachment from tools like WebSearch, Fetch
type MessageSource struct {
	ID       string  `json:"id"`
	ToolName string  `json:"tool_name"`
	Title    *string `json:"title,omitempty"`
	URL      *string `json:"url,omitempty"`
	Snippet  *string `json:"snippet,omitempty"`
}

// MessageMedia represents media attachment from DrawImage, EditImage, UserUpload
type MessageMedia struct {
	ID          string  `json:"id"` // ShortID
	Source      string  `json:"source"`
	Type        string  `json:"type"` // image, audio
	Prompt      *string `json:"prompt,omitempty"`
	Description *string `json:"description,omitempty"`
	SignedURL   string  `json:"signed_url"` // 24h temporary R2 signed URL
}

// MessageCode represents code execution result
type MessageCode struct {
	ID              string  `json:"id"`
	Language        string  `json:"language"`
	Code            string  `json:"code"`
	Stdout          *string `json:"stdout,omitempty"`
	Stderr          *string `json:"stderr,omitempty"`
	ExecutionTimeMs *int    `json:"execution_time_ms,omitempty"`
}

// Run represents an async run.
type Run struct {
	ID          string         `json:"id"`
	Object      string         `json:"object"`
	CreatedAt   int64          `json:"created_at"`
	ThreadID    string         `json:"thread_id"`
	AssistantID string         `json:"assistant_id"`
	Status      string         `json:"status"` // queued, in_progress, completed, failed
	Model       string         `json:"model"`
	Message     *ThreadMessage `json:"message,omitempty"` // Populated when completed
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

// VisionUsage represents token usage for vision API calls.
type VisionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// VisionAnalyzeRequest represents a request to analyze an image.
type VisionAnalyzeRequest struct {
	ImageURL   string  `json:"image_url"`
	UserPrompt *string `json:"user_prompt,omitempty"`
}

// VisionAnalyzeResponse represents the response from image analysis.
type VisionAnalyzeResponse struct {
	Result string       `json:"result"`
	Usage  VisionUsage  `json:"usage"`
}

// VisionTagsRequest represents a request to generate tags for an image.
type VisionTagsRequest struct {
	ImageURL string `json:"image_url"`
}

// VisionTagsResponse represents the response from tag generation.
type VisionTagsResponse struct {
	Result string      `json:"result"`
	Usage  VisionUsage `json:"usage"`
}

// VisionOCRRequest represents a request to extract text from an image.
type VisionOCRRequest struct {
	ImageURL string `json:"image_url"`
}

// VisionOCRResponse represents the response from OCR.
type VisionOCRResponse struct {
	Result string      `json:"result"`
	Usage  VisionUsage `json:"usage"`
}

// VisionVideoRequest represents a request to analyze a video.
type VisionVideoRequest struct {
	VideoURL   string  `json:"video_url"`
	UserPrompt *string `json:"user_prompt,omitempty"`
}

// VisionVideoResponse represents the response from video analysis.
type VisionVideoResponse struct {
	Result string      `json:"result"`
	Usage  VisionUsage `json:"usage"`
}

// ModerationTextRequest represents a request to moderate text.
type ModerationTextRequest struct {
	Prompt string `json:"prompt"`
}

// ModerationMediaRequest represents a request to moderate image/video.
type ModerationMediaRequest struct {
	MediaURL string `json:"media_url"`
	IsVideo  bool   `json:"is_video"`
}

// ModerationResponse represents the response from moderation.
type ModerationResponse struct {
	Category string      `json:"category"`
	Score    float64     `json:"score"`
	Usage    VisionUsage `json:"usage"`
}
