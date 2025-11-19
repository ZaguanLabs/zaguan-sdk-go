package zaguansdk

import (
	"testing"
)

func TestChatRequest_BasicFields(t *testing.T) {
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	if req.Model != "openai/gpt-4o" {
		t.Errorf("Model = %v, want openai/gpt-4o", req.Model)
	}
	if len(req.Messages) != 1 {
		t.Errorf("len(Messages) = %d, want 1", len(req.Messages))
	}
}

func TestChatRequest_WithTemperature(t *testing.T) {
	temp := float32(0.7)
	req := ChatRequest{
		Model:       "openai/gpt-4o",
		Temperature: &temp,
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	if req.Temperature == nil {
		t.Error("Temperature should not be nil")
	}
	if *req.Temperature != 0.7 {
		t.Errorf("Temperature = %v, want 0.7", *req.Temperature)
	}
}

func TestChatRequest_WithMaxTokens(t *testing.T) {
	maxTokens := 1000
	req := ChatRequest{
		Model:     "openai/gpt-4o",
		MaxTokens: &maxTokens,
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	if req.MaxTokens == nil {
		t.Error("MaxTokens should not be nil")
	}
	if *req.MaxTokens != 1000 {
		t.Errorf("MaxTokens = %v, want 1000", *req.MaxTokens)
	}
}

func TestChatRequest_WithTools(t *testing.T) {
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "What's the weather?"},
		},
		Tools: []Tool{
			{
				Type: "function",
				Function: FunctionDefinition{
					Name:        "get_weather",
					Description: "Get current weather",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
			},
		},
	}

	if len(req.Tools) != 1 {
		t.Errorf("len(Tools) = %d, want 1", len(req.Tools))
	}
	if req.Tools[0].Function.Name != "get_weather" {
		t.Errorf("Tool name = %v, want get_weather", req.Tools[0].Function.Name)
	}
}

func TestChatRequest_WithStop(t *testing.T) {
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Stop: []string{"\n", "END"},
	}

	if len(req.Stop) != 2 {
		t.Errorf("len(Stop) = %d, want 2", len(req.Stop))
	}
}

func TestChatRequest_WithPenalties(t *testing.T) {
	presence := float32(0.5)
	frequency := float32(0.3)
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		PresencePenalty:  &presence,
		FrequencyPenalty: &frequency,
	}

	if req.PresencePenalty == nil || *req.PresencePenalty != 0.5 {
		t.Error("PresencePenalty not set correctly")
	}
	if req.FrequencyPenalty == nil || *req.FrequencyPenalty != 0.3 {
		t.Error("FrequencyPenalty not set correctly")
	}
}

func TestChatRequest_WithLogitBias(t *testing.T) {
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		LogitBias: map[string]float32{
			"50256": -100,
			"50257": 100,
		},
	}

	if len(req.LogitBias) != 2 {
		t.Errorf("len(LogitBias) = %d, want 2", len(req.LogitBias))
	}
}

func TestChatRequest_WithUser(t *testing.T) {
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		User: "user-123",
	}

	if req.User != "user-123" {
		t.Errorf("User = %v, want user-123", req.User)
	}
}

func TestChatRequest_WithReasoningEffort(t *testing.T) {
	req := ChatRequest{
		Model: "openai/o1",
		Messages: []Message{
			{Role: "user", Content: "Solve this"},
		},
		ReasoningEffort: "high",
	}

	if req.ReasoningEffort != "high" {
		t.Errorf("ReasoningEffort = %v, want high", req.ReasoningEffort)
	}
}

func TestChatRequest_WithAudio(t *testing.T) {
	req := ChatRequest{
		Model: "openai/gpt-4o-audio",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Audio: &AudioConfig{
			Voice:  "alloy",
			Format: "mp3",
		},
	}

	if req.Audio == nil {
		t.Error("Audio should not be nil")
	}
	if req.Audio.Voice != "alloy" {
		t.Errorf("Audio.Voice = %v, want alloy", req.Audio.Voice)
	}
}

func TestChatRequest_WithModalities(t *testing.T) {
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		Modalities: []string{"text", "audio"},
	}

	if len(req.Modalities) != 2 {
		t.Errorf("len(Modalities) = %d, want 2", len(req.Modalities))
	}
}

func TestChatRequest_ZaguanExtensions(t *testing.T) {
	req := ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
		ProviderOptions: map[string]interface{}{
			"custom_param": "value",
		},
		VirtualModelID: "my-model-alias",
		Metadata: map[string]interface{}{
			"user_id": "123",
		},
	}

	if len(req.ProviderOptions) != 1 {
		t.Error("ProviderOptions not set")
	}
	if req.VirtualModelID != "my-model-alias" {
		t.Error("VirtualModelID not set")
	}
	if len(req.Metadata) != 1 {
		t.Error("Metadata not set")
	}
}

