// Package zaguansdk provides audio functionality for the Zaguan SDK.
//
// This file implements the Audio API for:
//   - Transcription: Converting audio to text (Whisper support)
//   - Translation: Translating audio to English
//   - Speech Synthesis: Converting text to spoken audio (TTS)
//
// Supports multiple audio formats including mp3, mp4, mpeg, mpga, m4a, wav, and webm.
package zaguansdk

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ZaguanLabs/zaguan-sdk-go/sdk/internal"
)

// AudioTranscriptionRequest represents a request to transcribe audio.
//
// Supports Whisper and other speech-to-text models.
type AudioTranscriptionRequest struct {
	// File is the audio file to transcribe.
	// Can be a file path (string) or io.Reader.
	// Required.
	File interface{}

	// FileName is the name of the file (required if File is io.Reader).
	FileName string

	// Model is the model identifier to use.
	// Example: "openai/whisper-1"
	// Required.
	Model string

	// Language is the language of the audio (ISO-639-1 format).
	// Optional (improves accuracy and latency).
	Language string

	// Prompt is optional text to guide the model's style.
	// Optional.
	Prompt string

	// ResponseFormat specifies the output format.
	// Values: "json", "text", "srt", "verbose_json", "vtt"
	// Optional (default: "json").
	ResponseFormat string

	// Temperature controls randomness (0.0 - 1.0).
	// Optional.
	Temperature *float64

	// TimestampGranularities specifies timestamp detail level.
	// Values: "word", "segment"
	// Optional.
	TimestampGranularities []string
}

// AudioTranscriptionResponse represents the response from transcription.
type AudioTranscriptionResponse struct {
	// Text is the transcribed text.
	Text string `json:"text"`

	// Language is the detected language (if not specified).
	Language string `json:"language,omitempty"`

	// Duration is the audio duration in seconds.
	Duration float64 `json:"duration,omitempty"`

	// Words contains word-level timestamps (if requested).
	Words []TranscriptionWord `json:"words,omitempty"`

	// Segments contains segment-level timestamps (if requested).
	Segments []TranscriptionSegment `json:"segments,omitempty"`
}

// TranscriptionWord represents a word with timestamp.
type TranscriptionWord struct {
	// Word is the transcribed word.
	Word string `json:"word"`

	// Start is the start time in seconds.
	Start float64 `json:"start"`

	// End is the end time in seconds.
	End float64 `json:"end"`
}

// TranscriptionSegment represents a segment with timestamp.
type TranscriptionSegment struct {
	// ID is the segment identifier.
	ID int `json:"id"`

	// Seek is the seek offset.
	Seek int `json:"seek"`

	// Start is the start time in seconds.
	Start float64 `json:"start"`

	// End is the end time in seconds.
	End float64 `json:"end"`

	// Text is the segment text.
	Text string `json:"text"`

	// Tokens are the token IDs.
	Tokens []int `json:"tokens"`

	// Temperature is the temperature used.
	Temperature float64 `json:"temperature"`

	// AvgLogprob is the average log probability.
	AvgLogprob float64 `json:"avg_logprob"`

	// CompressionRatio is the compression ratio.
	CompressionRatio float64 `json:"compression_ratio"`

	// NoSpeechProb is the no-speech probability.
	NoSpeechProb float64 `json:"no_speech_prob"`
}

// AudioTranslationRequest represents a request to translate audio to English.
//
// Translates audio in any language to English text.
type AudioTranslationRequest struct {
	// File is the audio file to translate.
	// Can be a file path (string) or io.Reader.
	// Required.
	File interface{}

	// FileName is the name of the file (required if File is io.Reader).
	FileName string

	// Model is the model identifier to use.
	// Example: "openai/whisper-1"
	// Required.
	Model string

	// Prompt is optional text to guide the model's style.
	// Optional.
	Prompt string

	// ResponseFormat specifies the output format.
	// Values: "json", "text", "srt", "verbose_json", "vtt"
	// Optional (default: "json").
	ResponseFormat string

	// Temperature controls randomness (0.0 - 1.0).
	// Optional.
	Temperature *float64
}

// AudioTranslationResponse represents the response from translation.
type AudioTranslationResponse struct {
	// Text is the translated text (in English).
	Text string `json:"text"`

	// Duration is the audio duration in seconds.
	Duration float64 `json:"duration,omitempty"`
}

// AudioSpeechRequest represents a request to generate speech from text.
//
// Converts text to spoken audio using TTS models.
type AudioSpeechRequest struct {
	// Model is the TTS model to use.
	// Example: "openai/tts-1", "openai/tts-1-hd"
	// Required.
	Model string `json:"model"`

	// Input is the text to convert to speech.
	// Required.
	Input string `json:"input"`

	// Voice is the voice to use.
	// OpenAI voices: "alloy", "echo", "fable", "onyx", "nova", "shimmer"
	// Required.
	Voice string `json:"voice"`

	// ResponseFormat specifies the audio format.
	// Values: "mp3", "opus", "aac", "flac", "wav", "pcm"
	// Optional (default: "mp3").
	ResponseFormat string `json:"response_format,omitempty"`

	// Speed controls the playback speed (0.25 - 4.0).
	// Optional (default: 1.0).
	Speed *float64 `json:"speed,omitempty"`
}

