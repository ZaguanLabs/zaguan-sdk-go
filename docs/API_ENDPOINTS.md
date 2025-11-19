# Zaguan CoreX API Endpoints

This document catalogs all API endpoints supported by Zaguan CoreX and how they map to SDK methods.

## OpenAI-Compatible Endpoints

### Chat Completions
- **Endpoint**: `POST /v1/chat/completions`
- **SDK Method**: `client.Chat(ctx, req, opts)`
- **Streaming**: `client.ChatStream(ctx, req, opts)`
- **Description**: Primary chat completion endpoint supporting all OpenAI parameters plus Zaguan extensions

### Models
- **Endpoint**: `GET /v1/models`
- **SDK Method**: `client.ListModels(ctx, opts)`
- **Description**: List all available models across all configured providers

### Embeddings
- **Endpoint**: `POST /v1/embeddings`
- **SDK Method**: `client.CreateEmbeddings(ctx, req, opts)`
- **Description**: Generate text embeddings for semantic search and similarity

### Audio Transcription
- **Endpoint**: `POST /v1/audio/transcriptions`
- **SDK Method**: `client.CreateTranscription(ctx, req, opts)`
- **Description**: Transcribe audio to text (Whisper)

### Audio Translation
- **Endpoint**: `POST /v1/audio/translations`
- **SDK Method**: `client.CreateTranslation(ctx, req, opts)`
- **Description**: Translate audio to English text

### Audio Speech
- **Endpoint**: `POST /v1/audio/speech`
- **SDK Method**: `client.CreateSpeech(ctx, req, opts)`
- **Description**: Text-to-speech synthesis

### Images Generation
- **Endpoint**: `POST /v1/images/generations`
- **SDK Method**: `client.CreateImage(ctx, req, opts)`
- **Description**: Generate images from text prompts (DALL-E)

### Images Edits
- **Endpoint**: `POST /v1/images/edits`
- **SDK Method**: `client.EditImage(ctx, req, opts)`
- **Description**: Edit images with prompts

### Images Variations
- **Endpoint**: `POST /v1/images/variations`
- **SDK Method**: `client.CreateImageVariation(ctx, req, opts)`
- **Description**: Create variations of an image

### Moderations
- **Endpoint**: `POST /v1/moderations`
- **SDK Method**: `client.CreateModeration(ctx, req, opts)`
- **Description**: Content moderation and safety checks

### Files
- **Endpoint**: `POST /v1/files`
- **SDK Method**: `client.UploadFile(ctx, req, opts)`
- **Description**: Upload files for use with assistants

- **Endpoint**: `GET /v1/files`
- **SDK Method**: `client.ListFiles(ctx, opts)`
- **Description**: List uploaded files

- **Endpoint**: `GET /v1/files/{file_id}`
- **SDK Method**: `client.GetFile(ctx, fileID, opts)`
- **Description**: Get file metadata

- **Endpoint**: `DELETE /v1/files/{file_id}`
- **SDK Method**: `client.DeleteFile(ctx, fileID, opts)`
- **Description**: Delete a file

### Fine-tuning
- **Endpoint**: `POST /v1/fine_tuning/jobs`
- **SDK Method**: `client.CreateFineTuningJob(ctx, req, opts)`
- **Description**: Create a fine-tuning job

- **Endpoint**: `GET /v1/fine_tuning/jobs`
- **SDK Method**: `client.ListFineTuningJobs(ctx, opts)`
- **Description**: List fine-tuning jobs

- **Endpoint**: `GET /v1/fine_tuning/jobs/{job_id}`
- **SDK Method**: `client.GetFineTuningJob(ctx, jobID, opts)`
- **Description**: Get fine-tuning job status

- **Endpoint**: `POST /v1/fine_tuning/jobs/{job_id}/cancel`
- **SDK Method**: `client.CancelFineTuningJob(ctx, jobID, opts)`
- **Description**: Cancel a fine-tuning job

### Batches
- **Endpoint**: `POST /v1/batches`
- **SDK Method**: `client.CreateBatch(ctx, req, opts)`
- **Description**: Create a batch processing job

- **Endpoint**: `GET /v1/batches/{batch_id}`
- **SDK Method**: `client.GetBatch(ctx, batchID, opts)`
- **Description**: Get batch status

- **Endpoint**: `POST /v1/batches/{batch_id}/cancel`
- **SDK Method**: `client.CancelBatch(ctx, batchID, opts)`
- **Description**: Cancel a batch

### Assistants
- **Endpoint**: `POST /v1/assistants`
- **SDK Method**: `client.CreateAssistant(ctx, req, opts)`
- **Description**: Create an assistant

- **Endpoint**: `GET /v1/assistants`
- **SDK Method**: `client.ListAssistants(ctx, opts)`
- **Description**: List assistants

- **Endpoint**: `GET /v1/assistants/{assistant_id}`
- **SDK Method**: `client.GetAssistant(ctx, assistantID, opts)`
- **Description**: Get assistant details

- **Endpoint**: `POST /v1/assistants/{assistant_id}`
- **SDK Method**: `client.UpdateAssistant(ctx, assistantID, req, opts)`
- **Description**: Update an assistant

- **Endpoint**: `DELETE /v1/assistants/{assistant_id}`
- **SDK Method**: `client.DeleteAssistant(ctx, assistantID, opts)`
- **Description**: Delete an assistant

### Threads
- **Endpoint**: `POST /v1/threads`
- **SDK Method**: `client.CreateThread(ctx, req, opts)`
- **Description**: Create a conversation thread

- **Endpoint**: `GET /v1/threads/{thread_id}`
- **SDK Method**: `client.GetThread(ctx, threadID, opts)`
- **Description**: Get thread details

