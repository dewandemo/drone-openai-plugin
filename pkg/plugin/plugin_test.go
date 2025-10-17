package plugin

import (
	"os"
	"testing"
)

func TestRun_MissingAPIKey(t *testing.T) {
	// Clear environment
	clearPluginEnv()
	defer clearPluginEnv()

	// Set only prompt (missing API key)
	os.Setenv("PLUGIN_PROMPT", "test prompt")

	err := Run()
	if err == nil {
		t.Error("Expected error for missing API key, got nil")
	}
}

func TestRun_MissingPrompt(t *testing.T) {
	// Clear environment
	clearPluginEnv()
	defer clearPluginEnv()

	// Set only API key (missing prompt)
	os.Setenv("PLUGIN_API_KEY", "test-key")

	err := Run()
	if err == nil {
		t.Error("Expected error for missing prompt, got nil")
	}
}

func TestRun_NonexistentFile(t *testing.T) {
	// Clear environment
	clearPluginEnv()
	defer clearPluginEnv()

	// Set required fields with nonexistent file
	os.Setenv("PLUGIN_API_KEY", "test-key")
	os.Setenv("PLUGIN_PROMPT", "test prompt")
	os.Setenv("PLUGIN_FILE", "/nonexistent/file.txt")

	err := Run()
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

// Note: We cannot easily test the full Run() function with actual OpenAI API calls
// in unit tests without mocking. The tests above verify the basic validation logic.
// For full integration testing with the OpenAI API, we would need:
// 1. A valid API key
// 2. Network access
// 3. Mocking infrastructure or dependency injection
//
// In a production environment, you would typically:
// - Use interfaces for the OpenAI client to enable mocking
// - Create integration tests that run separately from unit tests
// - Use test doubles or stubs for external dependencies

func TestRun_ValidConfigButInvalidAPIKey(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=true to run.")
	}

	// Clear environment
	clearPluginEnv()
	defer clearPluginEnv()

	// Set config with invalid API key
	os.Setenv("PLUGIN_API_KEY", "invalid-key-12345")
	os.Setenv("PLUGIN_PROMPT", "Say hello")
	os.Setenv("PLUGIN_MODEL", "gpt-4o-mini")

	err := Run()
	if err == nil {
		t.Error("Expected error for invalid API key, got nil")
	}
}

// clearPluginEnv clears all PLUGIN_* environment variables
func clearPluginEnv() {
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

// Example of how to structure a test with a real API key (for manual testing only)
// This test is skipped by default and requires explicit environment variable to run
func TestRun_WithRealAPIKey(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY_FOR_TESTING") == "" {
		t.Skip("Skipping test with real API. Set OPENAI_API_KEY_FOR_TESTING to run.")
	}

	// Clear environment
	clearPluginEnv()
	defer clearPluginEnv()

	// Set valid configuration
	os.Setenv("PLUGIN_API_KEY", os.Getenv("OPENAI_API_KEY_FOR_TESTING"))
	os.Setenv("PLUGIN_PROMPT", "Say 'Hello from test' in exactly those words")
	os.Setenv("PLUGIN_MODEL", "gpt-4o-mini")
	os.Setenv("PLUGIN_MAX_TOKENS", "20")

	err := Run()
	if err != nil {
		t.Errorf("Run() unexpected error with valid API key: %v", err)
	}
}