// CreateTranscription transcribes audio to text.
//
// Supports various audio formats including mp3, mp4, mpeg, mpga, m4a, wav, and webm.
//
// Example:
//
//	resp, err := client.CreateTranscription(ctx, zaguansdk.AudioTranscriptionRequest{
//		File:  "/path/to/audio.mp3",
//		Model: "openai/whisper-1",
//		Language: "en",
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(resp.Text)
func (c *Client) CreateTranscription(ctx context.Context, req AudioTranscriptionRequest, opts *RequestOptions) (*AudioTranscriptionResponse, error) {
	// Validate request
	if err := validateAudioTranscriptionRequest(&req); err != nil {
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "creating audio transcription", "model", req.Model)

	// Create multipart form
	body, contentType, err := createAudioMultipartForm(req.File, req.FileName, map[string]string{
		"model":           req.Model,
		"language":        req.Language,
		"prompt":          req.Prompt,
		"response_format": req.ResponseFormat,
		"temperature":     floatPtrToString(req.Temperature),
	})
	if err != nil {
		return nil, err
	}

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/audio/transcriptions",
		Body:   body,
		Headers: http.Header{
			"Content-Type": []string{contentType},
		},
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			for k, v := range opts.Headers {
				reqCfg.Headers[k] = v
			}
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp AudioTranscriptionResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "create transcription request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "create transcription request succeeded")

	return &resp, nil
}

// CreateTranslation translates audio to English text.
//
// Supports various audio formats. The audio can be in any language,
// and the output will always be in English.
//
// Example:
//
//	resp, err := client.CreateTranslation(ctx, zaguansdk.AudioTranslationRequest{
//		File:  "/path/to/audio.mp3",
//		Model: "openai/whisper-1",
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(resp.Text)
func (c *Client) CreateTranslation(ctx context.Context, req AudioTranslationRequest, opts *RequestOptions) (*AudioTranslationResponse, error) {
	// Validate request
	if err := validateAudioTranslationRequest(&req); err != nil {
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "creating audio translation", "model", req.Model)

	// Create multipart form
	body, contentType, err := createAudioMultipartForm(req.File, req.FileName, map[string]string{
		"model":           req.Model,
		"prompt":          req.Prompt,
		"response_format": req.ResponseFormat,
		"temperature":     floatPtrToString(req.Temperature),
	})
	if err != nil {
		return nil, err
	}

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/audio/translations",
		Body:   body,
		Headers: http.Header{
			"Content-Type": []string{contentType},
		},
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			for k, v := range opts.Headers {
				reqCfg.Headers[k] = v
			}
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	var resp AudioTranslationResponse
	if err := c.internalHTTP.DoJSON(ctx, reqCfg, &resp); err != nil {
		c.log(ctx, LogLevelError, "create translation request failed", "error", err)
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "create translation request succeeded")

	return &resp, nil
}

// CreateSpeech generates audio from text using text-to-speech.
//
// Returns an io.ReadCloser containing the audio data. The caller is
// responsible for closing the reader.
//
// Example:
//
//	audio, err := client.CreateSpeech(ctx, zaguansdk.AudioSpeechRequest{
//		Model: "openai/tts-1",
//		Input: "Hello, world!",
//		Voice: "alloy",
//	}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer audio.Close()
//
//	// Save to file
//	out, _ := os.Create("speech.mp3")
//	defer out.Close()
//	io.Copy(out, audio)
func (c *Client) CreateSpeech(ctx context.Context, req AudioSpeechRequest, opts *RequestOptions) (io.ReadCloser, error) {
	// Validate request
	if err := validateAudioSpeechRequest(&req); err != nil {
		return nil, err
	}

	c.log(ctx, LogLevelDebug, "creating speech", "model", req.Model, "voice", req.Voice)

	// Build request config
	reqCfg := internal.RequestConfig{
		Method: "POST",
		Path:   "/v1/audio/speech",
		Body:   req,
	}

	// Apply request options
	if opts != nil {
		if opts.Timeout > 0 {
			reqCfg.Timeout = opts.Timeout
		}
		if opts.RequestID != "" {
			reqCfg.RequestID = opts.RequestID
		}
		if opts.Headers != nil {
			reqCfg.Headers = opts.Headers
		}
	} else if c.timeout > 0 {
		reqCfg.Timeout = c.timeout
	}

	// Execute request
	resp, err := c.internalHTTP.Do(ctx, reqCfg)
	if err != nil {
		c.log(ctx, LogLevelError, "create speech request failed", "error", err)
		return nil, err
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, internal.ParseErrorResponse(resp)
	}

	c.log(ctx, LogLevelDebug, "create speech request succeeded")

	return resp.Body, nil
}

// createAudioMultipartForm creates a multipart form for audio requests.
func createAudioMultipartForm(file interface{}, fileName string, fields map[string]string) (io.Reader, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file
	var fileReader io.Reader
	var fileNameToUse string

	switch v := file.(type) {
	case string:
		// File path
		f, err := os.Open(v)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		fileReader = f
		fileNameToUse = filepath.Base(v)
	case io.Reader:
		// Reader
		fileReader = v
		if fileName == "" {
			return nil, "", &ValidationError{
				Field:   "file_name",
				Message: "file_name is required when file is io.Reader",
			}
		}
		fileNameToUse = fileName
	default:
		return nil, "", &ValidationError{
			Field:   "file",
			Message: "file must be a string path or io.Reader",
		}
	}

	// Create form file
	part, err := writer.CreateFormFile("file", fileNameToUse)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file data
	if _, err := io.Copy(part, fileReader); err != nil {
		return nil, "", fmt.Errorf("failed to copy file data: %w", err)
	}

	// Add other fields
	for key, value := range fields {
		if value != "" {
			if err := writer.WriteField(key, value); err != nil {
				return nil, "", fmt.Errorf("failed to write field %s: %w", key, err)
			}
		}
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return &buf, writer.FormDataContentType(), nil
}

// floatPtrToString converts a float pointer to string, or returns empty string if nil.
func floatPtrToString(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%f", *f)
}
