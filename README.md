# PixiGPT Go Client

Production-grade Go client for the [PixiGPT API](https://pixigpt.com).

## Features

- 🚀 **High Performance**: Connection pooling, keep-alive, optimized for high volume
- 🔄 **Smart Retries**: Exponential backoff with configurable retry logic
- ⏱️ **Proper Timeouts**: Prevents hanging requests
- 🎯 **Context Support**: Full `context.Context` integration for cancellation
- 📦 **Minimal Dependencies**: Standard library + optional godotenv for examples
- 🔧 **OpenAI Compatible**: Familiar API surface for OpenAI users

## Installation

```bash
go get github.com/PixiGPT/pixigpt-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/PixiGPT/pixigpt-go/client"
)

func main() {
    // Create client
    c := client.New("sk-proj-YOUR_API_KEY", "https://pixigpt.com/v1")

    // Send chat completion
    resp, err := c.CreateChatCompletion(context.Background(), client.ChatCompletionRequest{
        AssistantID: "your-assistant-id",  // Optional - omit for pure OpenAI mode
        Messages: []client.Message{
            {Role: "user", Content: "Hello!"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)

    // Access chain of thought reasoning if available
    if resp.Choices[0].ReasoningContent != "" {
        fmt.Printf("Thinking: %s\n", resp.Choices[0].ReasoningContent)
    }
}
```

## Configuration

The client uses production-grade defaults optimized for high volume:

- **Connection Pooling**: 100 max idle connections, 10 per host
- **Timeouts**: 30s client timeout, 10s dial, 5s TLS handshake
- **Keep-Alive**: Enabled for connection reuse
- **Retries**: Up to 3 attempts with exponential backoff

### Custom Configuration

```go
import (
    "net/http"
    "time"
)

// Custom HTTP client with different timeouts
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}

c := client.New(apiKey, baseURL,
    client.WithHTTPClient(httpClient),
    client.WithRetryMax(5),
)
```

## API Methods

### Chat Completions (Stateless)

Simplest method - no thread management needed:

```go
resp, err := c.CreateChatCompletion(ctx, client.ChatCompletionRequest{
    AssistantID: assistantID,  // Optional - omit to provide your own system prompt
    Messages: []client.Message{
        {Role: "user", Content: "What's the weather?"},
    },
    Temperature: 0.7,
    MaxTokens: 2000,
})
```

**With Tool Calling:**

```go
resp, err := c.CreateChatCompletion(ctx, client.ChatCompletionRequest{
    AssistantID: assistantID,
    Messages: []client.Message{
        {Role: "user", Content: "What's the weather in Paris?"},
    },
    Tools: []client.Tool{
        {
            "type": "function",
            "function": map[string]interface{}{
                "name": "get_weather",
                "description": "Get current weather for a location",
                "parameters": map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "location": map[string]interface{}{"type": "string"},
                    },
                },
            },
        },
    },
})

// Check for tool calls
if resp.Choices[0].FinishReason == "tool_calls" {
    for _, toolCall := range resp.Choices[0].Message.ToolCalls {
        fmt.Printf("Tool: %s(%s)\n", toolCall.Function.Name, toolCall.Function.Arguments)
        // Execute tool and send result back...
    }
}
```

**Pure OpenAI Mode (No Assistant):**

```go
resp, err := c.CreateChatCompletion(ctx, client.ChatCompletionRequest{
    // No AssistantID - provide your own system prompt
    Messages: []client.Message{
        {Role: "system", Content: "You are a helpful assistant."},
        {Role: "user", Content: "Hello!"},
    },
})
```

### Threads (Async with Memory)

For multi-turn conversations with persistent memory:

```go
// 1. Create thread
thread, err := c.CreateThread(ctx)

// 2. Add messages
msg, err := c.CreateMessage(ctx, thread.ID, "user", "Hello!")

// 3. Run assistant
run, err := c.CreateRun(ctx, thread.ID, assistantID, true)

// 4. Wait for completion
completedRun, err := c.WaitForRun(ctx, thread.ID, run.ID)

// 5. Get messages
messages, err := c.ListMessages(ctx, thread.ID, 10)
```

### Assistants

Manage AI assistants:

```go
// List
assistants, err := c.ListAssistants(ctx)

// Create
assistant, err := c.CreateAssistant(ctx, "My Assistant", "You are helpful.", nil)

// Update
assistant, err := c.UpdateAssistant(ctx, id, "Updated Name", "New instructions", nil)

// Delete
err := c.DeleteAssistant(ctx, assistantID)
```

## Error Handling

The client provides typed errors for common cases:

```go
resp, err := c.CreateChatCompletion(ctx, req)
if err != nil {
    if client.IsAuthError(err) {
        log.Fatal("Invalid API key")
    }
    if client.IsRateLimitError(err) {
        log.Fatal("Rate limit exceeded")
    }
    log.Fatal(err)
}
```

## Examples

See the [examples/](examples/) directory for complete working examples:

- [`chat.go`](examples/chat.go) - Simple chat completion
- [`thread.go`](examples/thread.go) - Multi-turn conversation with threads
- [`assistant.go`](examples/assistant.go) - Assistant management

To run examples:

```bash
# Copy .env.example to .env and add your API key
cp .env.example .env

# Run chat example
cd examples
go run chat.go
```

## Production Considerations

### Connection Pooling

The client reuses HTTP connections automatically. For very high volume (1000+ req/s), consider tuning:

```go
transport := &http.Transport{
    MaxIdleConns:        500,  // Increase pool size
    MaxIdleConnsPerHost: 50,   // More per-host connections
}

httpClient := &http.Client{Transport: transport}
c := client.New(apiKey, baseURL, client.WithHTTPClient(httpClient))
```

### Context Cancellation

Always use context with timeout for production:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := c.CreateChatCompletion(ctx, req)
```

### Retry Strategy

By default, the client retries 5xx errors (server issues) but not 4xx errors (client errors). Network errors are retried with exponential backoff.

To disable retries:

```go
c := client.New(apiKey, baseURL, client.WithRetryMax(0))
```

## License

MIT

## Contributing

Pull requests welcome! This is open source - keep it simple, fast, and production-ready.
