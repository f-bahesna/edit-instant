FROM golang:1.21-alpine

# Install git, FFmpeg, and ca-certificates (CGO is NOT needed for pure Go SQLite!)
RUN apk update && apk add --no-cache git ffmpeg ca-certificates

# Install Air for Go live-reloading
RUN go install github.com/cosmtrek/air@v1.49.0

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Expose backend port
EXPOSE 8080

# Run Air for live-reloading in development
CMD ["air", "-c", ".air.toml"]
