package zaguansdk

import (
	"testing"
)

func TestValidateAudioTranscriptionRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     AudioTranscriptionRequest
		wantErr bool
	}{
		{
			name: "valid request with file path",
			req: AudioTranscriptionRequest{
				File:  "test.mp3",
				Model: "openai/whisper-1",
			},
			wantErr: false,
		},
		{
			name: "valid request with all parameters",
			req: AudioTranscriptionRequest{
				File:           "test.mp3",
				Model:          "openai/whisper-1",
				Language:       "en",
				Prompt:         "Test prompt",
				ResponseFormat: "json",
				Temperature:    floatPtr(0.5),
			},
			wantErr: false,
		},
		{
			name: "missing file",
			req: AudioTranscriptionRequest{
				Model: "openai/whisper-1",
			},
			wantErr: true,
		},
		{
			name: "missing model",
			req: AudioTranscriptionRequest{
				File: "test.mp3",
			},
			wantErr: true,
		},
		{
			name: "temperature too low",
			req: AudioTranscriptionRequest{
				File:        "test.mp3",
				Model:       "openai/whisper-1",
				Temperature: floatPtr(-0.1),
			},
			wantErr: true,
		},
		{
			name: "temperature too high",
			req: AudioTranscriptionRequest{
				File:        "test.mp3",
				Model:       "openai/whisper-1",
				Temperature: floatPtr(1.1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAudioTranscriptionRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAudioTranscriptionRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAudioTranslationRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     AudioTranslationRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: AudioTranslationRequest{
				File:  "test.mp3",
				Model: "openai/whisper-1",
			},
			wantErr: false,
		},
		{
			name: "valid request with all parameters",
			req: AudioTranslationRequest{
				File:           "test.mp3",
				Model:          "openai/whisper-1",
				Prompt:         "Test prompt",
				ResponseFormat: "text",
				Temperature:    floatPtr(0.3),
			},
			wantErr: false,
		},
		{
			name: "missing file",
			req: AudioTranslationRequest{
				Model: "openai/whisper-1",
			},
			wantErr: true,
		},
		{
			name: "missing model",
			req: AudioTranslationRequest{
				File: "test.mp3",
			},
			wantErr: true,
		},
		{
			name: "temperature out of range",
			req: AudioTranslationRequest{
				File:        "test.mp3",
				Model:       "openai/whisper-1",
				Temperature: floatPtr(2.0),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAudioTranslationRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAudioTranslationRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAudioSpeechRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     AudioSpeechRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: AudioSpeechRequest{
				Model: "openai/tts-1",
				Input: "Hello, world!",
				Voice: "alloy",
			},
			wantErr: false,
		},
		{
			name: "valid request with all parameters",
			req: AudioSpeechRequest{
				Model:          "openai/tts-1-hd",
				Input:          "Test speech",
				Voice:          "nova",
				ResponseFormat: "mp3",
				Speed:          floatPtr(1.5),
			},
			wantErr: false,
		},
		{
			name: "missing model",
			req: AudioSpeechRequest{
				Input: "Hello",
				Voice: "alloy",
			},
			wantErr: true,
		},
		{
			name: "missing input",
			req: AudioSpeechRequest{
				Model: "openai/tts-1",
				Voice: "alloy",
			},
			wantErr: true,
		},
		{
			name: "missing voice",
			req: AudioSpeechRequest{
				Model: "openai/tts-1",
				Input: "Hello",
			},
			wantErr: true,
		},
		{
			name: "speed too low",
			req: AudioSpeechRequest{
				Model: "openai/tts-1",
				Input: "Hello",
				Voice: "alloy",
				Speed: floatPtr(0.2),
			},
			wantErr: true,
		},
		{
			name: "speed too high",
			req: AudioSpeechRequest{
				Model: "openai/tts-1",
				Input: "Hello",
				Voice: "alloy",
				Speed: floatPtr(5.0),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAudioSpeechRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAudioSpeechRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloatPtrToString(t *testing.T) {
	tests := []struct {
		name  string
		input *float64
		want  string
	}{
		{
			name:  "nil pointer",
			input: nil,
			want:  "",
		},
		{
			name:  "zero value",
			input: floatPtr(0.0),
			want:  "0.000000",
		},
		{
			name:  "positive value",
			input: floatPtr(1.5),
			want:  "1.500000",
		},
		{
			name:  "negative value",
			input: floatPtr(-2.3),
			want:  "-2.300000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := floatPtrToString(tt.input)
			if got != tt.want {
				t.Errorf("floatPtrToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function for tests
func floatPtr(f float64) *float64 {
	return &f
}
