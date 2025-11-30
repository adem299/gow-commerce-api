# Build stage
FROM golang:1.21 as builder

WORKDIR /gow-commerce

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o gow-commerce .

# Final stage
FROM alpine:latest  

# Install CA certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN adduser -D -g '' appuser

WORKDIR /home/appuser

# Copy the binary and .env file
COPY --from=builder /gow-commerce/gow-commerce .
COPY --from=builder /gow-commerce/.env .

# Set proper permissions
RUN chown -R appuser:appuser .

# Switch to non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./gow-commerce"]
