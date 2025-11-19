# Zaguan Go SDK Examples

This directory contains examples demonstrating how to use the Zaguan Go SDK.

## Prerequisites

1. A running Zaguan CoreX instance
2. A valid API key
3. Go 1.21 or later

## Setup

Set your environment variables:

```bash
export ZAGUAN_BASE_URL="http://localhost:8080"  # or your Zaguan instance URL
export ZAGUAN_API_KEY="your-api-key-here"
```

## Examples

### Basic Chat (`basic_chat/`)

Demonstrates basic chat completion using the OpenAI-compatible API.

```bash
cd basic_chat
go run main.go
```

**Features shown:**
- Creating a client
- Sending a simple chat request
- Handling responses
- Accessing token usage information

### Anthropic Messages (`anthropic_messages/`)

Demonstrates using Anthropic's native Messages API with extended thinking.

```bash
cd anthropic_messages
go run main.go
```

**Features shown:**
- Using the Messages API
- Enabling extended thinking
- Handling multiple content blocks
- Accessing thinking output
- Cache token information

### Streaming Chat (`streaming_chat/`)

Demonstrates streaming chat completions for real-time responses.

```bash
cd streaming_chat
go run main.go
```

**Features shown:**
- Streaming responses
- Handling SSE events
- Real-time token updates
- Proper cleanup

### Credits Tracking (`credits_tracking/`)

Demonstrates credit balance and usage tracking.

```bash
cd credits_tracking
go run main.go
```

**Features shown:**
- Checking credit balance
- Viewing usage history
- Accessing statistics
- Filtering by date/model/provider

### Provider-Specific Features (`provider_features/`)

Demonstrates using provider-specific parameters.

```bash
cd provider_features
go run main.go
```

**Features shown:**
- Google Gemini reasoning control
- Anthropic extended thinking
- DeepSeek thinking parameter
- Perplexity search options

### Error Handling (`error_handling/`)

Demonstrates proper error handling and recovery.

```bash
cd error_handling
go run main.go
```

**Features shown:**
- Detecting error types
- Handling insufficient credits
- Band access errors
- Rate limiting
- Retry logic

### Multimodal (`multimodal/`)

Demonstrates vision and audio capabilities.

```bash
cd multimodal
go run main.go
```

**Features shown:**
- Sending images (URL and base64)
- Audio input
- Audio output
- Content parts

### Tool Calling (`tool_calling/`)

Demonstrates function/tool calling.

```bash
cd tool_calling
go run main.go
```

**Features shown:**
- Defining tools
- Handling tool calls
- Providing tool results
- Multi-turn conversations

## Running All Examples

```bash
#!/bin/bash
for dir in */; do
    if [ -f "$dir/main.go" ]; then
        echo "Running $dir..."
        (cd "$dir" && go run main.go)
        echo ""
    fi
done
```

## Notes

- All examples use environment variables for configuration
- Examples assume a local Zaguan instance by default
- Error handling is simplified for clarity
- Production code should include more robust error handling

## Troubleshooting

### "ZAGUAN_API_KEY environment variable is required"

Set your API key:
```bash
export ZAGUAN_API_KEY="your-key-here"
```

### "connection refused"

Ensure Zaguan CoreX is running:
```bash
# Check if the server is running
curl http://localhost:8080/v1/models
```

### "insufficient credits"

Check your credit balance using the `credits_tracking` example or upgrade your tier.

## Further Reading

- [SDK Documentation](../README.md)
- [API Endpoints](../API_ENDPOINTS.md)
- [Implementation Plan](../IMPLEMENTATION_PLAN.md)
