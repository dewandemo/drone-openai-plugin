package file

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dewan-ahmed/drone-openai-plugin/internal/openai"
)

// Processor handles file operations for the plugin
type Processor struct {
	logger *slog.Logger
}

// NewProcessor creates a new file processor
func NewProcessor(logger *slog.Logger) *Processor {
	return &Processor{
		logger: logger,
	}
}

// ProcessFileContent reads and processes a file, returning a message with the file content
func (p *Processor) ProcessFileContent(prompt, filePath string) (openai.Message, error) {
	p.logger.Info("processing file", "path", filePath)
	
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return openai.Message{}, fmt.Errorf("error reading file: %w", err)
	}

	// Check if it's an image file
	if p.isImageFile(filePath) {
		p.logger.Info("detected image file", "type", p.getMimeType(filePath))
		return p.createImageMessage(prompt, filePath, fileData), nil
	}

	// For text files, append content to prompt
	p.logger.Info("detected text file", "size_bytes", len(fileData))
	fileContent := string(fileData)
	combinedPrompt := fmt.Sprintf("%s\n\nFile content:\n%s", prompt, fileContent)
	
	return openai.Message{
		Role:    "user",
		Content: combinedPrompt,
	}, nil
}

// createImageMessage creates a multimodal message for image files
func (p *Processor) createImageMessage(prompt, filePath string, fileData []byte) openai.Message {
	mimeType := p.getMimeType(filePath)
	base64Image := base64.StdEncoding.EncodeToString(fileData)
	dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image)

	return openai.Message{
		Role: "user",
		Content: []openai.MessagePart{
			{
				Type: "text",
				Text: prompt,
			},
			{
				Type: "image_url",
				ImageURL: &openai.ImageURL{
					URL: dataURL,
				},
			},
		},
	}
}

// isImageFile checks if the file is an image based on extension
func (p *Processor) isImageFile(filePath string) bool {
	lower := strings.ToLower(filePath)
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, ext := range imageExts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// getMimeType returns the MIME type for the file
func (p *Processor) getMimeType(filePath string) string {
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
