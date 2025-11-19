# Zaguan Go SDK Implementation Plan

This document outlines the phased implementation plan for the Zaguan Go SDK.

## Phase 1: Foundation (Week 1)

### 1.1 Project Setup
- [x] Create directory structure
- [x] Create `version.go` with SDK version
- [x] Create comprehensive `doc.go` package documentation
- [ ] Initialize Go module (`go.mod`)
- [ ] Set up GitHub repository structure
- [ ] Create `LICENSE` and `README.md`

### 1.2 Core Types
**Files**: `client.go`, `option.go`, `errors.go`

- [ ] Implement `Config` struct with validation
- [ ] Implement `Client` struct with HTTP client management
- [ ] Implement `NewClient()` constructor
- [ ] Implement `RequestOptions` for per-request overrides
- [ ] Implement `APIError` with proper error wrapping
- [ ] Add user-agent header with SDK version
- [ ] Add request ID generation (UUID v4)

### 1.3 Chat Completions (OpenAI)
**Files**: `chat.go`, `stream.go`

- [ ] Define `ChatRequest` with all OpenAI fields
- [ ] Define `ChatResponse` with usage details
- [ ] Define `Message`, `Choice`, `Usage`, `TokenDetails`
- [ ] Implement `Chat()` method (non-streaming)
- [ ] Implement `ChatStream()` method with SSE parsing
- [ ] Implement `ChatStream.Recv()` iterator pattern
- [ ] Implement `ChatStream.Close()` cleanup
- [ ] Handle context cancellation properly

### 1.4 Testing Infrastructure
**Files**: `client_test.go`, `chat_test.go`

- [ ] Set up `httptest.Server` for mocking
- [ ] Write unit tests for request marshaling
- [ ] Write unit tests for response unmarshaling
- [ ] Write tests for error handling
- [ ] Write tests for streaming
- [ ] Add integration test suite (optional, requires live API)

## Phase 2: Anthropic Native API (Week 2)

### 2.1 Messages API
**Files**: `messages.go`, `messages_stream.go`

- [ ] Define `MessagesRequest` struct
- [ ] Define `MessagesResponse` struct
- [ ] Define `AnthropicMessage`, `AnthropicContentBlock`
- [ ] Define `AnthropicThinkingConfig` for extended thinking
- [ ] Implement `Messages()` method (non-streaming)
- [ ] Implement `MessagesStream()` method
- [ ] Implement `MessagesStream.Recv()` iterator
- [ ] Handle thinking blocks properly

### 2.2 Token Counting
**Files**: `messages.go`

- [ ] Implement `CountTokens()` method
- [ ] Define `CountTokensResponse` struct
- [ ] Add proper error handling

### 2.3 Anthropic Batches
**Files**: `messages_batch.go`

- [ ] Define `MessagesBatchRequest` struct
- [ ] Define `MessagesBatchResponse` struct
- [ ] Implement `CreateMessagesBatch()`
- [ ] Implement `GetMessagesBatch()`
- [ ] Implement `ListMessagesBatches()`
- [ ] Implement `CancelMessagesBatch()`
- [ ] Implement `GetMessagesBatchResults()` (returns JSONL)

### 2.4 Testing
**Files**: `messages_test.go`

- [ ] Unit tests for Messages API
- [ ] Unit tests for streaming
- [ ] Unit tests for thinking configuration
- [ ] Unit tests for token counting
- [ ] Unit tests for batch operations

## Phase 3: Models & Capabilities (Week 3)

### 3.1 Models
**Files**: `models.go`

- [ ] Define `Model` struct
- [ ] Define `ModelsResponse` struct
- [ ] Implement `ListModels()`
- [ ] Implement `GetModel()` (if needed)
- [ ] Add filtering/pagination support

### 3.2 Capabilities
**Files**: `capabilities.go`

- [ ] Define `ModelCapabilities` struct
- [ ] Define `CapabilitiesResponse` struct
- [ ] Implement `GetCapabilities()`
- [ ] Add helper methods (e.g., `SupportsVision()`)

### 3.3 Testing
**Files**: `models_test.go`, `capabilities_test.go`

- [ ] Unit tests for model listing
- [ ] Unit tests for capabilities
- [ ] Mock responses for all providers

## Phase 4: Credits & Usage Tracking (Week 3)

### 4.1 Credits Balance
**Files**: `credits.go`

- [ ] Define `CreditsBalance` struct
- [ ] Implement `GetCreditsBalance()`
- [ ] Add tier and band information

### 4.2 Credits History
**Files**: `credits.go`

- [ ] Define `CreditsHistoryEntry` struct
- [ ] Define `CreditsHistoryResponse` with pagination
- [ ] Implement `GetCreditsHistory()`
- [ ] Add filtering by date range, model, provider

### 4.3 Credits Stats
**Files**: `credits.go`

- [ ] Define `CreditsStats` struct
- [ ] Implement `GetCreditsStats()`
- [ ] Add aggregation by period, band, provider

### 4.4 Testing
**Files**: `credits_test.go`

- [ ] Unit tests for balance retrieval
- [ ] Unit tests for history pagination
- [ ] Unit tests for stats aggregation

## Phase 5: Extended OpenAI Features (Week 4)

### 5.1 Embeddings
**Files**: `embeddings.go`

- [ ] Define `EmbeddingsRequest` struct
- [ ] Define `EmbeddingsResponse` struct
- [ ] Implement `CreateEmbeddings()`

### 5.2 Audio
**Files**: `audio.go`

- [ ] Define `TranscriptionRequest` struct
- [ ] Define `TranslationRequest` struct
- [ ] Define `SpeechRequest` struct
- [ ] Implement `CreateTranscription()`
- [ ] Implement `CreateTranslation()`
- [ ] Implement `CreateSpeech()`
- [ ] Handle multipart form data for audio uploads

