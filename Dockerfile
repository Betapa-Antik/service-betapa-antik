# Step 1: Build the binary
FROM golang:1.23.1-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install git, SSL/TLS certificates and other build essentials
RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application as a statically compiled binary
# -ldflags="-w -s" reduces binary size by stripping debug symbols
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o betapa-antik-service main.go

# Step 2: Create a minimal production runner image
FROM alpine:3.19

# Set metadata labels
LABEL maintainer="Antigravity Team"
LABEL version="1.0.0"

# Install SSL certificates and timezone database (tzdata)
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user and group for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/betapa-antik-service .

# Change ownership of the application directory to the non-root user
RUN chown -R appuser:appgroup /app

# Use the non-root user to run the application
USER appuser

# Expose the default application port
EXPOSE 8081

# Run the application
CMD ["./betapa-antik-service"]
