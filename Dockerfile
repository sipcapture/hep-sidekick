# Builder stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /app/hep-sidekick ./cmd/hep-sidekick

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata libcap
COPY --from=builder /app/hep-sidekick /hep-sidekick
RUN /usr/sbin/setcap cap_bpf,cap_net_raw,cap_net_admin=eip hep-sidekick
CMD ["./heplify", "-h"]
