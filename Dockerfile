# Stage 1: Build the binary
FROM golang:1.22-alpine AS builder

# Install git and other required tools
RUN apk add --no-cache git

WORKDIR /app

# Copy source code
COPY . .

# Optionally, load environment variables from a .env file
COPY .env /app/.env

# Build the binary
RUN go build -o server api_server/cmd/main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest
WORKDIR /app

# Copy the built binary
COPY --from=builder /app/server /app/server

# Copy the .env file for runtime use
COPY .env /app/.env

# Set default environment variables (can be overridden at runtime)
ENV PORT=8080
ENV LOG_LEVEL=info

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./server"]
