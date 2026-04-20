# Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code and build
COPY . .
RUN go build -o /fs-node

# Run Stage
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /fs-node /app/fs-node

# Start the server
CMD ["/app/fs-node"]