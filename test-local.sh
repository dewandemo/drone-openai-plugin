#!/bin/bash

# Local testing script for Drone OpenAI Plugin
# Usage: ./test-local.sh

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if API key is set
if [ -z "$OPENAI_API_KEY" ] && [ -z "$PLUGIN_API_KEY" ]; then
    echo -e "${RED}Error: OpenAI API key not set${NC}"
    echo "Set either:"
    echo "  export OPENAI_API_KEY='sk-...'"
    echo "  or"
    echo "  export PLUGIN_API_KEY='sk-...'"
    echo ""
    echo "Or create a .env.local file with your API key"
    exit 1
fi

# Use OPENAI_API_KEY if PLUGIN_API_KEY is not set
if [ -z "$PLUGIN_API_KEY" ]; then
    export PLUGIN_API_KEY="$OPENAI_API_KEY"
fi

echo -e "${BLUE}=== Drone OpenAI Plugin - Local Testing ===${NC}\n"

# Build the plugin
echo -e "${BLUE}Step 1: Building plugin...${NC}"
go build -o drone-openai-plugin ./cmd/plugin
echo -e "${GREEN}✓ Build successful${NC}\n"

# Test 1: Simple prompt
echo -e "${YELLOW}Test 1: Simple text prompt${NC}"
echo "Command: Simple greeting request"
export PLUGIN_PROMPT="Say 'Hello from your local test!' in a friendly way"
./drone-openai-plugin
echo -e "${GREEN}✓ Test 1 passed${NC}\n"

# Test 2: With file (using our own main.go)
if [ -f "cmd/plugin/main.go" ]; then
    echo -e "${YELLOW}Test 2: Code analysis with file${NC}"
    echo "Command: Analyzing cmd/plugin/main.go"
    export PLUGIN_PROMPT="Summarize what this code does in one sentence"
    export PLUGIN_FILE="cmd/plugin/main.go"
    ./drone-openai-plugin
    unset PLUGIN_FILE
    echo -e "${GREEN}✓ Test 2 passed${NC}\n"
fi

# Test 3: Save to output file
echo -e "${YELLOW}Test 3: Save response to file${NC}"
echo "Command: Saving response to test-output.txt"
export PLUGIN_PROMPT="Write a short motivational quote about software development"
export PLUGIN_OUTPUT_FILE="test-output.txt"
./drone-openai-plugin
if [ -f "test-output.txt" ]; then
    echo -e "\n${BLUE}Saved content:${NC}"
    cat test-output.txt
    rm test-output.txt
    echo -e "\n${GREEN}✓ Test 3 passed${NC}\n"
fi
unset PLUGIN_OUTPUT_FILE

# Test 4: Different model and parameters
echo -e "${YELLOW}Test 4: Custom parameters${NC}"
echo "Command: Using gpt-4o-mini with custom temperature"
export PLUGIN_PROMPT="Count from 1 to 5"
export PLUGIN_MODEL="gpt-4o-mini"
export PLUGIN_TEMPERATURE="0.1"
export PLUGIN_MAX_TOKENS="50"
./drone-openai-plugin
echo -e "${GREEN}✓ Test 4 passed${NC}\n"

# Summary
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ All tests completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Your plugin is ready to use in Drone CI!"
echo ""
echo "Next steps:"
echo "  1. Build Docker image: docker build -t yourusername/drone-openai-plugin ."
echo "  2. Push to Docker Hub: docker push yourusername/drone-openai-plugin"
echo "  3. Configure Drone secrets in your repo"
echo "  4. Add plugin step to .drone.yml"
echo ""
echo "See LOCAL_TESTING.md for more testing options"

