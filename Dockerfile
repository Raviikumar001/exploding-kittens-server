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

# Expose default port (configurable with PORT)
EXPOSE 8080

USER nonroot:nonroot

ENV PORT=8080

ENTRYPOINT ["/app/server"]
