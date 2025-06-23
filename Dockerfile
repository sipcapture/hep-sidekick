# Builder stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go application for a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /hep-sidekick ./cmd/hep-sidekick

# Final stage
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /hep-sidekick /hep-sidekick

# Set the entrypoint
ENTRYPOINT ["/hep-sidekick"] 