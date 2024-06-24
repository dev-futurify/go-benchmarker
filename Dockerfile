# Start with a base image containing Go
FROM golang:1.22.4 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN go build -o main .

# Start a new stage to create a lean production image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=build /app/main .

# Install necessary utilities (if any)
RUN apk add --no-cache bash

# Set the entry point to the Go binary
ENTRYPOINT ["./main"]

# Default command-line arguments
CMD ["--cpu", "1", "--memory", "256", "--disk", "100"]
