# Multi-stage Dockerfile for Go Fiber server
FROM golang:1.22 AS builder

WORKDIR /app

# Dependencies layer
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

# Source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o exploding-kittens-server ./server.go

# Final minimal image
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary and environment file
COPY --from=builder /app/exploding-kittens-server /app/server
COPY --from=builder /app/app.env /app/app.env

# Railway assigns PORT dynamically, but expose common port for local dev
EXPOSE 8080

USER nonroot:nonroot

# Railway will override PORT, this is just a fallback
ENV PORT=8080

ENTRYPOINT ["/app/server"]
