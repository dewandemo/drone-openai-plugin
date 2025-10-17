package openai

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Client wraps the official OpenAI SDK client
type Client struct {
	client *openai.Client
	logger *slog.Logger
}

// NewClient creates a new OpenAI client using the official SDK
func NewClient(apiKey string, logger *slog.Logger) *Client {
	oaiClient := openai.NewClient(option.WithAPIKey(apiKey))
	return &Client{
		client: &oaiClient,
		logger: logger,
	}
}

// Message represents a chat message
type Message struct {
	Role    string
	Content interface{} // Can be string or []MessagePart for multimodal
}

// MessagePart represents a part of a multimodal message
type MessagePart struct {
	Type     string // "text" or "image_url"
	Text     string
	ImageURL *ImageURL
}

// ImageURL represents an image URL in a message
type ImageURL struct {
	URL string
}

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Model       string
	Messages    []Message
	Temperature float64
	MaxTokens   int64
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64
}

// ChatCompletionResponse represents the response from OpenAI
type ChatCompletionResponse struct {
	Content string
	Usage   Usage
}

// CreateChatCompletion sends a request to OpenAI and returns the response
func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	c.logger.Info("creating chat completion",
		"model", req.Model,
		"temperature", req.Temperature,
		"max_tokens", req.MaxTokens,
		"num_messages", len(req.Messages),
	)

	// Convert our messages to OpenAI SDK format
	messages := make([]openai.ChatCompletionMessageParamUnion, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = convertMessage(msg)
	}

	// Create the request using the official SDK
	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    req.Model,
	}

	if req.Temperature > 0 {
		params.Temperature = openai.Float(req.Temperature)
	}
	if req.MaxTokens > 0 {
		params.MaxTokens = openai.Int(req.MaxTokens)
	}

	// Make the API call
	resp, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		c.logger.Error("OpenAI API call failed", "error", err)
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	// Extract the response content
	if len(resp.Choices) == 0 {
		c.logger.Error("no choices in response")
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content
	if content == "" {
		c.logger.Warn("empty content in response")
		return nil, fmt.Errorf("empty response from OpenAI")
	}

	c.logger.Info("chat completion successful",
		"prompt_tokens", resp.Usage.PromptTokens,
		"completion_tokens", resp.Usage.CompletionTokens,
		"total_tokens", resp.Usage.TotalTokens,
	)

	return &ChatCompletionResponse{
		Content: content,
		Usage: Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

// convertMessage converts our internal Message type to OpenAI SDK format
func convertMessage(msg Message) openai.ChatCompletionMessageParamUnion {
	switch msg.Role {
	case "system":
		return openai.SystemMessage(msg.Content.(string))
	case "user":
		// Check if it's multimodal content
		if parts, ok := msg.Content.([]MessagePart); ok {
			// Build multimodal content
			contentParts := make([]openai.ChatCompletionContentPartUnionParam, len(parts))
			for i, part := range parts {
				if part.Type == "text" {
					contentParts[i] = openai.TextContentPart(part.Text)
				} else if part.Type == "image_url" && part.ImageURL != nil {
					contentParts[i] = openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
						URL: part.ImageURL.URL,
					})
				}
			}
			return openai.UserMessage(contentParts)
		}
		// Simple text message
		return openai.UserMessage(msg.Content.(string))
	case "assistant":
		return openai.AssistantMessage(msg.Content.(string))
	default:
		return openai.UserMessage(msg.Content.(string))
	}
}
