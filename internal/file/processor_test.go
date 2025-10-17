package file

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dewan-ahmed/drone-openai-plugin/internal/openai"
)

func TestProcessFileContent_TextFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	processor := NewProcessor(logger)

	// Create a temporary text file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "This is test content"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	prompt := "Analyze this file"
	msg, err := processor.ProcessFileContent(prompt, testFile)
	if err != nil {
		t.Fatalf("ProcessFileContent() error = %v", err)
	}

	if msg.Role != "user" {
		t.Errorf("Expected role user, got %v", msg.Role)
	}

	expectedContent := prompt + "\n\nFile content:\n" + testContent
	contentStr, ok := msg.Content.(string)
	if !ok {
		t.Errorf("Expected string content for text file")
	} else if contentStr != expectedContent {
		t.Errorf("Content mismatch.\nExpected: %q\nGot: %q", expectedContent, contentStr)
	}
}

func TestProcessFileContent_ImageFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	processor := NewProcessor(logger)

	tests := []struct {
		name     string
		filename string
		mimeType string
	}{
		{"PNG image", "test.png", "image/png"},
		{"JPEG image", "test.jpg", "image/jpeg"},
		{"JPEG image (alternative)", "test.jpeg", "image/jpeg"},
		{"GIF image", "test.gif", "image/gif"},
		{"WebP image", "test.webp", "image/webp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, tt.filename)
			testContent := []byte("fake image data")
			if err := os.WriteFile(testFile, testContent, 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			prompt := "Describe this image"
			msg, err := processor.ProcessFileContent(prompt, testFile)
			if err != nil {
				t.Fatalf("ProcessFileContent() error = %v", err)
			}

			if msg.Role != "user" {
				t.Errorf("Expected role user, got %v", msg.Role)
			}

			parts, ok := msg.Content.([]openai.MessagePart)
			if !ok {
				t.Fatalf("Expected []MessagePart content for image file")
			}

			if len(parts) != 2 {
				t.Fatalf("Expected 2 content parts, got %d", len(parts))
			}

			// Check text part
			textPart := parts[0]
			if textPart.Type != "text" {
				t.Errorf("First part type should be text, got %v", textPart.Type)
			}
			if textPart.Text != prompt {
				t.Errorf("Text part = %q, want %q", textPart.Text, prompt)
			}

			// Check image part
			imagePart := parts[1]
			if imagePart.Type != "image_url" {
				t.Errorf("Second part type should be image_url, got %v", imagePart.Type)
			}
			if imagePart.ImageURL == nil {
				t.Fatal("ImageURL should not be nil")
			}
			if !strings.HasPrefix(imagePart.ImageURL.URL, "data:"+tt.mimeType+";base64,") {
				t.Errorf("ImageURL should start with 'data:%s;base64,', got %s", tt.mimeType, imagePart.ImageURL.URL[:50])
			}
		})
	}
}

func TestProcessFileContent_FileNotFound(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	processor := NewProcessor(logger)

	_, err := processor.ProcessFileContent("test", "/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestIsImageFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	processor := NewProcessor(logger)

	tests := []struct {
		filename string
		isImage  bool
	}{
		{"test.jpg", true},
		{"test.JPG", true},
		{"test.jpeg", true},
		{"test.JPEG", true},
		{"test.png", true},
		{"test.PNG", true},
		{"test.gif", true},
		{"test.GIF", true},
		{"test.webp", true},
		{"test.WEBP", true},
		{"test.txt", false},
		{"test.go", false},
		{"test.pdf", false},
		{"test.doc", false},
		{"image.bmp", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := processor.isImageFile(tt.filename)
			if result != tt.isImage {
				t.Errorf("isImageFile(%q) = %v, want %v", tt.filename, result, tt.isImage)
			}
		})
	}
}

func TestGetMimeType(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	processor := NewProcessor(logger)

	tests := []struct {
		filename string
		mimeType string
	}{
		{"test.jpg", "image/jpeg"},
		{"test.JPG", "image/jpeg"},
		{"test.jpeg", "image/jpeg"},
		{"test.JPEG", "image/jpeg"},
		{"test.png", "image/png"},
		{"test.PNG", "image/png"},
		{"test.gif", "image/gif"},
		{"test.GIF", "image/gif"},
		{"test.webp", "image/webp"},
		{"test.WEBP", "image/webp"},
		{"test.unknown", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := processor.getMimeType(tt.filename)
			if result != tt.mimeType {
				t.Errorf("getMimeType(%q) = %q, want %q", tt.filename, result, tt.mimeType)
			}
		})
	}
}
