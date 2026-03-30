# syntax=docker/dockerfile:1

# Builder stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd

# Runtime stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /usr/local/bin/server

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["/usr/local/bin/server", "-health"]
ENTRYPOINT ["/usr/local/bin/server"]
