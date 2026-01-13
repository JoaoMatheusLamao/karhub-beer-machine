# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (needed for go modules)
RUN apk add --no-cache git

# Copy go mod files first (cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o api ./cmd/api/main.go

# Build the CLI binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o cli ./cmd/cli/main.go

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/api /app/api
COPY --from=builder /app/cli /app/cli

# Expose HTTP port
EXPOSE 8080

# Run API by default
ENTRYPOINT ["/app/api"]
