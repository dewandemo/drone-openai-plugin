# Drone OpenAI Plugin

A Drone CI plugin that integrates with OpenAI's API to process prompts with optional file attachments and output AI-generated responses.

## Features

- Send custom prompts to OpenAI models (GPT-4, GPT-3.5, etc.)
- Attach files (text or images) for context-aware processing
- Configurable model parameters (temperature, max tokens)
- Save responses to files
- Support for system prompts to control AI behavior
- Automatic handling of different file types (text and images)

## Usage

### Basic Example

```yaml
steps:
  - name: openai-task
    image: yourdockerhub/drone-openai-plugin:latest
    settings:
      api_key:
        from_secret: openai_api_key
      prompt: "What are best practices for writing secure Go code?"
```

### Full Configuration

```yaml
steps:
  - name: openai-analysis
    image: yourdockerhub/drone-openai-plugin:latest
    settings:
      api_key:
        from_secret: openai_api_key
      model: gpt-4o-mini # OpenAI model to use
      prompt: "Your prompt here" # Required: The prompt to send
      file: path/to/file.ext # Optional: File to include
      system_prompt: "You are..." # System prompt for context
      temperature: 0.7 # Creativity level (0-2)
      max_tokens: 1500 # Maximum response length
      output_file: result.txt # Save response to file
      timeout: 60 # Request timeout in seconds
```

## Parameters

| Parameter       | Description                                                    | Default                        | Required |
| --------------- | -------------------------------------------------------------- | ------------------------------ | -------- |
| `api_key`       | OpenAI API key                                                 | -                              | Yes      |
| `model`         | OpenAI model to use (e.g., gpt-4o, gpt-4o-mini, gpt-3.5-turbo) | gpt-4o-mini                    | No       |
| `prompt`        | The prompt to send to OpenAI                                   | -                              | Yes      |
| `file`          | Path to file to include with prompt                            | -                              | No       |
| `system_prompt` | System message to set AI behavior                              | "You are a helpful assistant." | No       |
| `temperature`   | Controls randomness (0-2)                                      | 0.7                            | No       |
| `max_tokens`    | Maximum tokens in response                                     | 1000                           | No       |
| `output_file`   | Path to save the response                                      | -                              | No       |
| `timeout`       | Request timeout in seconds                                     | 60                             | No       |

## Supported File Types

### Text Files

- `.txt`, `.md`, `.py`, `.js`, `.go`, `.java`, `.yaml`, `.json`, etc.
- Content is appended to the prompt

### Image Files

- `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`
- Sent as base64-encoded data for vision-capable models

## Building the Plugin

### Quick Start with Makefile (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/drone-openai-plugin.git
cd drone-openai-plugin

# Build the binary
make build

# Run locally for testing
make run API_KEY="your-api-key" PROMPT="Test prompt"
```

### Manual Build

```bash
# Build the binary
go build -o drone-openai-plugin ./cmd/plugin

# Run locally for testing
export PLUGIN_API_KEY="your-api-key"
export PLUGIN_PROMPT="Test prompt"
./drone-openai-plugin
```

See [LOCAL_TESTING.md](./LOCAL_TESTING.md) for comprehensive testing guide.

### Docker Build

```bash
# Build Docker image
docker build -t drone-openai-plugin .

# Run with Docker
docker run --rm \
  -e PLUGIN_API_KEY="your-api-key" \
  -e PLUGIN_PROMPT="Test prompt" \
  -v $(pwd):/workspace \
  drone-openai-plugin
```

### Publish to Docker Hub

```bash
docker tag drone-openai-plugin yourdockerhub/drone-openai-plugin:latest
docker push yourdockerhub/drone-openai-plugin:latest
```

## Use Cases

### 1. Code Analysis and Review

```yaml
- name: code-review
  image: yourdockerhub/drone-openai-plugin:latest
  settings:
    api_key:
      from_secret: openai_api_key
    prompt: "Review this code for security issues and best practices"
    file: src/main.go
    output_file: review.md
```

### 2. Documentation Generation

```yaml
- name: generate-docs
  image: yourdockerhub/drone-openai-plugin:latest
  settings:
    api_key:
      from_secret: openai_api_key
    system_prompt: "You are a technical writer"
    prompt: "Generate user documentation for this API"
    file: openapi.yaml
    output_file: docs/api.md
```

### 3. Image Analysis

```yaml
- name: analyze-diagram
  image: yourdockerhub/drone-openai-plugin:latest
  settings:
    api_key:
      from_secret: openai_api_key
    model: gpt-4o # Vision-capable model
    prompt: "Describe this architecture and identify potential improvements"
    file: docs/architecture.png
```

### 4. Test Generation

```yaml
- name: generate-tests
  image: yourdockerhub/drone-openai-plugin:latest
  settings:
    api_key:
      from_secret: openai_api_key
    prompt: "Generate comprehensive unit tests for this module"
    file: src/auth.py
    output_file: tests/test_auth.py
```

## Environment Variables

The plugin reads configuration from environment variables prefixed with `PLUGIN_`:

- `PLUGIN_API_KEY` - OpenAI API key
- `PLUGIN_MODEL` - Model selection
- `PLUGIN_PROMPT` - Main prompt
- `PLUGIN_FILE` - File path
- `PLUGIN_SYSTEM_PROMPT` - System prompt
- `PLUGIN_TEMPERATURE` - Temperature setting
- `PLUGIN_MAX_TOKENS` - Max tokens
- `PLUGIN_OUTPUT_FILE` - Output file path
- `PLUGIN_TIMEOUT` - Timeout in seconds

## Error Handling

The plugin will fail with an appropriate error message if:

- API key is not provided
- Prompt is empty
- File specified doesn't exist
- OpenAI API returns an error
- Network timeout occurs

## Security Considerations

- Store API keys in Drone secrets, never hardcode them
- Be cautious with file permissions when saving outputs
- Consider rate limiting for production use
- Review generated content before using in production

## Limitations

- Maximum file size depends on OpenAI's API limits
- Image analysis requires vision-capable models (gpt-4o, gpt-4-turbo)
- Token limits vary by model
- API costs apply based on usage

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see [LICENSE file](./LICENSE) for details

## Support

For issues, questions, or contributions, please open an issue on GitHub.