func TestMessage_WithName(t *testing.T) {
	msg := Message{
		Role:    "user",
		Name:    "Alice",
		Content: "Hello",
	}

	if msg.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", msg.Name)
	}
}

func TestMessage_WithToolCallID(t *testing.T) {
	msg := Message{
		Role:       "tool",
		ToolCallID: "call_123",
		Content:    "Result",
	}

	if msg.ToolCallID != "call_123" {
		t.Errorf("ToolCallID = %v, want call_123", msg.ToolCallID)
	}
}

func TestMessage_MultimodalContent(t *testing.T) {
	msg := Message{
		Role: "user",
		Content: []ContentPart{
			{
				Type: "text",
				Text: "What's in this image?",
			},
			{
				Type: "image_url",
				ImageURL: &ImageURL{
					URL:    "https://example.com/image.jpg",
					Detail: "high",
				},
			},
		},
	}

	if parts, ok := msg.Content.([]ContentPart); ok {
		if len(parts) != 2 {
			t.Errorf("len(Content) = %d, want 2", len(parts))
		}
	} else {
		t.Error("Content should be []ContentPart")
	}
}

func TestImageURL_Detail(t *testing.T) {
	img := ImageURL{
		URL:    "https://example.com/image.jpg",
		Detail: "low",
	}

	if img.Detail != "low" {
		t.Errorf("Detail = %v, want low", img.Detail)
	}
}

func TestInputAudio_Format(t *testing.T) {
	audio := InputAudio{
		Data:   "base64data",
		Format: "wav",
	}

	if audio.Format != "wav" {
		t.Errorf("Format = %v, want wav", audio.Format)
	}
}

func TestAudioConfig_AllFields(t *testing.T) {
	config := AudioConfig{
		Voice:  "nova",
		Format: "opus",
	}

	if config.Voice != "nova" {
		t.Errorf("Voice = %v, want nova", config.Voice)
	}
	if config.Format != "opus" {
		t.Errorf("Format = %v, want opus", config.Format)
	}
}

func TestFunctionDefinition_WithStrict(t *testing.T) {
	strict := true
	fn := FunctionDefinition{
		Name:        "get_weather",
		Description: "Get weather",
		Parameters:  map[string]interface{}{"type": "object"},
		Strict:      &strict,
	}

	if fn.Strict == nil || !*fn.Strict {
		t.Error("Strict should be true")
	}
}

func TestToolCall_Structure(t *testing.T) {
	call := ToolCall{
		ID:   "call_123",
		Type: "function",
		Function: FunctionCall{
			Name:      "get_weather",
			Arguments: `{"location": "SF"}`,
		},
	}

	if call.ID != "call_123" {
		t.Errorf("ID = %v, want call_123", call.ID)
	}
	if call.Function.Name != "get_weather" {
		t.Error("Function name not set")
	}
}

func TestChatResponse_WithSystemFingerprint(t *testing.T) {
	resp := ChatResponse{
		ID:                "chatcmpl-123",
		Object:            "chat.completion",
		Created:           1677652288,
		Model:             "openai/gpt-4o",
		SystemFingerprint: "fp_123",
		Choices: []Choice{
			{
				Index:        0,
				Message:      &Message{Role: "assistant", Content: "Hello"},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
	}

	if resp.SystemFingerprint != "fp_123" {
		t.Errorf("SystemFingerprint = %v, want fp_123", resp.SystemFingerprint)
	}
}

func TestTokenDetails_AllFields(t *testing.T) {
	details := TokenDetails{
		ReasoningTokens:          100,
		CachedTokens:             50,
		AudioTokens:              25,
		AcceptedPredictionTokens: 10,
		RejectedPredictionTokens: 5,
	}

	if details.ReasoningTokens != 100 {
		t.Error("ReasoningTokens not set")
	}
	if details.CachedTokens != 50 {
		t.Error("CachedTokens not set")
	}
	if details.AudioTokens != 25 {
		t.Error("AudioTokens not set")
	}
}

func TestChoice_WithDelta(t *testing.T) {
	choice := Choice{
		Index: 0,
		Delta: &Message{
			Role:    "assistant",
			Content: "Hello",
		},
		FinishReason: "stop",
	}

	if choice.Delta == nil {
		t.Error("Delta should not be nil")
	}
	if choice.Delta.Role != "assistant" {
		t.Error("Delta role not set")
	}
}
