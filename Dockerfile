# Multi-stage build for optimal image size
FROM golang:1.21-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Set Go proxy for better connectivity in China
ENV GOPROXY=https://goproxy.cn,direct

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o figureya-recommend \
    main.go load_env.go

# Final stage: minimal runtime image
FROM scratch

# Copy timezone data and CA certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Set timezone
ENV TZ=Asia/Shanghai

# Copy the binary
COPY --from=builder /app/figureya-recommend /figureya-recommend

# Copy data files and static assets
COPY --from=builder /app/figureya_docs_llm.json /figureya_docs_llm.json
COPY --from=builder /app/static /static

# Create non-root user (using existing nobody user)
USER nobody:nobody

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/figureya-recommend", "-health"] || exit 1

# Set default environment variables
ENV PORT=8080
ENV GIN_MODE=release

# Run the application
ENTRYPOINT ["/figureya-recommend"]