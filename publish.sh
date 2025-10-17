#!/bin/bash

# Publishing script for Drone OpenAI Plugin
# Usage: ./publish.sh [version]
# Example: ./publish.sh 1.0.0

set -e

# Configuration (CHANGE THESE!)
DOCKER_USERNAME="${DOCKER_USERNAME:-yourusername}"  # Set your Docker Hub username
IMAGE_NAME="drone-openai-plugin"
VERSION="${1:-1.0.0}"  # Use provided version or default to 1.0.0

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Drone OpenAI Plugin - Publishing Tool    ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════╝${NC}"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running${NC}"
    echo "Please start Docker and try again"
    exit 1
fi

# Check if logged in to Docker Hub
if ! docker info 2>&1 | grep -q "Username"; then
    echo -e "${YELLOW}Not logged in to Docker Hub${NC}"
    echo "Running: docker login"
    docker login
fi

# Validate configuration
if [ "$DOCKER_USERNAME" == "yourusername" ]; then
    echo -e "${RED}Error: Please configure DOCKER_USERNAME${NC}"
    echo "Set it in this script or export DOCKER_USERNAME='your-username'"
    exit 1
fi

# Check if OPENAI_API_KEY is set for testing
if [ -z "$OPENAI_API_KEY" ] && [ -z "$PLUGIN_API_KEY" ]; then
    echo -e "${YELLOW}Warning: No API key set for testing${NC}"
    echo "Set OPENAI_API_KEY or PLUGIN_API_KEY to test the image"
    echo "Continue anyway? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        exit 1
    fi
    SKIP_TEST=true
fi

echo ""
echo -e "${BLUE}Configuration:${NC}"
echo -e "  Docker Username: ${GREEN}${DOCKER_USERNAME}${NC}"
echo -e "  Image Name:      ${GREEN}${IMAGE_NAME}${NC}"
echo -e "  Version:         ${GREEN}${VERSION}${NC}"
echo -e "  Full Image:      ${GREEN}${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}${NC}"
echo ""

# Confirm
echo -e "${YELLOW}Ready to build and publish?${NC} (y/N)"
read -r response
if [[ ! "$response" =~ ^[Yy]$ ]]; then
    echo "Cancelled"
    exit 0
fi

echo ""
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}Step 1/4: Building Docker image...${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"

docker build \
    --platform linux/amd64 \
    -t ${DOCKER_USERNAME}/${IMAGE_NAME}:latest \
    -t ${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION} \
    .

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Build successful${NC}"

# Get image info
IMAGE_SIZE=$(docker images ${DOCKER_USERNAME}/${IMAGE_NAME}:latest --format "{{.Size}}")
echo -e "  Image size: ${GREEN}${IMAGE_SIZE}${NC}"

# Test the image if API key is available
if [ "$SKIP_TEST" != "true" ]; then
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════${NC}"
    echo -e "${BLUE}Step 2/4: Testing image locally...${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════${NC}"

    # Use OPENAI_API_KEY if PLUGIN_API_KEY is not set
    if [ -z "$PLUGIN_API_KEY" ]; then
        export PLUGIN_API_KEY="$OPENAI_API_KEY"
    fi

    docker run --rm \
        -e PLUGIN_API_KEY="${PLUGIN_API_KEY}" \
        -e PLUGIN_PROMPT="Respond with exactly: 'Docker image test successful!'" \
        -e PLUGIN_MAX_TOKENS="20" \
        ${DOCKER_USERNAME}/${IMAGE_NAME}:latest

    if [ $? -ne 0 ]; then
        echo -e "${RED}✗ Test failed!${NC}"
        echo "Do you want to continue anyway? (y/N)"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        echo -e "${GREEN}✓ Test successful${NC}"
    fi
else
    echo ""
    echo -e "${YELLOW}Skipping test (Step 2/4)${NC}"
fi

echo ""
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}Step 3/4: Pushing to Docker Hub...${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"

echo "Pushing latest..."
docker push ${DOCKER_USERNAME}/${IMAGE_NAME}:latest

echo "Pushing version ${VERSION}..."
docker push ${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Push failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Push successful${NC}"

echo ""
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}Step 4/4: Creating git tag...${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"

# Check if tag already exists
if git rev-parse "v${VERSION}" >/dev/null 2>&1; then
    echo -e "${YELLOW}Tag v${VERSION} already exists${NC}"
else
    echo "Create git tag v${VERSION}? (y/N)"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        git tag -a "v${VERSION}" -m "Release version ${VERSION}"
        echo -e "${GREEN}✓ Tag created: v${VERSION}${NC}"
        echo ""
        echo "Push tag to remote? (y/N)"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            git push origin "v${VERSION}"
            echo -e "${GREEN}✓ Tag pushed to remote${NC}"
        fi
    else
        echo "Skipped"
    fi
fi

# Summary
echo ""
echo -e "${GREEN}╔════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║          ✓ Published Successfully!         ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}Your plugin is now available at:${NC}"
echo -e "  ${GREEN}https://hub.docker.com/r/${DOCKER_USERNAME}/${IMAGE_NAME}${NC}"
echo ""
echo -e "${BLUE}Use in Drone CI:${NC}"
echo ""
echo "  kind: pipeline"
echo "  name: default"
echo ""
echo "  steps:"
echo "    - name: openai-task"
echo -e "      image: ${GREEN}${DOCKER_USERNAME}/${IMAGE_NAME}:latest${NC}"
echo "      settings:"
echo "        api_key:"
echo "          from_secret: openai_api_key"
echo "        prompt: \"Your prompt here\""
echo ""
echo -e "${BLUE}Or use specific version:${NC}"
echo -e "      image: ${GREEN}${DOCKER_USERNAME}/${IMAGE_NAME}:${VERSION}${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "  1. Visit Docker Hub to verify: hub.docker.com/r/${DOCKER_USERNAME}/${IMAGE_NAME}"
echo "  2. Update your .drone.yml to use the published image"
echo "  3. Configure OpenAI API key as a secret in your Drone repository"
echo "  4. Test in a real Drone pipeline"
echo ""
echo "Happy building! 🚀"

