# Stage 1: Build the Go binaries
FROM golang:1.22 AS stage1

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy the Go workspace and source code
COPY go.work ./go.work
COPY core ./core
COPY go.mod ./go.mod
COPY go.sum ./go.sum

COPY builder ./builder
COPY server ./server

# Sync the Go workspace and download dependencies
RUN go work sync && go mod download

# Build the binaries for all services
RUN go build -o /app/bin/builder ./builder/cmd
RUN go build -o /app/bin/server ./server/cmd

# Stage 2: Create a minimal image for each service
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY --from=stage1 /app/bin/builder .
RUN chmod +x ./builder


CMD ["./builder"]

FROM golang:1.22-alpine AS server
WORKDIR /app
COPY --from=stage1 /app/bin/server .
RUN chmod +x ./server

