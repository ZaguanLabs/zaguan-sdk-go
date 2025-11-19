package zaguansdk

import (
	"context"
	"io"
	"testing"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal/testutil"
)

func TestClient_ChatStream(t *testing.T) {
	tests := []struct {
		name    string
		req     ChatRequest
		events  []string
		wantErr bool
	}{
		{
			name: "valid streaming request",
			req: ChatRequest{
				Model: "openai/gpt-4o",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			},
			events: []string{
				testutil.ChatStreamEventFixture("Hello"),
				testutil.ChatStreamEventFixture(" there"),
				testutil.ChatStreamEventFixture("!"),
			},
			wantErr: false,
		},
		{
			name: "invalid request - missing model",
			req: ChatRequest{
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				// For error cases, use default mock server
				mockServer := testutil.NewMockServer(nil)
				defer mockServer.Close()

				client := NewClient(Config{
					BaseURL: mockServer.URL(),
					APIKey:  "test-key",
				})

				_, err := client.ChatStream(context.Background(), tt.req, nil)
				if err == nil {
					t.Error("ChatStream() should have returned validation error")
				}
				return
			}

			// Create streaming mock server
			mockServer := testutil.NewMockServer(
				testutil.StreamingHandler(tt.events),
			)
			defer mockServer.Close()

			client := NewClient(Config{
				BaseURL: mockServer.URL(),
				APIKey:  "test-key",
			})

			stream, err := client.ChatStream(context.Background(), tt.req, nil)
			if err != nil {
				t.Fatalf("ChatStream() error = %v", err)
			}
			defer stream.Close()

			// Read all events
			eventCount := 0
			for {
				event, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("stream.Recv() error = %v", err)
				}
				if event == nil {
					t.Error("stream.Recv() returned nil event")
				}
				eventCount++
			}

			if eventCount != len(tt.events) {
				t.Errorf("received %d events, want %d", eventCount, len(tt.events))
			}
		})
	}
}

func TestClient_MessagesStream(t *testing.T) {
	tests := []struct {
		name    string
		req     MessagesRequest
		events  []string
		wantErr bool
	}{
		{
			name: "valid streaming request",
			req: MessagesRequest{
				Model:     "anthropic/claude-3-5-sonnet-20241022",
				MaxTokens: 1024,
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			events: []string{
				testutil.MessagesStreamEventFixture("Hello"),
				testutil.MessagesStreamEventFixture(" there"),
				`{"type":"message_stop"}`,
			},
			wantErr: false,
		},
		{
			name: "invalid request - missing max_tokens",
			req: MessagesRequest{
				Model: "anthropic/claude-3-5-sonnet-20241022",
				Messages: []AnthropicMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mockServer := testutil.NewMockServer(nil)
				defer mockServer.Close()

				client := NewClient(Config{
					BaseURL: mockServer.URL(),
					APIKey:  "test-key",
				})

				_, err := client.MessagesStream(context.Background(), tt.req, nil)
				if err == nil {
					t.Error("MessagesStream() should have returned validation error")
				}
				return
			}

			mockServer := testutil.NewMockServer(
				testutil.StreamingHandler(tt.events),
			)
			defer mockServer.Close()

			client := NewClient(Config{
				BaseURL: mockServer.URL(),
				APIKey:  "test-key",
			})

			stream, err := client.MessagesStream(context.Background(), tt.req, nil)
			if err != nil {
				t.Fatalf("MessagesStream() error = %v", err)
			}
			defer stream.Close()

			// Read all events
			eventCount := 0
			for {
				event, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("stream.Recv() error = %v", err)
				}
				if event == nil {
					t.Error("stream.Recv() returned nil event")
				}
				eventCount++
			}

			// Should receive all events except the stop event (which triggers EOF)
			expectedCount := len(tt.events) - 1 // message_stop triggers EOF
			if eventCount != expectedCount {
				t.Errorf("received %d events, want %d", eventCount, expectedCount)
			}
		})
	}
}

func TestChatStream_Close(t *testing.T) {
	mockServer := testutil.NewMockServer(
		testutil.StreamingHandler([]string{
			testutil.ChatStreamEventFixture("test"),
		}),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	stream, err := client.ChatStream(context.Background(), ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}, nil)
	if err != nil {
		t.Fatalf("ChatStream() error = %v", err)
	}

	// Close should work
	if err := stream.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Second close should be safe
	if err := stream.Close(); err != nil {
		t.Errorf("Second Close() error = %v", err)
	}

	// Recv after close should error
	_, err = stream.Recv()
	if err == nil {
		t.Error("Recv() after Close() should return error")
	}
}

func TestChatStream_ContextCancellation(t *testing.T) {
	mockServer := testutil.NewMockServer(
		testutil.StreamingHandler([]string{
			testutil.ChatStreamEventFixture("test"),
		}),
	)
	defer mockServer.Close()

	client := NewClient(Config{
		BaseURL: mockServer.URL(),
		APIKey:  "test-key",
	})

	ctx, cancel := context.WithCancel(context.Background())

	stream, err := client.ChatStream(ctx, ChatRequest{
		Model: "openai/gpt-4o",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}, nil)
	if err != nil {
		t.Fatalf("ChatStream() error = %v", err)
	}
	defer stream.Close()

	// Cancel context
	cancel()

	// Recv should return context error
	_, err = stream.Recv()
	if err == nil {
		t.Error("Recv() should return error after context cancellation")
	}
}
