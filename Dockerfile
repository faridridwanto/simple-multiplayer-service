FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Create a minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/test-client.html .

# Expose the port the server listens on
EXPOSE 8080

# Set environment variables
ENV SESSION_LIMIT=10

# Run the server
CMD ["./server"]
