FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o drone-openai-plugin ./cmd/plugin
RUN ls -la /app/drone-openai-plugin

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /bin
COPY --from=builder /app/drone-openai-plugin /bin/drone-openai-plugin
RUN chmod +x /bin/drone-openai-plugin
RUN ls -la /bin/drone-openai-plugin
ENTRYPOINT ["/bin/drone-openai-plugin"]