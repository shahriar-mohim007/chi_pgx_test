# Builder Stage
FROM golang:1.23.2-alpine AS builder

WORKDIR /app

# Install git for go mod download and upx for compression
RUN apk add --no-cache git upx

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source files and build with optimization flags
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o api cmd/api/main.go

# Compress binary with upx
RUN upx --best --lzma /app/api

# Final Stage
FROM scratch

WORKDIR /app
COPY --from=builder /app/api /usr/local/bin/api

# Run the compiled and compressed binary
CMD ["/usr/local/bin/api"]
