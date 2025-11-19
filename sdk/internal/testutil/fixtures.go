package testutil

// ChatCompletionFixture returns a sample chat completion response.
func ChatCompletionFixture() map[string]interface{} {
	return map[string]interface{}{
		"id":      "chatcmpl-123",
		"object":  "chat.completion",
		"created": 1677652288,
		"model":   "openai/gpt-4o-mini",
		"choices": []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": "Hello! How can I help you today?",
				},
				"finish_reason": "stop",
			},
		},
		"usage": map[string]interface{}{
			"prompt_tokens":     10,
			"completion_tokens": 9,
			"total_tokens":      19,
		},
	}
}

// MessagesFixture returns a sample Anthropic messages response.
func MessagesFixture() map[string]interface{} {
	return map[string]interface{}{
		"id":   "msg_123",
		"type": "message",
		"role": "assistant",
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "Hello! How can I help you today?",
			},
		},
		"model":       "anthropic/claude-3-5-sonnet-20241022",
		"stop_reason": "end_turn",
		"usage": map[string]interface{}{
			"input_tokens":  10,
			"output_tokens": 9,
		},
	}
}

// ModelsFixture returns a sample models list response.
func ModelsFixture() map[string]interface{} {
	return map[string]interface{}{
		"object": "list",
		"data": []map[string]interface{}{
			{
				"id":       "openai/gpt-4o",
				"object":   "model",
				"created":  1677652288,
				"owned_by": "openai",
			},
			{
				"id":       "anthropic/claude-3-5-sonnet-20241022",
				"object":   "model",
				"created":  1677652288,
				"owned_by": "anthropic",
			},
		},
	}
}

// CreditsBalanceFixture returns a sample credits balance response.
func CreditsBalanceFixture() map[string]interface{} {
	return map[string]interface{}{
		"credits_remaining": 1000,
		"credits_total":     2000,
		"credits_used":      1000,
		"credits_percent":   50.0,
		"tier":              "pro",
		"bands":             []string{"A", "B", "C"},
		"reset_date":        "2025-12-01T00:00:00Z",
	}
}

// ChatStreamEventFixture returns a sample streaming event.
func ChatStreamEventFixture(content string) string {
	return `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"openai/gpt-4o-mini","choices":[{"index":0,"delta":{"content":"` + content + `"},"finish_reason":null}]}`
}

// MessagesStreamEventFixture returns a sample Anthropic streaming event.
func MessagesStreamEventFixture(content string) string {
	return `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"` + content + `"}}`
}
