package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name: "default values",
			envVars: map[string]string{
				"PLUGIN_API_KEY": "test-key",
				"PLUGIN_PROMPT":  "test prompt",
			},
			expected: Config{
				APIKey:       "test-key",
				Model:        "gpt-4o-mini",
				Prompt:       "test prompt",
				FilePath:     "",
				Temperature:  0.7,
				MaxTokens:    1000,
				SystemPrompt: "You are a helpful assistant.",
				OutputFile:   "",
				Timeout:      60,
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"PLUGIN_API_KEY":       "custom-key",
				"PLUGIN_MODEL":         "gpt-4o",
				"PLUGIN_PROMPT":        "custom prompt",
				"PLUGIN_FILE":          "/path/to/file.txt",
				"PLUGIN_TEMPERATURE":   "0.9",
				"PLUGIN_MAX_TOKENS":    "2000",
				"PLUGIN_SYSTEM_PROMPT": "Custom system prompt",
				"PLUGIN_OUTPUT_FILE":   "output.txt",
				"PLUGIN_TIMEOUT":       "120",
			},
			expected: Config{
				APIKey:       "custom-key",
				Model:        "gpt-4o",
				Prompt:       "custom prompt",
				FilePath:     "/path/to/file.txt",
				Temperature:  0.9,
				MaxTokens:    2000,
				SystemPrompt: "Custom system prompt",
				OutputFile:   "output.txt",
				Timeout:      120,
			},
		},
		{
			name: "invalid numeric values fall back to defaults",
			envVars: map[string]string{
				"PLUGIN_API_KEY":     "test-key",
				"PLUGIN_PROMPT":      "test",
				"PLUGIN_TEMPERATURE": "invalid",
				"PLUGIN_MAX_TOKENS":  "not-a-number",
				"PLUGIN_TIMEOUT":     "bad",
			},
			expected: Config{
				APIKey:       "test-key",
				Model:        "gpt-4o-mini",
				Prompt:       "test",
				Temperature:  0.7,
				MaxTokens:    1000,
				SystemPrompt: "You are a helpful assistant.",
				Timeout:      60,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			clearEnv()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer clearEnv()

			// Load config
			cfg := Load()

			// Validate fields
			if cfg.APIKey != tt.expected.APIKey {
				t.Errorf("APIKey = %v, want %v", cfg.APIKey, tt.expected.APIKey)
			}
			if cfg.Model != tt.expected.Model {
				t.Errorf("Model = %v, want %v", cfg.Model, tt.expected.Model)
			}
			if cfg.Prompt != tt.expected.Prompt {
				t.Errorf("Prompt = %v, want %v", cfg.Prompt, tt.expected.Prompt)
			}
			if cfg.FilePath != tt.expected.FilePath {
				t.Errorf("FilePath = %v, want %v", cfg.FilePath, tt.expected.FilePath)
			}
			if cfg.Temperature != tt.expected.Temperature {
				t.Errorf("Temperature = %v, want %v", cfg.Temperature, tt.expected.Temperature)
			}
			if cfg.MaxTokens != tt.expected.MaxTokens {
				t.Errorf("MaxTokens = %v, want %v", cfg.MaxTokens, tt.expected.MaxTokens)
			}
			if cfg.SystemPrompt != tt.expected.SystemPrompt {
				t.Errorf("SystemPrompt = %v, want %v", cfg.SystemPrompt, tt.expected.SystemPrompt)
			}
			if cfg.OutputFile != tt.expected.OutputFile {
				t.Errorf("OutputFile = %v, want %v", cfg.OutputFile, tt.expected.OutputFile)
			}
			if cfg.Timeout != tt.expected.Timeout {
				t.Errorf("Timeout = %v, want %v", cfg.Timeout, tt.expected.Timeout)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				APIKey: "test-key",
				Prompt: "test prompt",
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			config: Config{
				APIKey: "",
				Prompt: "test prompt",
			},
			wantErr: true,
			errMsg:  "API_KEY is required",
		},
		{
			name: "missing prompt",
			config: Config{
				APIKey: "test-key",
				Prompt: "",
			},
			wantErr: true,
			errMsg:  "PROMPT is required",
		},
		{
			name: "missing both",
			config: Config{
				APIKey: "",
				Prompt: "",
			},
			wantErr: true,
			errMsg:  "API_KEY is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("Validate() expected error but got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

// clearEnv clears all PLUGIN_* environment variables
func clearEnv() {
	envVars := []string{
		"PLUGIN_API_KEY",
		"PLUGIN_MODEL",
		"PLUGIN_PROMPT",
		"PLUGIN_FILE",
		"PLUGIN_TEMPERATURE",
		"PLUGIN_MAX_TOKENS",
		"PLUGIN_SYSTEM_PROMPT",
		"PLUGIN_OUTPUT_FILE",
		"PLUGIN_TIMEOUT",
	}
	for _, key := range envVars {
		os.Unsetenv(key)
	}
}
