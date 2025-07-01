# Build stage
FROM golang:1.21.4-alpine AS builder

WORKDIR /app

# Copy only the module files first for caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source
COPY . .

# Build
RUN go build -o /app/main .

# Runtime stage
FROM alpine:latest

COPY --from=builder /app/main /main

CMD ["/main"]