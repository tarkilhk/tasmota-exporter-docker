# Build stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o tasmota-exporter ./cmd/tasmota-exporter

# Final stage
FROM alpine:latest
RUN apk add --no-cache tzdata

WORKDIR /app

# Copy the binary from builder and ensure it's executable
COPY --from=builder /app/tasmota-exporter .
RUN chmod +x /app/tasmota-exporter

# Expose the port the app runs on
EXPOSE 9090

# Run the application
CMD ["./tasmota-exporter"] 