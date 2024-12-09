# Build Stage
FROM golang:1.22.1-bullseye AS builder

# Set the working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Run `go mod tidy` for the main module
RUN go mod tidy

# Run `go mod tidy` for submodules
WORKDIR /app/api_server
RUN go mod tidy

WORKDIR /app/builder
RUN go mod tidy

# Build the binary
WORKDIR /app
RUN go build -o server api_server/cmd/main.go

# Runtime Stage
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
# Copy the binary from the builder stage
COPY --from=builder /app/server .
RUN chmod +x ./server

# Expose the application's port
EXPOSE 8080

# Command to run the binary
CMD ["./server"]