- **Endpoint**: `POST /v1/threads/{thread_id}`
- **SDK Method**: `client.UpdateThread(ctx, threadID, req, opts)`
- **Description**: Update a thread

- **Endpoint**: `DELETE /v1/threads/{thread_id}`
- **SDK Method**: `client.DeleteThread(ctx, threadID, opts)`
- **Description**: Delete a thread

### Runs
- **Endpoint**: `POST /v1/threads/{thread_id}/runs`
- **SDK Method**: `client.CreateRun(ctx, threadID, req, opts)`
- **Description**: Create a run on a thread

- **Endpoint**: `GET /v1/threads/{thread_id}/runs/{run_id}`
- **SDK Method**: `client.GetRun(ctx, threadID, runID, opts)`
- **Description**: Get run status

- **Endpoint**: `POST /v1/threads/{thread_id}/runs/{run_id}/cancel`
- **SDK Method**: `client.CancelRun(ctx, threadID, runID, opts)`
- **Description**: Cancel a run

## Anthropic-Native Endpoints

### Messages
- **Endpoint**: `POST /v1/messages`
- **SDK Method**: `client.Messages(ctx, req, opts)`
- **Streaming**: `client.MessagesStream(ctx, req, opts)`
- **Description**: Anthropic's native messages API with extended thinking support

### Count Tokens
- **Endpoint**: `POST /v1/messages/count_tokens`
- **SDK Method**: `client.CountTokens(ctx, req, opts)`
- **Description**: Count tokens for a messages request (Anthropic-specific)

### Batch Messages (Anthropic)
- **Endpoint**: `POST /v1/messages/batches`
- **SDK Method**: `client.CreateMessagesBatch(ctx, req, opts)`
- **Description**: Create a batch of message requests

- **Endpoint**: `GET /v1/messages/batches/{batch_id}`
- **SDK Method**: `client.GetMessagesBatch(ctx, batchID, opts)`
- **Description**: Get batch status

- **Endpoint**: `GET /v1/messages/batches`
- **SDK Method**: `client.ListMessagesBatches(ctx, opts)`
- **Description**: List message batches

- **Endpoint**: `POST /v1/messages/batches/{batch_id}/cancel`
- **SDK Method**: `client.CancelMessagesBatch(ctx, batchID, opts)`
- **Description**: Cancel a message batch

- **Endpoint**: `GET /v1/messages/batches/{batch_id}/results`
- **SDK Method**: `client.GetMessagesBatchResults(ctx, batchID, opts)`
- **Description**: Get batch results (JSONL format)

## Zaguan-Specific Endpoints

### Capabilities
- **Endpoint**: `GET /v1/capabilities`
- **SDK Method**: `client.GetCapabilities(ctx, opts)`
- **Description**: Get detailed capability information for all models

### Credits Balance
- **Endpoint**: `GET /v1/credits/balance`
- **SDK Method**: `client.GetCreditsBalance(ctx, opts)`
- **Description**: Get current credit balance, tier, and bands

### Credits History
- **Endpoint**: `GET /v1/credits/history`
- **SDK Method**: `client.GetCreditsHistory(ctx, opts)`
- **Description**: Get credit usage history with pagination

### Credits Stats
- **Endpoint**: `GET /v1/credits/stats`
- **SDK Method**: `client.GetCreditsStats(ctx, opts)`
- **Description**: Get aggregated credit statistics

### Virtual Models
- **Endpoint**: `GET /v1/virtual-models`
- **SDK Method**: `client.ListVirtualModels(ctx, opts)`
- **Description**: List configured virtual models

- **Endpoint**: `GET /v1/virtual-models/{model_id}`
- **SDK Method**: `client.GetVirtualModel(ctx, modelID, opts)`
- **Description**: Get virtual model configuration

### Provider Status
- **Endpoint**: `GET /v1/providers/status`
- **SDK Method**: `client.GetProviderStatus(ctx, opts)`
- **Description**: Get health status of all configured providers

### Circuit Breaker Status
- **Endpoint**: `GET /v1/circuit-breaker/status`
- **SDK Method**: `client.GetCircuitBreakerStatus(ctx, opts)`
- **Description**: Get circuit breaker status for all providers

## Real-time API (WebSocket)

### OpenAI Realtime
- **Endpoint**: `WS /v1/realtime`
- **SDK Method**: `client.Realtime(ctx, opts)`
- **Description**: WebSocket connection for real-time audio/text streaming

## Admin Endpoints

### Reload Configuration
- **Endpoint**: `POST /v1/admin/reload`
- **SDK Method**: `client.ReloadConfig(ctx, opts)`
- **Description**: Hot-reload server configuration (requires admin privileges)

### Provider Admin
- **Endpoint**: `POST /v1/admin/providers/{provider}/enable`
- **SDK Method**: `client.EnableProvider(ctx, provider, opts)`
- **Description**: Enable a provider at runtime

- **Endpoint**: `POST /v1/admin/providers/{provider}/disable`
- **SDK Method**: `client.DisableProvider(ctx, provider, opts)`
- **Description**: Disable a provider at runtime

## Endpoint Categories

### Priority 1 (Core - Must Implement)
- Chat Completions (OpenAI)
- Messages (Anthropic)
- Models
- Capabilities
- Credits (Balance, History, Stats)

### Priority 2 (Extended - Should Implement)
- Embeddings
- Audio (Transcription, Translation, Speech)
- Images (Generation, Edits, Variations)
- Count Tokens
- Provider Status

### Priority 3 (Advanced - Nice to Have)
- Assistants, Threads, Runs
- Fine-tuning
- Batches (OpenAI and Anthropic)
- Virtual Models
- Circuit Breaker Status
- Admin endpoints

### Priority 4 (Future)
- Real-time API (WebSocket)
- Moderations
- Files
