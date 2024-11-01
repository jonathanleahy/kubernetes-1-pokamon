# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set a working directory
WORKDIR /app

# Install necessary Alpine packages
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set a working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy templates directory
COPY --from=builder /app/templates ./templates

# Copy static directory
COPY --from=builder /app/static ./static

# Minimal dependencies (if any are needed)
RUN apk --no-cache add ca-certificates

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]