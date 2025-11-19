# Changelog

All notable changes to the Zaguan Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-11-19

### Added
- **Comprehensive input validation** for all request types
  - Temperature, top_p, and penalty range validation
  - Required field validation (model, messages, max_tokens)
  - Enum validation (reasoning_effort, thinking.type)
  - Model ID and config validation
- **Complete test suite** with 221 tests achieving 63% coverage
  - Unit tests for all core functionality
  - Integration tests with mock servers
  - Error handling tests
  - Streaming tests
  - Type definition tests
- **Performance benchmarks** for key operations
- **Test infrastructure** including mock HTTP servers and reusable fixtures
- **Internal package tests** with 73.1% coverage

### Fixed
- **Security**: Fixed all 7 gosec warnings (unhandled Close() errors)
- **Code quality**: Removed all 3 staticcheck warnings (unused code)
  - Removed unused `buildURL()` function
  - Removed unused `lastErr` fields from stream types
- **Error handling**: Properly handle all Close() errors with explicit ignores
- **Race conditions**: Verified zero data races with race detector

### Changed
- Improved error messages with explicit validation feedback
- Enhanced documentation with comprehensive audit reports

### Quality Metrics
- ✅ 0 security issues (was 7)
- ✅ 0 code quality warnings (was 3)
- ✅ 63% test coverage (was 0%)
- ✅ 221 comprehensive tests (was 0)
- ✅ 0 race conditions detected
- ✅ All static analysis clean

## [Unreleased]

### Added - 2025-11-19

#### Core Infrastructure
- HTTP client wrapper with full request/response handling (`internal/http.go`)
- Client struct with HTTP integration (`client.go`)
- Streaming support for both Chat and Messages APIs (`stream.go`)
- Comprehensive error handling with specialized error types
- Request options for per-request customization
- Logger interface for observability
- Request ID generation (UUID v4)
- Context-aware timeout handling
- Query parameter support for filtering and pagination

#### Chat Completions API
- `Chat()` method - Non-streaming chat completions
- `ChatStream()` method - Streaming chat completions with SSE
- Full OpenAI API compatibility
- Zaguan extensions support (routing, provider-specific parameters)
- ChatStreamEvent, ChatStreamChoice, ChatStreamDelta types
- Comprehensive ChatRequest and ChatResponse types
- Tool calling support
- Multimodal (vision) support
- Helper methods: HasReasoningTokens(), HasCachedTokens()

#### Anthropic Messages API
- `Messages()` method - Non-streaming Messages API
- `MessagesStream()` method - Streaming Messages API with SSE
- Native Anthropic format support
- Extended thinking configuration
- MessagesStreamEvent, MessagesStreamDelta types
- Multiple event types (message_start, content_block_delta, etc.)
- Thinking block streaming
- System prompt support
- Cache token tracking

#### Models API
- `ListModels()` method - List all available models
- `GetModel()` method - Get specific model details
- `DeleteModel()` method - Delete fine-tuned models
- Model struct with metadata and permissions
- Provider-prefixed model IDs

#### Capabilities API
- `GetCapabilities()` method - Get all model capabilities
- `GetModelCapabilities()` method - Get specific model capabilities
- Helper methods: SupportsVision(), SupportsTools(), SupportsReasoning()
- ModelCapabilities with detailed feature flags
- Support for both map and array response formats
- Cost information (per 1M tokens)
- Context window and output limit information

#### Credits Tracking API
- `GetCreditsBalance()` method - Get current credit balance
- `GetCreditsHistory()` method - Get usage history with pagination
- `GetCreditsStats()` method - Get aggregated statistics
- Cursor-based pagination support
- Filtering by date, model, provider, band, status
- Period-based grouping (day, week, month, all_time)
- Helper methods: IsLowCredits(), DaysUntilReset(), ParseResetDate()
- CreditsBalance, CreditsHistoryEntry, CreditsStats types
- Provider, Model, Band, and Daily statistics

#### Documentation
- Comprehensive README with quickstart guide
- SDK_OUTLINE.md - Complete design document
- API_ENDPOINTS.md - Endpoint catalog
- IMPLEMENTATION_PLAN.md - Development roadmap
- STATUS.md - Current development status
- IMPLEMENTATION_SUMMARY.md - Detailed implementation summary
- QUICK_REFERENCE.md - Concise API reference
- CHANGELOG.md - This file

#### Examples
- basic_chat - Simple chat completion example
- anthropic_messages - Anthropic Messages with extended thinking
- Examples README with setup instructions

### Technical Details

#### Dependencies
- Go 1.21+
- github.com/google/uuid v1.6.0 (only external dependency)

#### Architecture
- Zero external dependencies (except UUID)
- Context-first design pattern
- Functional options pattern for flexibility
- Interface-based logging
- Structured error types
- Proper resource cleanup
- Iterator pattern for streaming

#### Code Metrics
- ~3,500+ lines of production code
- ~1,500 lines of documentation
- 15+ SDK files
- 100% compilation success
- Examples compile successfully

### Testing Status
- ⏳ Unit tests - Not yet implemented
- ⏳ Integration tests - Not yet implemented
- ✅ Build verification - All code compiles
- ✅ Examples verification - All examples compile

## [0.1.0] - Unreleased

### Status
- Core implementation complete
- Ready for initial testing
- Pending: comprehensive test suite
- Pending: additional examples
- Target: Beta release with 80%+ test coverage

### Breaking Changes
None - Initial release

### Known Issues
- Test coverage at 0% (tests not yet implemented)
- Limited examples (2 of 10+ planned)
- Documentation could be expanded with more GoDoc examples

### Migration Guide
N/A - Initial release

---

## Release Planning

### v0.1.0-beta.1 (Target: 2 weeks)
- [ ] Core Chat and Messages APIs ✅ DONE
- [ ] Streaming support ✅ DONE
- [ ] Basic examples ✅ DONE
- [ ] Initial documentation ✅ DONE
- [ ] Unit tests (pending)
- [ ] Integration tests (optional)

### v0.1.0-beta.2 (Target: 4 weeks)
- [ ] Models & Capabilities ✅ DONE
- [ ] Credits tracking ✅ DONE
- [ ] Extended examples (pending)
- [ ] Test coverage >50% (pending)

### v0.1.0 (Target: 8 weeks)
- [ ] All Priority 1-2 endpoints
- [ ] Test coverage >80%
- [ ] Complete documentation
- [ ] Production-ready

### v1.0.0 (Target: 12 weeks)
- [ ] All endpoints implemented
- [ ] Comprehensive test suite
- [ ] Migration guides
- [ ] Community feedback incorporated

---

**Legend:**
- ✅ Complete
- 🚧 In Progress
- ⏳ Pending
- ❌ Blocked

**Note:** Dates are targets and subject to change based on testing and feedback.