### 5.3 Images
**Files**: `images.go`

- [ ] Define `ImageRequest` struct
- [ ] Define `ImageResponse` struct
- [ ] Implement `CreateImage()`
- [ ] Implement `EditImage()`
- [ ] Implement `CreateImageVariation()`
- [ ] Handle multipart form data for image uploads

### 5.4 Testing
**Files**: `embeddings_test.go`, `audio_test.go`, `images_test.go`

- [ ] Unit tests for all new endpoints
- [ ] Mock multipart uploads

## Phase 6: Provider-Specific Helpers (Week 5)

### 6.1 Helper Functions
**Files**: `helpers.go`

- [ ] `WithGoogleReasoning(effort, budget)` helper
- [ ] `WithAnthropicThinking(enabled, budget)` helper
- [ ] `WithPerplexitySearch(domains, citations)` helper
- [ ] `WithDeepSeekThinking(enabled)` helper
- [ ] `WithOpenAIFlex(enabled, priority)` helper

### 6.2 Provider Constants
**Files**: `providers.go`

- [ ] Define provider name constants
- [ ] Define model name constants (common ones)
- [ ] Define reasoning effort levels
- [ ] Define audio voices and formats

### 6.3 Documentation
**Files**: `PROVIDERS.md`

- [ ] Document each provider's specific features
- [ ] Provide usage examples for each provider
- [ ] Document reasoning token behavior per provider

## Phase 7: Advanced Features (Week 6)

### 7.1 Virtual Models
**Files**: `virtual_models.go`

- [ ] Define `VirtualModel` struct
- [ ] Implement `ListVirtualModels()`
- [ ] Implement `GetVirtualModel()`

### 7.2 Provider Status
**Files**: `provider_status.go`

- [ ] Define `ProviderStatus` struct
- [ ] Implement `GetProviderStatus()`
- [ ] Implement `GetCircuitBreakerStatus()`

### 7.3 Batches (OpenAI)
**Files**: `batches.go`

- [ ] Define `BatchRequest` struct
- [ ] Define `BatchResponse` struct
- [ ] Implement `CreateBatch()`
- [ ] Implement `GetBatch()`
- [ ] Implement `CancelBatch()`

### 7.4 Admin Operations
**Files**: `admin.go`

- [ ] Implement `ReloadConfig()`
- [ ] Implement `EnableProvider()`
- [ ] Implement `DisableProvider()`

## Phase 8: Polish & Documentation (Week 7)

### 8.1 Examples
**Directory**: `examples/`

- [ ] Basic chat example
- [ ] Streaming chat example
- [ ] Anthropic messages example
- [ ] Provider-specific features example
- [ ] Credits tracking example
- [ ] Error handling example
- [ ] Multimodal (vision) example
- [ ] Tool calling example

### 8.2 Documentation
**Files**: Various markdown files

- [ ] Complete README.md with quickstart
- [ ] API reference documentation
- [ ] Migration guide from OpenAI SDK
- [ ] Best practices guide
- [ ] Troubleshooting guide
- [ ] CHANGELOG.md

### 8.3 Code Quality
- [ ] Run `go vet` and fix issues
- [ ] Run `golangci-lint` and fix issues
- [ ] Ensure 80%+ test coverage
- [ ] Add GoDoc comments to all public APIs
- [ ] Review and improve error messages

### 8.4 CI/CD
**Files**: `.github/workflows/`

- [ ] Set up GitHub Actions for tests
- [ ] Set up code coverage reporting
- [ ] Set up automated releases
- [ ] Set up documentation generation

## Phase 9: Beta Release (Week 8)

### 9.1 Pre-release Checklist
- [ ] All Priority 1 & 2 endpoints implemented
- [ ] All tests passing
- [ ] Documentation complete
- [ ] Examples working
- [ ] Security review completed
- [ ] Performance benchmarks run

### 9.2 Release Tasks
- [ ] Tag v0.1.0-beta.1
- [ ] Publish to GitHub
- [ ] Announce in community channels
- [ ] Gather feedback

### 9.3 Beta Iteration
- [ ] Fix reported bugs
- [ ] Improve documentation based on feedback
- [ ] Add missing features
- [ ] Optimize performance

## Phase 10: Stable Release (Week 10)

### 10.1 Stabilization
- [ ] Address all beta feedback
- [ ] Finalize API surface
- [ ] Complete all Priority 1-3 features
- [ ] Achieve 90%+ test coverage
- [ ] Performance optimization

### 10.2 Release v1.0.0
- [ ] Tag v1.0.0
- [ ] Publish to pkg.go.dev
- [ ] Update all documentation
- [ ] Create release announcement
- [ ] Submit to awesome-go lists

## Success Metrics

### Code Quality
- Test coverage: >80% (target: 90%)
- Zero critical security issues
- All public APIs documented
- Passes all linters

### Performance
- Chat completion latency: <50ms overhead
- Streaming: <10ms first token overhead
- Memory: <10MB baseline usage
- Zero memory leaks

### Documentation
- README with quickstart
- 10+ working examples
- Complete API reference
- Migration guides

### Community
- 100+ GitHub stars (6 months)
- 10+ contributors
- Active issue tracking
- Regular releases

## Dependencies

### Required
- `encoding/json` (stdlib)
- `net/http` (stdlib)
- `context` (stdlib)
- `io` (stdlib)

### Optional
- `github.com/google/uuid` - UUID generation
- Testing framework (stdlib `testing` sufficient)

### No External Dependencies for Core
The SDK should minimize external dependencies to reduce supply chain risk and simplify adoption.
