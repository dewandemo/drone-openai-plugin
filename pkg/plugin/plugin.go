package plugin

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dewan-ahmed/drone-openai-plugin/internal/config"
	"github.com/dewan-ahmed/drone-openai-plugin/internal/file"
	"github.com/dewan-ahmed/drone-openai-plugin/internal/openai"
	"github.com/dewan-ahmed/drone-openai-plugin/internal/output"
)

// Run executes the plugin workflow
func Run() error {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("starting drone openai plugin")

	// Load configuration from environment
	cfg := config.Load()
	logger.Info("configuration loaded",
		"model", cfg.Model,
		"temperature", cfg.Temperature,
		"max_tokens", cfg.MaxTokens,
		"timeout", cfg.Timeout,
		"has_file", cfg.FilePath != "",
		"has_output_file", cfg.OutputFile != "",
	)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Error("configuration validation failed", "error", err)
		return fmt.Errorf("configuration error: %w", err)
	}

	// Create component instances
	fileProcessor := file.NewProcessor(logger)
	openaiClient := openai.NewClient(cfg.APIKey, logger)
	outputWriter := output.NewWriter(logger)

	// Build messages for OpenAI
	messages := []openai.Message{
		{
			Role:    "system",
			Content: cfg.SystemPrompt,
		},
	}

	// Process user message with optional file
	if cfg.FilePath != "" {
		userMessage, err := fileProcessor.ProcessFileContent(cfg.Prompt, cfg.FilePath)
		if err != nil {
			logger.Error("file processing failed", "error", err)
			return fmt.Errorf("error processing file: %w", err)
		}
		messages = append(messages, userMessage)
	} else {
		messages = append(messages, openai.Message{
			Role:    "user",
			Content: cfg.Prompt,
		})
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	// Call OpenAI API
	logger.Info("calling openai api")
	response, err := openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       cfg.Model,
		Messages:    messages,
		Temperature: cfg.Temperature,
		MaxTokens:   int64(cfg.MaxTokens),
	})
	if err != nil {
		logger.Error("openai api call failed", "error", err)
		return fmt.Errorf("error calling OpenAI: %w", err)
	}

	// Output the response
	if err := outputWriter.WriteResponse(response.Content, response.Usage, cfg.OutputFile); err != nil {
		logger.Error("output writing failed", "error", err)
		return fmt.Errorf("error writing output: %w", err)
	}

	logger.Info("plugin execution completed successfully")
	fmt.Println("\n✓ OpenAI plugin execution completed successfully")
	return nil
}
