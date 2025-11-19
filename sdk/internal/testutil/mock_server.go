package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
)

// MockServer is a test HTTP server for mocking API responses.
type MockServer struct {
	Server   *httptest.Server
	Requests []*http.Request
	mu       sync.Mutex
}

// NewMockServer creates a new mock server.
func NewMockServer(handler http.Handler) *MockServer {
	ms := &MockServer{
		Requests: make([]*http.Request, 0),
	}

	if handler == nil {
		handler = ms.DefaultHandler()
	}

	ms.Server = httptest.NewServer(handler)
	return ms
}

// Close closes the mock server.
func (ms *MockServer) Close() {
	ms.Server.Close()
}

// URL returns the base URL of the mock server.
func (ms *MockServer) URL() string {
	return ms.Server.URL
}

// RecordRequest records a request for later inspection.
func (ms *MockServer) RecordRequest(r *http.Request) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.Requests = append(ms.Requests, r)
}

// GetRequests returns all recorded requests.
func (ms *MockServer) GetRequests() []*http.Request {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.Requests
}

// DefaultHandler returns a default handler that returns 404 for all requests.
func (ms *MockServer) DefaultHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ms.RecordRequest(r)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"message": "Not found",
				"type":    "not_found",
			},
		})
	})
}

// ChatCompletionHandler returns a handler for chat completions.
func ChatCompletionHandler(response interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/v1/chat/completions") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// MessagesHandler returns a handler for Anthropic messages.
func MessagesHandler(response interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/v1/messages") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// ErrorHandler returns a handler that returns an error response.
func ErrorHandler(statusCode int, errorType, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"type":    errorType,
				"message": message,
			},
		})
	}
}

// StreamingHandler returns a handler for streaming responses.
func StreamingHandler(events []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		for _, event := range events {
			w.Write([]byte("data: " + event + "\n\n"))
			flusher.Flush()
		}

		w.Write([]byte("data: [DONE]\n\n"))
		flusher.Flush()
	}
}
