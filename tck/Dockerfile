FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the SDK files first (from parent directory)
COPY . /app/hiero-sdk-go/

# Set the working directory for the TCK files
WORKDIR /app/hiero-sdk-go/tck

# Copy go.mod and go.sum first
COPY tck/go.mod tck/go.sum ./

RUN go mod tidy
RUN go build -o server cmd/server.go


FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/hiero-sdk-go/tck/server .
RUN chmod +x /app/server

# Install necessary libraries
RUN apk add --no-cache libc6-compat

CMD ["./server"]