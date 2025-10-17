package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the plugin
type Config struct {
	APIKey       string
	Model        string
	Prompt       string
	FilePath     string
	Temperature  float64
	MaxTokens    int
	SystemPrompt string
	OutputFile   string
	Timeout      int
}

// Load creates a new Config from environment variables
func Load() *Config {
	return &Config{
		APIKey:       getEnv("PLUGIN_API_KEY", ""),
		Model:        getEnv("PLUGIN_MODEL", "gpt-4o-mini"),
		Prompt:       getEnv("PLUGIN_PROMPT", ""),
		FilePath:     getEnv("PLUGIN_FILE", ""),
		Temperature:  getEnvFloat("PLUGIN_TEMPERATURE", 0.7),
		MaxTokens:    getEnvInt("PLUGIN_MAX_TOKENS", 1000),
		SystemPrompt: getEnv("PLUGIN_SYSTEM_PROMPT", "You are a helpful assistant."),
		OutputFile:   getEnv("PLUGIN_OUTPUT_FILE", ""),
		Timeout:      getEnvInt("PLUGIN_TIMEOUT", 60),
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API_KEY is required")
	}
	if c.Prompt == "" {
		return fmt.Errorf("PROMPT is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

