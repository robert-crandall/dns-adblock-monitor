# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files first for better cache utilization
COPY go.mod ./

# Copy go.sum if it exists (won't fail if missing)
COPY go.sum* ./

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o /dns-checker ./src

# Final stage
FROM alpine:latest

WORKDIR /

# Copy binary from builder
COPY --from=builder /dns-checker /dns-checker

# Default environment variables
ENV DNS_HOSTS=ads.example.com,ads.example.net
ENV BLOCKING_IPV4=0.0.0.0/8,127.0.0.0/8
ENV BLOCKING_IPV6=::/128,::1/128,fc00::/7
ENV DNS_RESOLVER=1.1.1.1:53

EXPOSE 8080

CMD ["/dns-checker"]
