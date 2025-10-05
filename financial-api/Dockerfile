
# Stage 1: Build the Go binary
FROM golang:1.24 AS builder

WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

# Stage 2: Create a minimal final image
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /server /server

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["/server"]
