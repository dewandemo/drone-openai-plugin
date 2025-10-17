FROM --platform=linux/amd64 golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o drone-openai-plugin ./cmd/plugin
RUN echo "=== Build complete ===" && \
    ls -lh /app/drone-openai-plugin && \
    echo "Binary size: $(du -h /app/drone-openai-plugin | cut -f1)"

FROM --platform=linux/amd64 alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /bin
COPY --from=builder /app/drone-openai-plugin /bin/drone-openai-plugin
RUN chmod +x /bin/drone-openai-plugin && \
    echo "=== Binary installed ===" && \
    ls -lh /bin/drone-openai-plugin && \
    echo "Architecture: linux/amd64"
ENTRYPOINT ["/bin/drone-openai-plugin"]