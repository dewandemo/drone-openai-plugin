package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	APIKey      string
	Model       string
	Prompt      string
	FilePath    string
	Temperature float64
	MaxTokens   int
	SystemPrompt string
	OutputFile  string
	Timeout     int
}

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func main() {
	config := loadConfig()
	
	if err := validateConfig(config); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Timeout)*time.Second)
	defer cancel()

	response, err := callOpenAI(ctx, config)
	if err != nil {
		log.Fatalf("Error calling OpenAI: %v", err)
	}

	if err := outputResult(response, config); err != nil {
		log.Fatalf("Error outputting result: %v", err)
	}

	fmt.Println("OpenAI plugin execution completed successfully")
}

func loadConfig() *Config {
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

func validateConfig(config *Config) error {
	if config.APIKey == "" {
		return fmt.Errorf("API_KEY is required")
	}
	if config.Prompt == "" {
		return fmt.Errorf("PROMPT is required")
	}
	return nil
}

func callOpenAI(ctx context.Context, config *Config) (*OpenAIResponse, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: config.SystemPrompt,
		},
	}

	// Handle user message with optional file
	if config.FilePath != "" {
		content, err := createMessageWithFile(config.Prompt, config.FilePath)
		if err != nil {
			return nil, fmt.Errorf("error processing file: %v", err)
		}
		messages = append(messages, Message{
			Role:    "user",
			Content: content,
		})
	} else {
		messages = append(messages, Message{
			Role:    "user",
			Content: config.Prompt,
		})
	}

	request := OpenAIRequest{
		Model:       config.Model,
		Messages:    messages,
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &openAIResp, nil
}

func createMessageWithFile(prompt, filePath string) (interface{}, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Check if it's an image file
	if isImageFile(filePath) {
		mimeType := getMimeType(filePath)
		base64Image := base64.StdEncoding.EncodeToString(fileData)
		
		return []ContentPart{
			{
				Type: "text",
				Text: prompt,
			},
			{
				Type: "image_url",
				ImageURL: &ImageURL{
					URL: fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image),
				},
			},
		}, nil
	}

	// For text files, append content to prompt
	fileContent := string(fileData)
	combinedPrompt := fmt.Sprintf("%s\n\nFile content:\n%s", prompt, fileContent)
	return combinedPrompt, nil
}

func isImageFile(filePath string) bool {
	lower := strings.ToLower(filePath)
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, ext := range imageExts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func getMimeType(filePath string) string {
	lower := strings.ToLower(filePath)
	switch {
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lower, ".png"):
		return "image/png"
	case strings.HasSuffix(lower, ".gif"):
		return "image/gif"
	case strings.HasSuffix(lower, ".webp"):
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func outputResult(response *OpenAIResponse, config *Config) error {
	if len(response.Choices) == 0 {
		return fmt.Errorf("no response from OpenAI")
	}

	// Get the response content
	var content string
	switch v := response.Choices[0].Message.Content.(type) {
	case string:
		content = v
	default:
		// Handle complex content types
		jsonContent, _ := json.Marshal(v)
		content = string(jsonContent)
	}

	// Print to stdout
	fmt.Println("=== OpenAI Response ===")
	fmt.Println(content)
	fmt.Printf("\n=== Usage Statistics ===\n")
	fmt.Printf("Prompt Tokens: %d\n", response.Usage.PromptTokens)
	fmt.Printf("Completion Tokens: %d\n", response.Usage.CompletionTokens)
	fmt.Printf("Total Tokens: %d\n", response.Usage.TotalTokens)

	// Save to file if specified
	if config.OutputFile != "" {
		if err := os.WriteFile(config.OutputFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("error writing output file: %v", err)
		}
		fmt.Printf("\nResponse saved to: %s\n", config.OutputFile)
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
		var intVal int
		if _, err := fmt.Sscanf(value, "%d", &intVal); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		var floatVal float64
		if _, err := fmt.Sscanf(value, "%f", &floatVal); err == nil {
			return floatVal
		}
	}
	return defaultValue
}