package zaguansdk

import (
	"testing"
)

func TestUsage_HasReasoningTokens(t *testing.T) {
	tests := []struct {
		name  string
		usage Usage
		want  bool
	}{
		{
			name: "has reasoning tokens",
			usage: Usage{
				CompletionTokensDetails: &TokenDetails{
					ReasoningTokens: 100,
				},
			},
			want: true,
		},
		{
			name: "no reasoning tokens",
			usage: Usage{
				CompletionTokensDetails: &TokenDetails{
					ReasoningTokens: 0,
				},
			},
			want: false,
		},
		{
			name:  "nil details",
			usage: Usage{},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.usage.HasReasoningTokens(); got != tt.want {
				t.Errorf("HasReasoningTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsage_HasCachedTokens(t *testing.T) {
	tests := []struct {
		name  string
		usage Usage
		want  bool
	}{
		{
			name: "has cached tokens",
			usage: Usage{
				PromptTokensDetails: &TokenDetails{
					CachedTokens: 50,
				},
			},
			want: true,
		},
		{
			name: "no cached tokens",
			usage: Usage{
				PromptTokensDetails: &TokenDetails{
					CachedTokens: 0,
				},
			},
			want: false,
		},
		{
			name:  "nil details",
			usage: Usage{},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.usage.HasCachedTokens(); got != tt.want {
				t.Errorf("HasCachedTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_Types(t *testing.T) {
	// Test basic message creation
	msg := Message{
		Role:    "user",
		Content: "Hello, world!",
	}

	if msg.Role != "user" {
		t.Errorf("Message.Role = %v, want user", msg.Role)
	}
	if msg.Content != "Hello, world!" {
		t.Errorf("Message.Content = %v, want Hello, world!", msg.Content)
	}
}

func TestMessage_WithToolCalls(t *testing.T) {
	msg := Message{
		Role: "assistant",
		ToolCalls: []ToolCall{
			{
				ID:   "call_123",
				Type: "function",
				Function: FunctionCall{
					Name:      "get_weather",
					Arguments: `{"location": "San Francisco"}`,
				},
			},
		},
	}

	if len(msg.ToolCalls) != 1 {
		t.Errorf("len(ToolCalls) = %d, want 1", len(msg.ToolCalls))
	}
	if msg.ToolCalls[0].ID != "call_123" {
		t.Errorf("ToolCall.ID = %v, want call_123", msg.ToolCalls[0].ID)
	}
}

func TestContentPart_Types(t *testing.T) {
	tests := []struct {
		name string
		part ContentPart
		want string
	}{
		{
			name: "text content",
			part: ContentPart{
				Type: "text",
				Text: "Hello",
			},
			want: "text",
		},
		{
			name: "image content",
			part: ContentPart{
				Type: "image_url",
				ImageURL: &ImageURL{
					URL: "https://example.com/image.jpg",
				},
			},
			want: "image_url",
		},
		{
			name: "audio content",
			part: ContentPart{
				Type: "input_audio",
				InputAudio: &InputAudio{
					Data:   "base64data",
					Format: "wav",
				},
			},
			want: "input_audio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.part.Type != tt.want {
				t.Errorf("ContentPart.Type = %v, want %v", tt.part.Type, tt.want)
			}
		})
	}
}

func TestTool_Creation(t *testing.T) {
	tool := Tool{
		Type: "function",
		Function: FunctionDefinition{
			Name:        "get_weather",
			Description: "Get the current weather",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
	}

	if tool.Type != "function" {
		t.Errorf("Tool.Type = %v, want function", tool.Type)
	}
	if tool.Function.Name != "get_weather" {
		t.Errorf("Tool.Function.Name = %v, want get_weather", tool.Function.Name)
	}
}

func TestChoice_Types(t *testing.T) {
	choice := Choice{
		Index: 0,
		Message: &Message{
			Role:    "assistant",
			Content: "Hello!",
		},
		FinishReason: "stop",
	}

	if choice.Index != 0 {
		t.Errorf("Choice.Index = %d, want 0", choice.Index)
	}
	if choice.FinishReason != "stop" {
		t.Errorf("Choice.FinishReason = %v, want stop", choice.FinishReason)
	}
}

func TestChatResponse_Structure(t *testing.T) {
	resp := ChatResponse{
		ID:      "chatcmpl-123",
		Object:  "chat.completion",
		Created: 1677652288,
		Model:   "openai/gpt-4o",
		Choices: []Choice{
			{
				Index: 0,
				Message: &Message{
					Role:    "assistant",
					Content: "Hello!",
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
	}

	if resp.ID != "chatcmpl-123" {
		t.Errorf("ChatResponse.ID = %v, want chatcmpl-123", resp.ID)
	}
	if len(resp.Choices) != 1 {
		t.Errorf("len(Choices) = %d, want 1", len(resp.Choices))
	}
	if resp.Usage.TotalTokens != 15 {
		t.Errorf("Usage.TotalTokens = %d, want 15", resp.Usage.TotalTokens)
	}
}

func TestAnthropicMessage_Types(t *testing.T) {
	msg := AnthropicMessage{
		Role:    "user",
		Content: "Hello, Claude!",
	}

	if msg.Role != "user" {
		t.Errorf("AnthropicMessage.Role = %v, want user", msg.Role)
	}
}

func TestAnthropicThinkingConfig(t *testing.T) {
	config := AnthropicThinkingConfig{
		Type:         "enabled",
		BudgetTokens: 5000,
	}

	if config.Type != "enabled" {
		t.Errorf("ThinkingConfig.Type = %v, want enabled", config.Type)
	}
	if config.BudgetTokens != 5000 {
		t.Errorf("ThinkingConfig.BudgetTokens = %d, want 5000", config.BudgetTokens)
	}
}

func TestMessagesResponse_Structure(t *testing.T) {
	resp := MessagesResponse{
		ID:   "msg_123",
		Type: "message",
		Role: "assistant",
		Content: []AnthropicContentBlock{
			{
				Type: "text",
				Text: "Hello!",
			},
		},
		Model:      "anthropic/claude-3-5-sonnet-20241022",
		StopReason: "end_turn",
		Usage: AnthropicUsage{
			InputTokens:  10,
			OutputTokens: 5,
		},
	}

	if resp.ID != "msg_123" {
		t.Errorf("MessagesResponse.ID = %v, want msg_123", resp.ID)
	}
	if len(resp.Content) != 1 {
		t.Errorf("len(Content) = %d, want 1", len(resp.Content))
	}
	if resp.Usage.InputTokens != 10 {
		t.Errorf("Usage.InputTokens = %d, want 10", resp.Usage.InputTokens)
	}
}

func TestAnthropicContentBlock_Types(t *testing.T) {
	tests := []struct {
		name  string
		block AnthropicContentBlock
		want  string
	}{
		{
			name: "text block",
			block: AnthropicContentBlock{
				Type: "text",
				Text: "Hello",
			},
			want: "text",
		},
		{
			name: "thinking block",
			block: AnthropicContentBlock{
				Type:     "thinking",
				Thinking: "Let me think...",
			},
			want: "thinking",
		},
		{
			name: "tool use block",
			block: AnthropicContentBlock{
				Type:  "tool_use",
				ID:    "tool_123",
				Name:  "get_weather",
				Input: map[string]interface{}{"location": "SF"},
			},
			want: "tool_use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.block.Type != tt.want {
				t.Errorf("ContentBlock.Type = %v, want %v", tt.block.Type, tt.want)
			}
		})
	}
}

func TestModelCapabilities_Structure(t *testing.T) {
	cap := ModelCapabilities{
		ModelID:           "openai/gpt-4o",
		Provider:          "openai",
		SupportsVision:    true,
		SupportsTools:     true,
		SupportsReasoning: false,
		MaxContextTokens:  128000,
		MaxOutputTokens:   4096,
		InputCostPer1M:    5.0,
		OutputCostPer1M:   15.0,
		Features:          []string{"json_mode", "structured_outputs"},
		Modalities:        []string{"text", "image"},
	}

	if cap.ModelID != "openai/gpt-4o" {
		t.Errorf("ModelID = %v, want openai/gpt-4o", cap.ModelID)
	}
	if !cap.SupportsVision {
		t.Error("SupportsVision should be true")
	}
	if len(cap.Features) != 2 {
		t.Errorf("len(Features) = %d, want 2", len(cap.Features))
	}
}

func TestCreditsHistoryEntry_Structure(t *testing.T) {
	entry := CreditsHistoryEntry{
		ID:               "entry_123",
		Timestamp:        "2025-11-19T12:00:00Z",
		RequestID:        "req_123",
		Model:            "openai/gpt-4o",
		Provider:         "openai",
		Band:             "A",
		PromptTokens:     10,
		CompletionTokens: 20,
		TotalTokens:      30,
		CreditsDebited:   5,
		Status:           "success",
	}

	if entry.Model != "openai/gpt-4o" {
		t.Errorf("Model = %v, want openai/gpt-4o", entry.Model)
	}
	if entry.CreditsDebited != 5 {
		t.Errorf("CreditsDebited = %d, want 5", entry.CreditsDebited)
	}
}
