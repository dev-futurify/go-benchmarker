# Use the official Golang image as a build stage
FROM golang:1.22.4 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod ./

# Download Go modules
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Add Maintainer info
LABEL maintainer="adri@satusky.com"

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the binary with default flags
CMD ["./main", "--cpu", "1", "--memory", "256", "--disk", "1000"]


