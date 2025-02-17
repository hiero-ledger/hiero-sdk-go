# Use the official Golang image as a base
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Create directory for the SDK
RUN mkdir -p /app/hiero-sdk-go

# Copy the SDK files first (from parent directory)
COPY . /app/hiero-sdk-go/

# Copy TCK files to their location
WORKDIR /app/hiero-sdk-go/tck

# Copy go.mod and go.sum first
COPY tck/go.mod tck/go.sum ./

# Then copy the rest of the TCK files
COPY tck/ .

ENV RUNNING_IN_DOCKER=true 

# Update and tidy the dependencies
RUN go mod tidy

# Build the Go application - specifically building server.go
RUN go build -o server cmd/server.go

# Command to run the server
CMD ["./server"]